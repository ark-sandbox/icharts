package icharts

import (
	"fmt"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"image/color"
	"log"
	"math"
)

type PieSection struct {
	Val        float64
	Desc       string
	angleStart float64
	angleEnd   float64
	Color      color.Color
}

func PieSectionNew(val float64, desc string) PieSection {
	sec := PieSection{}
	sec.Val = val
	sec.Desc = desc
	return sec
}

type PieChart struct {
	Chart
	PieChartState

	Sections  []PieSection
	RandColor bool
}

type PieChartState struct {
	cursorPos         Point
	sectionNearCursor int
	radius            float64
	center            Point
}

func PieChartNew() *PieChart {
	chart := PieChart{}
	chart.GtkCanvas, _ = gtk.DrawingAreaNew()

	chart.Init()
	chart.sectionNearCursor = -1
	chart.RandColor = true

	//Call Draw method by connecting to draw event.
	chart.GtkCanvas.Connect("draw", chart.Draw)
	chart.GtkCanvas.Connect("motion-notify-event", chart.MotionNotify)
	chart.GtkCanvas.SetEvents(
		chart.GtkCanvas.GetEvents() | int(gdk.POINTER_MOTION_MASK))
	return &chart
}

func (chart *PieChart) Draw(da *gtk.DrawingArea, cr *cairo.Context) {
	awidth := da.GetAllocatedWidth()
	aheight := da.GetAllocatedHeight()

	numSections := len(chart.Sections)

	if numSections < 1 {
		return
	}

	pieRadius := math.Min(float64(awidth), float64(aheight))
	pieRadius = pieRadius / 2.0 * 0.8
	log.Println("pie radius", pieRadius)

	//Center of the pie.
	xc, yc := float64(awidth)/2.0, float64(aheight)/2.0
	chart.center = Point{xc, yc}
	chart.radius = pieRadius

	chart.computeSections()
	for i, sec := range chart.Sections {
		r := pieRadius
		if chart.sectionNearCursor == i {
			log.Println("radius increased")
			r = r * 1.2
		}
		if chart.RandColor {
			CairoSetColor(cr, BrightColors[i])
		} else {
			CairoSetColor(cr, sec.Color)
		}
		cr.Arc(xc, yc, r, sec.angleStart, sec.angleEnd)
		cr.LineTo(xc, yc)
		cr.FillPreserve()

		cr.SetLineWidth(4.0)
		cr.SetSourceRGB(1.0, 1.0, 1.0)
		cr.Stroke()
	}

	//Choosing font.
	cr.SelectFontFace("Sans", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(12.0)

	//Highlight point near cursor and its value.
	idx := chart.sectionNearCursor
	if idx != -1 {
		sec := chart.Sections[idx]
		r, angle := chart.radius, sec.angleStart+((sec.angleEnd-sec.angleStart)/2.0)
		pt := Cartesian(r, angle)
		pos := MoveOriginTo(chart.center, Point{0, 0}, pt)

		textToShow := fmt.Sprintf("%s at %.2f",
			chart.Sections[idx].Desc, chart.Sections[idx].Val)
		extents := cr.TextExtents(textToShow)

		cr.SetSourceRGBA(0.0, 0.0, 0.4, 0.5)
		cr.Rectangle(pos.X+5.0, pos.Y,
			extents.Width, extents.Height+5.0)
		cr.Fill()

		cr.SetSourceRGB(0.0, 0.0, 0.0)
		cr.MoveTo(pos.X+5.0, pos.Y+(extents.Height+5.0/2.0))
		cr.ShowText(textToShow)
	}

}

func (chart *PieChart) computeSections() {
	sum := 0.0
	currentAngle := 0.0

	for _, sec := range chart.Sections {
		sum = sum + sec.Val
	}
	for i, sec := range chart.Sections {
		sec.angleStart = currentAngle
		per := (sec.Val / sum) * (2 * math.Pi)
		sec.angleEnd = sec.angleStart + per
		currentAngle = sec.angleEnd

		chart.Sections[i] = sec
	}
}

func (chart *PieChart) MotionNotify(da *gtk.DrawingArea, ev *gdk.Event) {
	motionEvt := gdk.EventMotionNewFromEvent(ev)
	x, y := motionEvt.MotionVal()
	pt := Point{x, y}
	pt = MoveOriginTo(Point{0.0, 0.0}, chart.center, pt)
	r, angle := Polar(pt)
	log.Println("polar co-ordiante", r, angle)
	for i, sec := range chart.Sections {
		if r < chart.radius && angle >= sec.angleStart && angle <= sec.angleEnd {
			log.Println(chart.radius, sec.angleStart, sec.angleEnd)
			chart.sectionNearCursor = i
			chart.GtkCanvas.QueueDraw()
			return
		}
	}
	chart.sectionNearCursor = -1
}
