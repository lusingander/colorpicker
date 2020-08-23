package colorpicker

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

var (
	markerFillColor   = color.RGBA{50, 50, 50, 120}
	markerStrokeColor = color.RGBA{50, 50, 50, 200}
)

type marker struct {
	*canvas.Circle
	center fyne.Position
	radius int
}

func newSelectColorMarker() *marker {
	marker := &marker{
		Circle: &canvas.Circle{
			FillColor:   markerFillColor,
			StrokeColor: markerStrokeColor,
			StrokeWidth: 1,
		},
		radius: 5,
	}
	marker.setPosition(fyne.NewPos(0, 0))
	return marker
}

func (m *marker) setPosition(p fyne.Position) {
	m.center = p
	m.Position1 = fyne.NewPos(p.X-m.radius, p.Y-m.radius)
	m.Position2 = fyne.NewPos(p.X+m.radius, p.Y+m.radius)
}

type selectVerticalBarMarker struct {
	*canvas.Circle
	radius float64
}

func newSelectVerticalBarMarker(w int) *selectVerticalBarMarker {
	marker := &selectVerticalBarMarker{
		Circle: &canvas.Circle{
			FillColor:   markerFillColor,
			StrokeColor: markerStrokeColor,
			StrokeWidth: 2,
		},
		radius: float64(w) / 2,
	}
	marker.setVerticalBarMarkerPosition(0)
	return marker
}

func (m *selectVerticalBarMarker) setVerticalBarMarkerPosition(h int) {
	m.updateMarkerPosition(fyne.NewPos(int(m.radius), h))
}

func (m *selectVerticalBarMarker) updateMarkerPosition(p fyne.Position) {
	r := int(round(m.radius))
	m.Circle.Position1 = fyne.NewPos(p.X-r, p.Y-r)
	m.Circle.Position2 = fyne.NewPos(p.X+r, p.Y+r)
}

type selectCircleMarker struct {
	*canvas.Circle
	cx, cy float64
	radius float64
}

func newSelectCircleMarker(w, h int) *selectCircleMarker {
	fw := float64(w)
	fh := float64(h)
	marker := &selectCircleMarker{
		Circle: &canvas.Circle{
			FillColor:   markerFillColor,
			StrokeColor: markerStrokeColor,
			StrokeWidth: 2,
		},
		cx:     fw / 2,
		cy:     fh / 2,
		radius: (fw / 10) / 2,
	}
	markerCenter := fyne.NewPos(int(round(fw-marker.radius)), int(round(fh/2)))
	marker.updateMarkerPosition(markerCenter)
	return marker
}

func (m *selectCircleMarker) setCircleMarkerPosition(p fyne.Position) {
	v := newVectorFromPoints(m.cx, m.cy, float64(p.X), float64(p.Y))
	nv := v.normalize()
	center := newVector(m.cx, m.cy)
	markerCenter := center.add(nv.multiply(m.cx - m.radius)).toPosition()
	m.updateMarkerPosition(markerCenter)
}

func (m *selectCircleMarker) setCircleMarekerPositionFromHue(hue float64) {
	rad := -2 * math.Pi * hue
	center := newVector(m.cx, m.cy)
	dir := newVector(1, 0).rotate(rad).multiply(m.cx - m.radius)
	markerCenter := center.add(dir).toPosition()
	m.updateMarkerPosition(markerCenter)
}

func (m *selectCircleMarker) updateMarkerPosition(p fyne.Position) {
	r := int(round(m.radius))
	m.Circle.Position1 = fyne.NewPos(p.X-r, p.Y-r)
	m.Circle.Position2 = fyne.NewPos(p.X+r, p.Y+r)
}

func (m *selectCircleMarker) calcHueFromCircleMarker(p fyne.Position) float64 {
	v := newVectorFromPoints(m.cx, m.cy, float64(p.X), float64(p.Y))
	baseV := newVector(1, 0)
	rad := math.Acos(baseV.dot(v) / (v.norm() * baseV.norm()))
	if float64(p.Y)-m.cy >= 0 {
		rad = math.Pi*2 - rad
	}
	hue := rad / (math.Pi * 2)
	return hue
}
