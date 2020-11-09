package icharts

import (
	"github.com/gotk3/gotk3/gtk"
)

type Chart struct {
	GtkCanvas           *gtk.DrawingArea
	minWidth, minHeight int
}

func (c *Chart) Init() {
	c.minWidth = 200
	c.minHeight = 200
	c.GtkCanvas.SetSizeRequest(c.minWidth, c.minHeight)
}

func (c *Chart) SetMinSize(w, h int) {
	c.minWidth = w
	c.minHeight = h
	c.GtkCanvas.SetSizeRequest(c.minWidth, c.minHeight)
}
