package bottompanel

import (
	"github.com/u00io/localports/system"
	"github.com/u00io/nuiforms/ui"
)

type BottomPanel struct {
	ui.Widget
}

func NewBottomPanel() *BottomPanel {
	var c BottomPanel
	c.InitWidget()
	c.SetLayout(`
		<label text="Based on NET.U00.IO project" />
	`, &c, nil)

	c.SetElevation(5)
	return &c
}

func (c *BottomPanel) HandleSystemEvent(event system.Event) {
}
