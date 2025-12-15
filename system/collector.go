package system

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	AF_INET                 = 2
	TCP_TABLE_OWNER_PID_ALL = 5
	UDP_TABLE_OWNER_PID     = 1
	MIB_TCP_STATE_LISTEN    = 2
)

// --------------------
// WinAPI structs
// --------------------

type MIB_TCPROW_OWNER_PID struct {
	State      uint32
	LocalAddr  uint32
	LocalPort  uint32
	RemoteAddr uint32
	RemotePort uint32
	OwningPid  uint32
}

type MIB_TCPTABLE_OWNER_PID struct {
	NumEntries uint32
	Table      [1]MIB_TCPROW_OWNER_PID
}

type MIB_UDPROW_OWNER_PID struct {
	LocalAddr uint32
	LocalPort uint32
	OwningPid uint32
}

type MIB_UDPTABLE_OWNER_PID struct {
	NumEntries uint32
	Table      [1]MIB_UDPROW_OWNER_PID
}

// --------------------
// DLL imports
// --------------------

var (
	modiphlpapi             = windows.NewLazySystemDLL("iphlpapi.dll")
	procGetExtendedTcpTable = modiphlpapi.NewProc("GetExtendedTcpTable")
	procGetExtendedUdpTable = modiphlpapi.NewProc("GetExtendedUdpTable")
)

// --------------------
// Helpers
// --------------------

func ntohs(port uint32) uint16 {
	p := uint16(port & 0xFFFF)
	return (p >> 8) | (p << 8)
}

func addrToString(addr uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(addr),
		byte(addr>>8),
		byte(addr>>16),
		byte(addr>>24))
}

func tcpStateToString(state uint32) string {
	switch state {
	case 1:
		return "CLOSED"
	case 2:
		return "LISTEN"
	case 3:
		return "SYN_SENT"
	case 4:
		return "SYN_RCVD"
	case 5:
		return "ESTABLISHED"
	case 6:
		return "FIN_WAIT1"
	case 7:
		return "FIN_WAIT2"
	case 8:
		return "CLOSE_WAIT"
	case 9:
		return "CLOSING"
	case 10:
		return "LAST_ACK"
	case 11:
		return "TIME_WAIT"
	case 12:
		return "DELETE_TCB"
	default:
		return "UNKNOWN"
	}
}

func processName(pid uint32) string {
	return Instance.GetProcessName(pid)
}

// --------------------
// Data model
// --------------------

type ProcInfo struct {
	PID  uint32
	Name string
}

type PortMap map[uint16][]ProcInfo

// ConnectionInfo contains detailed information about a single network connection
type ConnectionInfo struct {
	Protocol    string // "TCP" or "UDP"
	LocalAddr   string // Local side IP address
	LocalPort   uint16 // Local side port
	RemoteAddr  string // Remote side IP address (for TCP)
	RemotePort  uint16 // Remote side port (for TCP)
	State       string // Connection state (for TCP)
	PID         uint32 // Process ID
	ProcessName string // Process name
}

// NetworkConnections contains all network connections
type NetworkConnections struct {
	Connections []ConnectionInfo
}

// --------------------
// Collectors
// --------------------

func collectTCP() PortMap {
	ports := make(PortMap)

	var size uint32
	procGetExtendedTcpTable.Call(
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		AF_INET,
		TCP_TABLE_OWNER_PID_ALL,
		0,
	)

	buf := make([]byte, size)
	ret, _, _ := procGetExtendedTcpTable.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		AF_INET,
		TCP_TABLE_OWNER_PID_ALL,
		0,
	)
	if ret != 0 {
		return ports
	}

	table := (*MIB_TCPTABLE_OWNER_PID)(unsafe.Pointer(&buf[0]))
	rows := (*[1 << 20]MIB_TCPROW_OWNER_PID)(
		unsafe.Pointer(&table.Table[0]),
	)[:table.NumEntries:table.NumEntries]

	for _, row := range rows {
		// Show only listening ports (LISTEN)
		if row.State != MIB_TCP_STATE_LISTEN {
			continue
		}

		port := ntohs(row.LocalPort)
		ports[port] = append(ports[port], ProcInfo{
			PID:  row.OwningPid,
			Name: processName(row.OwningPid),
		})
	}

	return ports
}

