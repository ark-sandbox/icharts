package icharts

import (
	"fmt"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"math"
	"time"
)

type DataPoint struct {
	Val  float64
	Time time.Time
}

type LineChart struct {
	Chart
	LineChartState

	Data        []DataPoint
	MaxPoints   int
	PointRadius float64
}

type LineChartState struct {
	cursorPos       Point
	cursorNearPoint int
	dataPoints      []Point
}

func LineChartNew() *LineChart {
	chart := LineChart{}
	chart.GtkCanvas, _ = gtk.DrawingAreaNew()

	chart.Init()
	chart.PointRadius = 3.0
	chart.MaxPoints = 10
	chart.cursorNearPoint = -1

	//Call Draw method by connecting to draw event.
	chart.GtkCanvas.Connect("draw", chart.Draw)
	chart.GtkCanvas.Connect("motion-notify-event", chart.MotionNotify)
	chart.GtkCanvas.SetEvents(
		chart.GtkCanvas.GetEvents() | int(gdk.POINTER_MOTION_MASK))
	return &chart
}

func (chart *LineChart) Draw(da *gtk.DrawingArea, cr *cairo.Context) {

	awidth := da.GetAllocatedWidth()
	aheight := da.GetAllocatedHeight()
	/*
		if awidth > chart.MinWidth && aheight > chart.MinHeight {
			chart.GtkCanvas.SetSizeRequest(awidth, aheight)
		}*/

	//Normalise data points between 0 and 1.
	numPoints := len(chart.Data)
	if len(chart.Data) > chart.MaxPoints {
		numPoints = numPoints
	}
	plotData := Normalise(chart.Data[:numPoints])
	colWidth := awidth / chart.MaxPoints

	cr.SetSourceRGB(0.0, 0.0, 0.0)

	chart.dataPoints = nil
	for i, data := range plotData {
		cr.LineTo(float64(i*colWidth), data.Val*float64(aheight))
		chart.dataPoints = append(chart.dataPoints,
			Point{float64(i * colWidth), data.Val * float64(aheight)})
	}
	cr.Stroke()
	log.Print("datapoints : %+v", chart.dataPoints)

	for i, data := range plotData {
		cr.SetSourceRGB(1.0, 1.0, 1.0)
		cr.Arc(float64(i*colWidth), data.Val*float64(aheight),
			chart.PointRadius, 0.0, 2*math.Pi)
		cr.FillPreserve()
		cr.SetSourceRGB(0.0, 0.0, 0.8)
		cr.Stroke()
	}

	//Choosing font.
	cr.SelectFontFace("Sans", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(12.0)

	//Highlight point near cursor and its value.
	idx := chart.cursorNearPoint
	if idx != -1 {
		point := chart.dataPoints[idx]
		cr.SetSourceRGB(0.0, 0.0, 0.4)
		cr.Arc(point.X, point.Y,
			chart.PointRadius+2.0, 0.0, 2*math.Pi)
		cr.Stroke()

		textToShow := fmt.Sprintf("%.2f at %s",
			chart.Data[idx].Val, chart.Data[idx].Time.Format("Mon Jan _2 15:04:05"))
		extents := cr.TextExtents(textToShow)

		cr.SetSourceRGBA(0.0, 0.0, 0.4, 0.5)
		cr.Rectangle(point.X+5.0, point.Y,
			extents.Width, extents.Height+5.0)
		cr.Fill()

		cr.SetSourceRGB(0.0, 0.0, 0.0)
		cr.MoveTo(point.X+5.0, point.Y+(extents.Height+5.0/2.0))
		cr.ShowText(textToShow)
	}

}

func (chart *LineChart) MotionNotify(da *gtk.DrawingArea, ev *gdk.Event) {
	motionEvt := gdk.EventMotionNewFromEvent(ev)
	x, y := motionEvt.MotionVal()
	log.Printf("MotionNotify is called...%+v, %+v, %+v\n",
		x, y, motionEvt.Type())
	chart.cursorPos = Point{x, y}

	//Highlight point near cursor and its value.
	chart.cursorNearPoint = -1
	for i, point := range chart.dataPoints {
		if Dist(point, chart.cursorPos) < 15.0 {
			chart.cursorNearPoint = i
			chart.GtkCanvas.QueueDraw()
			log.Println("updated cursorNearPoint to ",
				chart.cursorNearPoint)
		}
	}
}
