package system

import (
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type System struct {
	mtx sync.Mutex

	events []Event

	filterType   string
	filterStatus string

	processNamesById map[uint32]string
}

type Event struct {
	Name      string
	Parameter string
}

var Instance *System

func NewSystem() *System {
	var c System
	return &c
}

func (c *System) Start() {
	go c.thUpdateProcesses()
}

func (c *System) Stop() {
}

func (c *System) thUpdateProcesses() {
	for {
		c.updateProcesses()
		time.Sleep(1 * time.Second)
	}
}

func (c *System) updateProcesses() {
	result := make(map[uint32]string)

	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err == nil {
		var entry windows.ProcessEntry32
		entry.Size = uint32(unsafe.Sizeof(entry))
		err = windows.Process32First(handle, &entry)
		for err == nil {
			nameSize := 0
			for i := 0; i < 260; i++ {
				if entry.ExeFile[nameSize] == 0 {
					break
				}
				nameSize++
			}

			id := int(entry.ProcessID)
			name := syscall.UTF16ToString(entry.ExeFile[:nameSize])
			result[uint32(id)] = name
			err = windows.Process32Next(handle, &entry)
		}

		_ = windows.CloseHandle(handle)
	}

	c.mtx.Lock()
	c.processNamesById = result
	c.mtx.Unlock()
}

func (c *System) SetFilterType(filterType string) {
	c.mtx.Lock()
	c.filterType = filterType
	c.mtx.Unlock()
}

func (c *System) SetFilterStatus(filterStatus string) {
	c.mtx.Lock()
	c.filterStatus = filterStatus
	c.mtx.Unlock()
}

func (c *System) GetFilterType() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.filterType
}

func (c *System) GetFilterStatus() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.filterStatus
}

func (c *System) EmitEvent(event string, parameter string) {
	c.mtx.Lock()
	c.events = append(c.events, Event{Name: event, Parameter: parameter})
	c.mtx.Unlock()
}

func (c *System) GetAndClearEvents() []Event {
	c.mtx.Lock()
	events := c.events
	c.events = make([]Event, 0)
	c.mtx.Unlock()
	return events
}

func (c *System) IsLocalAreaNetwork(ip string) bool {
	if len(ip) >= 4 {
		if ip[:4] == "10." {
			return true
		}
		if len(ip) >= 8 {
			if ip[:8] == "192.168." {
				return true
			}
			if len(ip) >= 7 {
				if ip[:7] == "172.16." {
					return true
				}
			}
		}
	}
	return false
}

func (c *System) GetProcessName(pid uint32) string {
	result := "?"
	c.mtx.Lock()
	if name, ok := c.processNamesById[pid]; ok {
		result = name
	}
	c.mtx.Unlock()
	return result
}

func (c *System) GetServiceByPort(port uint16) string {
	if service, ok := portServiceMap[port]; ok {
		return service
	}
	return ""
}

var portServiceMap = map[uint16]string{
	// --- Well-known ports (0–1023)
	20:  "FTP Data",
	21:  "FTP Control",
	22:  "SSH",
	23:  "Telnet",
	25:  "SMTP",
	37:  "Time",
	42:  "WINS Replication",
	43:  "WHOIS",
	49:  "TACACS",
	53:  "DNS",
	67:  "DHCP Server",
	68:  "DHCP Client",
	69:  "TFTP",
	70:  "Gopher",
	79:  "Finger",
	80:  "HTTP",
	88:  "Kerberos",
	110: "POP3",
	119: "NNTP",
	123: "NTP",
	135: "MS RPC",
	137: "NetBIOS Name",
	138: "NetBIOS Datagram",
	139: "NetBIOS Session",
	143: "IMAP",
	161: "SNMP",
	162: "SNMP Trap",
	179: "BGP",
	194: "IRC",
	199: "SMUX",
	389: "LDAP",
	427: "SLP",
	443: "HTTPS",
	445: "SMB/CIFS",
	465: "SMTPS",
	500: "ISAKMP / IKE",
	512: "rexec",
	513: "rlogin",
	514: "syslog",
	515: "LPD",
	520: "RIP",
	587: "SMTP Submission",
	636: "LDAPS",
	989: "FTPS Data",
	990: "FTPS Control",
	993: "IMAPS",
	995: "POP3S",

	// --- Registered ports (1024–49151)
	1080: "SOCKS Proxy",
	1433: "MSSQL",
	1521: "Oracle DB",
	2049: "NFS",
	2082: "cPanel",
	2083: "cPanel SSL",
	2086: "WHM",
	2087: "WHM SSL",
	2181: "Zookeeper",
	2375: "Docker",
	2376: "Docker TLS",
	2483: "Oracle TCPS",
	2484: "Oracle TCPS",
	3000: "Generic Web App",
	3001: "Generic Web App",
	3306: "MySQL",
	3389: "RDP",
	3690: "Subversion",
	4000: "Generic Web App",
	4444: "Metasploit",
	5432: "PostgreSQL",
	5601: "Kibana",
	5672: "AMQP (RabbitMQ)",
	5900: "VNC",
	5985: "WinRM HTTP",
	5986: "WinRM HTTPS",
	6060: "pprof",
	6379: "Redis",
	6443: "Kubernetes API",
	6667: "IRC",
	7001: "WebLogic",
	7002: "WebLogic SSL",
	7077: "Spark",
	8000: "HTTP Alt",
	8008: "HTTP Alt",
	8080: "HTTP Proxy",
	8081: "HTTP Alt",
	8443: "HTTPS Alt",
	9000: "Generic Service",
	9042: "Cassandra",
	9092: "Kafka",
	9200: "Elasticsearch",
	9418: "Git",
	9999: "Debug Service",

	// --- Common dev / infra
	27017: "MongoDB",
	27018: "MongoDB Shard",
	27019: "MongoDB Config",
	28017: "MongoDB Web",

	// --- Industrial / IoT / SCADA
	502:   "Modbus TCP",
	102:   "Siemens S7",
	1911:  "Tridium Niagara",
	20000: "DNP3",

	// --- VPN
	1194: "OpenVPN",
	1701: "L2TP",
	1723: "PPTP",
	4500: "IPsec NAT-T",

	// --- Blockchain
	8332:  "Bitcoin RPC",
	8333:  "Bitcoin P2P",
	30303: "Ethereum",
	8545:  "Ethereum RPC",
	26656: "Tendermint P2P",
	26657: "Tendermint RPC",

	// --- Monitoring
	9090: "Prometheus",
	9100: "Node Exporter",
}
