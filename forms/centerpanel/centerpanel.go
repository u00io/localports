package centerpanel

import (
	"fmt"
	"image/color"
	"sort"
	"sync"
	"time"

	"github.com/u00io/localports/system"
	"github.com/u00io/nuiforms/ui"
)

type CenterPanel struct {
	ui.Widget

	mtx  sync.Mutex
	data system.NetworkConnections

	orderColumnIndex int
	orderAsc         bool

	tableResults *ui.Table
}

func NewCenterPanel() *CenterPanel {
	var c CenterPanel
	c.InitWidget()
	c.tableResults = ui.NewTable()
	curstomWidgets := map[string]ui.Widgeter{
		"tableresults": c.tableResults,
	}
	c.SetLayout(`
		<row>
			<widget id="tableresults" />
		</row>
	`, &c, curstomWidgets)

	c.orderColumnIndex = 1

	// c.tableResults

	c.tableResults.SetColumnCount(10)
	c.tableResults.SetColumnName(0, "Type")
	c.tableResults.SetColumnName(1, "Local Port")
	c.tableResults.SetColumnName(2, "Local Address")
	c.tableResults.SetColumnName(3, "Remote Address")
	c.tableResults.SetColumnName(4, "Remote Port")
	c.tableResults.SetColumnName(5, "Status")
	c.tableResults.SetColumnName(6, "PID")
	c.tableResults.SetColumnName(7, "Program")
	c.tableResults.SetColumnName(8, "Service")
	c.tableResults.SetColumnName(9, "Country")

	c.tableResults.SetColumnWidth(0, 120)
	c.tableResults.SetColumnWidth(1, 160)
	c.tableResults.SetColumnWidth(2, 220)
	c.tableResults.SetColumnWidth(3, 220)
	c.tableResults.SetColumnWidth(4, 160)
	c.tableResults.SetColumnWidth(5, 160)
	c.tableResults.SetColumnWidth(6, 80)
	c.tableResults.SetColumnWidth(7, 220)
	c.tableResults.SetColumnWidth(8, 220)
	c.tableResults.SetColumnWidth(9, 180)

	c.tableResults.SetOnColumnClick(c.OnColumnHeaderClicked)
	c.updateColumns()

	go c.thUpdateData()

	return &c
}

