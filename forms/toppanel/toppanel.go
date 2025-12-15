package toppanel

import (
	"github.com/u00io/localports/system"
	"github.com/u00io/nuiforms/ui"
)

type TopPanel struct {
	ui.Widget

	autoupdateOn bool

	filterType   string
	filterStatus string

	firstUpdateDone bool
}

func NewTopPanel() *TopPanel {
	var c TopPanel
	c.InitWidget()
	c.SetElevation(5)
	c.SetLayout(`
		<row>
			<column pagging="0" spacing="0">
				<label text="Updating" textAlign="center"/>
				<panel />
				<frame autofillbackground="true" padding="2" />
				<panel />
				<row padding="0" spacing="0">
					<button text="Update" onclick="OnUpdateClick" />
					<panel />
					<button id="btnAutoupdate" text="Auto" onclick="OnAutoUpdateClick" />
				</row>
			</column>		

			<panel padding="2" autofillbackground="true"/>

			<column pagging="0" spacing="0">
				<label text="Type" textAlign="center"/>
				<panel />
				<frame autofillbackground="true" padding="2" />
				<panel />
				<row padding="0" spacing="0">
					<button id="btnTcp" text="TCP" onclick="OnTcpClick" />
					<panel />
					<button id="btnUdp" text="UDP" onclick="OnUdpClick" />
					<panel />
					<button id="btnAll" text="ALL" onclick="OnAllClick" />
				</row>
			</column>

			<panel padding="2" autofillbackground="true"/>

			<column pagging="0" spacing="0">
				<label id="lblStatus" text="Status" textAlign="center"/>
				<panel />
				<frame autofillbackground="true" padding="2" />
				<panel />
				<row padding="0" spacing="0">
					<button id="btnStatusListen" text="LISTEN" onclick="OnStatusListenClick" />
					<panel />
					<button id="btnStatusEstablished" text="ESTABLISHED" onclick="OnStatusEstablishedClick" />
					<panel />
					<button id="btnStatusOther" text="OTHER" onclick="OnStatusOtherClick" />
					<panel />
					<button id="btnStatusAll" text="ALL" onclick="OnStatusAllClick" />
				</row>
			</column>

			<hspacer />
		</row>
	`, &c, nil)

	c.autoupdateOn = true
	c.AddTimer(1000, c.timerUpdate)
	c.updateAutoupdateButton()

	c.filterType = "tcp"
	c.updateTypeButtons()
	system.Instance.SetFilterType(c.filterType)

	c.filterStatus = "LISTEN"
	c.updateStatusButtons()
	system.Instance.SetFilterStatus(c.filterStatus)

	btnEstablished, ok := c.FindWidgetByName("btnStatusEstablished").(*ui.Button)
	if ok {
		btnEstablished.SetMinWidth(150)
	}

	return &c
}

func (c *TopPanel) timerUpdate() {
	if !c.firstUpdateDone {
		c.firstUpdateDone = true
		c.EmitUpdateEvent()
	}
	if c.autoupdateOn {
		c.EmitUpdateEvent()
	}
}

func (c *TopPanel) EmitUpdateEvent() {
	system.Instance.EmitEvent("update", "")
}

func (c *TopPanel) OnUpdateClick() {
	c.EmitUpdateEvent()
}

func (c *TopPanel) OnTcpClick() {
	c.filterType = "tcp"
	system.Instance.SetFilterType(c.filterType)
	c.updateTypeButtons()
	c.EmitUpdateEvent()
	c.updateStatusButtons()
}

func (c *TopPanel) OnUdpClick() {
	c.filterType = "udp"
	system.Instance.SetFilterType(c.filterType)
	c.updateTypeButtons()
	c.EmitUpdateEvent()
	c.updateStatusButtons()
}

func (c *TopPanel) OnAllClick() {
	c.filterType = "all"
	system.Instance.SetFilterType(c.filterType)
	c.updateTypeButtons()
	c.EmitUpdateEvent()
	c.updateStatusButtons()
}