func collectUDP() PortMap {
	ports := make(PortMap)

	var size uint32
	procGetExtendedUdpTable.Call(
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		AF_INET,
		UDP_TABLE_OWNER_PID,
		0,
	)

	buf := make([]byte, size)
	ret, _, _ := procGetExtendedUdpTable.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		AF_INET,
		UDP_TABLE_OWNER_PID,
		0,
	)
	if ret != 0 {
		return ports
	}

	table := (*MIB_UDPTABLE_OWNER_PID)(unsafe.Pointer(&buf[0]))
	rows := (*[1 << 20]MIB_UDPROW_OWNER_PID)(
		unsafe.Pointer(&table.Table[0]),
	)[:table.NumEntries:table.NumEntries]

	for _, row := range rows {
		port := ntohs(row.LocalPort)
		ports[port] = append(ports[port], ProcInfo{
			PID:  row.OwningPid,
			Name: processName(row.OwningPid),
		})
	}

	return ports
}

func collectAllTCPConnections() []ConnectionInfo {
	var connections []ConnectionInfo

	var size uint32
	procGetExtendedTcpTable.Call(
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		AF_INET,
		TCP_TABLE_OWNER_PID_ALL,
		0,
	)

	buf := make([]byte, size)
	ret, _, _ := procGetExtendedTcpTable.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		AF_INET,
		TCP_TABLE_OWNER_PID_ALL,
		0,
	)
	if ret != 0 {
		return connections
	}

	table := (*MIB_TCPTABLE_OWNER_PID)(unsafe.Pointer(&buf[0]))
	rows := (*[1 << 20]MIB_TCPROW_OWNER_PID)(
		unsafe.Pointer(&table.Table[0]),
	)[:table.NumEntries:table.NumEntries]

	for _, row := range rows {
		connections = append(connections, ConnectionInfo{
			Protocol:    "TCP",
			LocalAddr:   addrToString(row.LocalAddr),
			LocalPort:   ntohs(row.LocalPort),
			RemoteAddr:  addrToString(row.RemoteAddr),
			RemotePort:  ntohs(row.RemotePort),
			State:       tcpStateToString(row.State),
			PID:         row.OwningPid,
			ProcessName: processName(row.OwningPid),
		})
	}

	return connections
}

func collectAllUDPConnections() []ConnectionInfo {
	var connections []ConnectionInfo

	var size uint32
	procGetExtendedUdpTable.Call(
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		AF_INET,
		UDP_TABLE_OWNER_PID,
		0,
	)

	buf := make([]byte, size)
	ret, _, _ := procGetExtendedUdpTable.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		AF_INET,
		UDP_TABLE_OWNER_PID,
		0,
	)
	if ret != 0 {
		return connections
	}

	table := (*MIB_UDPTABLE_OWNER_PID)(unsafe.Pointer(&buf[0]))
	rows := (*[1 << 20]MIB_UDPROW_OWNER_PID)(
		unsafe.Pointer(&table.Table[0]),
	)[:table.NumEntries:table.NumEntries]

	for _, row := range rows {
		connections = append(connections, ConnectionInfo{
			Protocol:    "UDP",
			LocalAddr:   addrToString(row.LocalAddr),
			LocalPort:   ntohs(row.LocalPort),
			RemoteAddr:  "",
			RemotePort:  0,
			State:       "",
			PID:         row.OwningPid,
			ProcessName: processName(row.OwningPid),
		})
	}

	return connections
}

// GetAllConnections returns information about all network connections (TCP and UDP)
func GetAllConnections() NetworkConnections {
	var result NetworkConnections

	// Collect TCP connections
	tcpConnections := collectAllTCPConnections()
	result.Connections = append(result.Connections, tcpConnections...)

	// Collect UDP connections
	udpConnections := collectAllUDPConnections()
	result.Connections = append(result.Connections, udpConnections...)

	return result
}
