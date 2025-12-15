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

	c.tableResults.SetColumnCount(9)
	c.tableResults.SetColumnName(0, "Type")
	c.tableResults.SetColumnName(1, "Port")
	c.tableResults.SetColumnName(2, "Status")
	c.tableResults.SetColumnName(3, "PID")
	c.tableResults.SetColumnName(4, "Program")

	c.tableResults.SetColumnName(5, "Remote Address")
	c.tableResults.SetColumnName(6, "Remote Port")
	c.tableResults.SetColumnName(7, "Local Address")
	c.tableResults.SetColumnName(8, "Service")

	c.tableResults.SetColumnWidth(0, 100)
	c.tableResults.SetColumnWidth(1, 100)
	c.tableResults.SetColumnWidth(2, 200)
	c.tableResults.SetColumnWidth(3, 100)
	c.tableResults.SetColumnWidth(4, 200)
	c.tableResults.SetColumnWidth(5, 200)
	c.tableResults.SetColumnWidth(6, 130)
	c.tableResults.SetColumnWidth(7, 200)
	c.tableResults.SetColumnWidth(8, 200)

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
			return conns[i].Protocol < conns[j].Protocol
		case 1:
			return conns[i].LocalPort < conns[j].LocalPort
		case 2:
			return conns[i].State < conns[j].State
		case 3:
			return conns[i].PID < conns[j].PID
		case 4:
			return conns[i].ProcessName < conns[j].ProcessName
		default:
			return true
		}
	})

	c.tableResults.SetRowCount(len(conns))
	for i, conn := range conns {
		c.tableResults.SetCellText2(i, 0, conn.Protocol)
		c.tableResults.SetCellText2(i, 1, fmt.Sprintf("%d", conn.LocalPort))
		c.tableResults.SetCellText2(i, 2, conn.State)
		c.tableResults.SetCellText2(i, 3, fmt.Sprintf("%d", conn.PID))
		c.tableResults.SetCellText2(i, 4, conn.ProcessName)

		c.tableResults.SetCellText2(i, 5, conn.RemoteAddr)
		c.tableResults.SetCellColor(i, 5, ui.ColorFromHex("#E57373"))

		if conn.RemoteAddr == "0.0.0.0" || conn.RemoteAddr == "::" || conn.RemoteAddr == "127.0.0.1" {
			c.tableResults.SetCellColor(i, 5, color.RGBA{100, 100, 100, 255})
		}

		if system.Instance.IsLocalAreaNetwork(conn.RemoteAddr) {
			c.tableResults.SetCellColor(i, 5, color.RGBA{100, 255, 100, 255})
		}

		if conn.RemotePort > 0 {
			c.tableResults.SetCellText2(i, 6, fmt.Sprintf("%d", conn.RemotePort))
		} else {
			c.tableResults.SetCellText2(i, 6, "")
		}

		c.tableResults.SetCellText2(i, 7, conn.LocalAddr)
		c.tableResults.SetCellColor(i, 7, color.RGBA{100, 100, 100, 255})

		service := system.Instance.GetServiceByPort(conn.LocalPort)
		if service == "" {
			service = system.Instance.GetServiceByPort(conn.RemotePort)
		}
		c.tableResults.SetCellText2(i, 8, service)
	}
}
