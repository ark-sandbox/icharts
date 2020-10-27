package main

import (
	. "../icharts"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"time"
)

func main() {
	gtk.Init(nil)

	// gui boilerplate
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	linechart1 := LineChartNew()
	datapoints1 := []DataPoint{{20.0, time.Now()},
		{30.0, time.Now().Add(-1 * time.Minute)},
		{40.0, time.Now().Add(-2 * time.Minute)},
		{50.0, time.Now().Add(-3 * time.Minute)},
		{60.0, time.Now().Add(-4 * time.Minute)},
		{20.0, time.Now().Add(-5 * time.Minute)},
		{80.0, time.Now().Add(-6 * time.Minute)},
	}
	linechart1.Data = datapoints1
	linechart2 := LineChartNew()
	linechart1.GtkCanvas.SetHExpand(true)
	linechart1.GtkCanvas.SetVExpand(true)

	linechart2.GtkCanvas.SetHExpand(true)
	linechart2.GtkCanvas.SetVExpand(true)
	grid.Add(linechart2.GtkCanvas)
	grid.Add(linechart1.GtkCanvas)

	topLabel, err := gtk.LabelNew("Grid End")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	_ = topLabel
	//grid.Add(topLabel)
	win.Add(grid)
	win.SetTitle("Arrow keys")
	win.Connect("destroy", gtk.MainQuit)
	win.ShowAll()

	gtk.Main()
}