func (c *CenterPanel) thUpdateData() {
	for {
		conns := system.GetAllConnections()
		c.mtx.Lock()
		c.data = conns
		c.mtx.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func (c *CenterPanel) HandleSystemEvent(event system.Event) {
	if event.Name == "update" {
		c.updateData()
	}
}

func (c *CenterPanel) ColumnName(index int) string {
	switch index {
	case 0:
		return "Type"
	case 1:
		return "Local Port"
	case 2:
		return "Local Address"
	case 3:
		return "Remote Address"
	case 4:
		return "Remote Port"
	case 5:
		return "Status"
	case 6:
		return "PID"
	case 7:
		return "Program"
	case 8:
		return "Service"
	default:
		return ""
	}
}

func (c *CenterPanel) OnColumnHeaderClicked(index int) {
	if c.orderColumnIndex == index {
		c.orderAsc = !c.orderAsc
	} else {
		c.orderColumnIndex = index
		c.orderAsc = true
	}
	c.updateColumns()
	c.updateData()
}

func (c *CenterPanel) updateColumns() {
	// set columns names based on sorting
	for i := 0; i < 9; i++ {
		name := c.ColumnName(i)
		if i == c.orderColumnIndex {
			if c.orderAsc {
				name = name + " [^]"
			} else {
				name = name + " [v]"
			}
		}
		c.tableResults.SetColumnName(i, name)
	}
}

func (c *CenterPanel) updateData() {
	conns := make([]system.ConnectionInfo, 0)

	filterType := system.Instance.GetFilterType()
	filterStatus := system.Instance.GetFilterStatus()
	for _, conn := range c.data.Connections {
		if filterType == "tcp" && conn.Protocol != "TCP" {
			continue
		}
		if filterType == "udp" && conn.Protocol != "UDP" {
			continue
		}
		if filterStatus == "LISTEN" && conn.State != "LISTEN" && (filterType == "tcp" || filterType == "all") {
			continue
		}
		if filterStatus == "ESTABLISHED" && conn.State != "ESTABLISHED" && (filterType == "tcp" || filterType == "all") {
			continue
		}
		if filterStatus == "OTHER" && (conn.State == "LISTEN" || conn.State == "ESTABLISHED") && (filterType == "tcp" || filterType == "all") {
			continue
		}
		conns = append(conns, conn)
	}

	sort.Slice(conns, func(i, j int) bool {
		switch c.orderColumnIndex {
		case 0:
			if c.orderAsc {
				return conns[i].Protocol < conns[j].Protocol
			} else {
				return conns[i].Protocol > conns[j].Protocol
			}
		case 1:
			if c.orderAsc {
				return conns[i].LocalPort < conns[j].LocalPort
			} else {
				return conns[i].LocalPort > conns[j].LocalPort
			}
		case 2:
			if c.orderAsc {
				return conns[i].LocalAddr < conns[j].LocalAddr
			} else {
				return conns[i].LocalAddr > conns[j].LocalAddr
			}
		case 3:
			if c.orderAsc {
				return conns[i].RemoteAddr < conns[j].RemoteAddr
			} else {
				return conns[i].RemoteAddr > conns[j].RemoteAddr
			}
		case 4:
			if c.orderAsc {
				return conns[i].RemotePort < conns[j].RemotePort
			} else {
				return conns[i].RemotePort > conns[j].RemotePort
			}
		case 5:
			if c.orderAsc {
				return conns[i].State < conns[j].State
			} else {
				return conns[i].State > conns[j].State
			}
		case 6:
			if c.orderAsc {
				return conns[i].PID < conns[j].PID
			} else {
				return conns[i].PID > conns[j].PID
			}
		case 7:
			if c.orderAsc {
				return conns[i].ProcessName < conns[j].ProcessName
			} else {
				return conns[i].ProcessName > conns[j].ProcessName
			}
		case 8:
			serviceA := system.Instance.GetServiceByPort(conns[i].LocalPort)
			if serviceA == "" {
				serviceA = system.Instance.GetServiceByPort(conns[i].RemotePort)
			}
			serviceB := system.Instance.GetServiceByPort(conns[j].LocalPort)
			if serviceB == "" {
				serviceB = system.Instance.GetServiceByPort(conns[j].RemotePort)
			}
			if c.orderAsc {
				return serviceA < serviceB
			}
			return serviceA > serviceB
		case 9:
			countryA, errA := system.GetCountryByIP(conns[i].RemoteAddr)
			if errA != nil {
				countryA = ""
			}
			countryB, errB := system.GetCountryByIP(conns[j].RemoteAddr)
			if errB != nil {
				countryB = ""
			}
			if c.orderAsc {
				return countryA < countryB
			} else {
				return countryA > countryB
			}
		default:
			return true
		}
	})

	c.tableResults.SetRowCount(len(conns))
	for i, conn := range conns {
		// TYPE
		c.tableResults.SetCellText2(i, 0, conn.Protocol)

		// LOCAL PORT
		c.tableResults.SetCellText2(i, 1, fmt.Sprintf("%d", conn.LocalPort))

		// LOCAL ADDRESS
		c.tableResults.SetCellText2(i, 2, conn.LocalAddr)
		c.tableResults.SetCellColor(i, 2, color.RGBA{100, 100, 100, 255})

		// REMOTE ADDRESS
		c.tableResults.SetCellText2(i, 3, conn.RemoteAddr)
		c.tableResults.SetCellColor(i, 3, ui.ColorFromHex("#E57373"))
		if conn.State == "LISTEN" {
			c.tableResults.SetCellText2(i, 3, "")
		}
		if conn.RemoteAddr == "0.0.0.0" || conn.RemoteAddr == "::" || conn.RemoteAddr == "127.0.0.1" {
			c.tableResults.SetCellColor(i, 3, color.RGBA{100, 100, 100, 255})
		}
		if system.Instance.IsLocalAreaNetwork(conn.RemoteAddr) {
			c.tableResults.SetCellColor(i, 3, color.RGBA{100, 255, 100, 255})
		}

		// REMOTE PORT
		if conn.RemotePort > 0 {
			c.tableResults.SetCellText2(i, 4, fmt.Sprintf("%d", conn.RemotePort))
		} else {
			c.tableResults.SetCellText2(i, 4, "")
		}

		// STATUS
		c.tableResults.SetCellText2(i, 5, conn.State)

		// PID
		c.tableResults.SetCellText2(i, 6, fmt.Sprintf("%d", conn.PID))
		c.tableResults.SetCellColor(i, 6, color.RGBA{100, 100, 100, 255})

		// PROGRAM
		c.tableResults.SetCellText2(i, 7, conn.ProcessName)

		// SERVICE
		service := system.Instance.GetServiceByPort(conn.LocalPort)
		if service == "" {
			service = system.Instance.GetServiceByPort(conn.RemotePort)
		}
		c.tableResults.SetCellText2(i, 8, service)

		// COUNTRY
		country, err := system.GetCountryByIP(conn.RemoteAddr)
		if err != nil {
			country = ""
		}
		c.tableResults.SetCellText2(i, 9, country)
	}
}