func (c *TopPanel) OnAutoUpdateClick() {
	c.autoupdateOn = !c.autoupdateOn
	if c.autoupdateOn {
		c.EmitUpdateEvent()
	}
	c.updateAutoupdateButton()
}

func (c *TopPanel) OnStatusListenClick() {
	c.filterStatus = "LISTEN"
	system.Instance.SetFilterStatus(c.filterStatus)
	c.updateStatusButtons()
	c.EmitUpdateEvent()
}

func (c *TopPanel) OnStatusEstablishedClick() {
	c.filterStatus = "ESTABLISHED"
	system.Instance.SetFilterStatus(c.filterStatus)
	c.updateStatusButtons()
	c.EmitUpdateEvent()
}

func (c *TopPanel) OnStatusOtherClick() {
	c.filterStatus = "OTHER"
	system.Instance.SetFilterStatus(c.filterStatus)
	c.updateStatusButtons()
	c.EmitUpdateEvent()
}

func (c *TopPanel) OnStatusAllClick() {
	c.filterStatus = "ALL"
	system.Instance.SetFilterStatus(c.filterStatus)
	c.updateStatusButtons()
	c.EmitUpdateEvent()
}

func (c *TopPanel) updateAutoupdateButton() {
	btnAutoupdate, ok := c.FindWidgetByName("btnAutoupdate").(*ui.Button)
	if !ok {
		return
	}
	if c.autoupdateOn {
		btnAutoupdate.SetRole("primary")
	} else {
		btnAutoupdate.SetRole("")
	}
}

func (c *TopPanel) updateTypeButtons() {
	btnTcp, ok := c.FindWidgetByName("btnTcp").(*ui.Button)
	if ok {
		if c.filterType == "tcp" {
			btnTcp.SetRole("primary")
		} else {
			btnTcp.SetRole("")
		}
	}

	btnUdp, ok := c.FindWidgetByName("btnUdp").(*ui.Button)
	if ok {
		if c.filterType == "udp" {
			btnUdp.SetRole("primary")
		} else {
			btnUdp.SetRole("")
		}
	}

	btnAll, ok := c.FindWidgetByName("btnAll").(*ui.Button)
	if ok {
		if c.filterType == "all" {
			btnAll.SetRole("primary")
		} else {
			btnAll.SetRole("")
		}
	}
}

func (c *TopPanel) updateStatusButtons() {
	btnListen, ok := c.FindWidgetByName("btnStatusListen").(*ui.Button)
	if !ok {
		return
	}
	btnEstablished, ok := c.FindWidgetByName("btnStatusEstablished").(*ui.Button)
	if !ok {
		return
	}

	btnOther, ok := c.FindWidgetByName("btnStatusOther").(*ui.Button)
	if !ok {
		return
	}

	btnAll, ok := c.FindWidgetByName("btnStatusAll").(*ui.Button)
	if !ok {
		return
	}

	lblStatus, ok := c.FindWidgetByName("lblStatus").(*ui.Label)
	if !ok {
		return
	}

	btnListen.SetRole("")
	btnEstablished.SetRole("")
	btnOther.SetRole("")
	btnAll.SetRole("")
	lblStatus.SetText("Status")

	if c.filterType == "udp" {
		btnListen.SetEnabled(false)
		btnEstablished.SetEnabled(false)
		btnOther.SetEnabled(false)
		btnAll.SetEnabled(false)
		lblStatus.SetEnabled(false)
		return
	}

	btnListen.SetEnabled(true)
	btnEstablished.SetEnabled(true)
	btnOther.SetEnabled(true)
	btnAll.SetEnabled(true)
	lblStatus.SetEnabled(true)

	if c.filterStatus == "LISTEN" {
		btnListen.SetRole("primary")
	}
	if c.filterStatus == "ESTABLISHED" {
		btnEstablished.SetRole("primary")
	}
	if c.filterStatus == "OTHER" {
		btnOther.SetRole("primary")
	}
	if c.filterStatus == "ALL" {
		btnAll.SetRole("primary")
	}
}

func (c *TopPanel) HandleSystemEvent(event system.Event) {
}
