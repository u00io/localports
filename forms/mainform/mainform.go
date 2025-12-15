package mainform

import (
	"fmt"

	"github.com/u00io/localports/forms/bottompanel"
	"github.com/u00io/localports/forms/centerpanel"
	"github.com/u00io/localports/forms/toppanel"
	"github.com/u00io/localports/system"
	"github.com/u00io/nuiforms/ui"
)

type MainForm struct {
	ui.Widget

	topPanel    *toppanel.TopPanel
	centerPanel *centerpanel.CenterPanel
	bottomPanel *bottompanel.BottomPanel
}

func NewMainForm() *MainForm {
	system.Instance = system.NewSystem()
	system.Instance.Start()

	var c MainForm
	c.InitWidget()

	c.topPanel = toppanel.NewTopPanel()
	c.centerPanel = centerpanel.NewCenterPanel()
	c.bottomPanel = bottompanel.NewBottomPanel()

	curstomWidgets := map[string]ui.Widgeter{
		"toppanel":    c.topPanel,
		"centerpanel": c.centerPanel,
		"bottompanel": c.bottomPanel,
	}
	c.SetLayout(`
<column>
	<widget id="toppanel" />
	<widget id="centerpanel"/>
	<widget id="bottompanel" />	
</column>
	`, &c, curstomWidgets)

	c.AddTimer(50, c.timerUpdate)

	return &c
}

func (c *MainForm) HandleSystemEvent(event system.Event) {
	fmt.Println("Event:", event)
	c.topPanel.HandleSystemEvent(event)
	c.centerPanel.HandleSystemEvent(event)
	c.bottomPanel.HandleSystemEvent(event)
}

func (c *MainForm) timerUpdate() {
	systemEvents := system.Instance.GetAndClearEvents()
	if len(systemEvents) > 0 {
		for _, ev := range systemEvents {
			c.HandleSystemEvent(ev)
		}
	}
}

func Run() {
	form := ui.NewForm()
	form.SetTitle("Local Ports")
	form.SetSize(1300, 800)
	form.Panel().AddWidgetOnGrid(NewMainForm(), 0, 0)
	form.Exec()
}
