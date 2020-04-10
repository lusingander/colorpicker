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

type selectColorMarker struct {
	*canvas.Circle
	center fyne.Position
	radius int
}

func newSelectColorMarker() *selectColorMarker {
	marker := &selectColorMarker{
		Circle: &canvas.Circle{
			FillColor:   markerFillColor,
			StrokeColor: markerStrokeColor,
			StrokeWidth: 1,
		},
		radius: 5,
	}
	marker.setColorMarkerPosition(fyne.NewPos(0, 0))
	return marker
}

func (m *selectColorMarker) setColorMarkerPosition(p fyne.Position) {
	m.center = p
	m.Position1 = fyne.NewPos(p.X-m.radius, p.Y-m.radius)
	m.Position2 = fyne.NewPos(p.X+m.radius, p.Y+m.radius)
}

type selectHueMarker struct {
	*canvas.Circle
	r float64
}

func newSelectHueMarker(w int) *selectHueMarker {
	r := float64(w) / 2
	marker := &selectHueMarker{
		Circle: &canvas.Circle{
			FillColor:   markerFillColor,
			StrokeColor: markerStrokeColor,
			StrokeWidth: 2,
		},
		r: r,
	}
	markerCenter := fyne.NewPos(int(r), 0)
	marker.updateMarkerPosition(markerCenter)
	return marker
}

func (m *selectHueMarker) setHueMarkerPosition(h int) {
	m.updateMarkerPosition(fyne.NewPos(int(m.r), h))
}

func (m *selectHueMarker) updateMarkerPosition(p fyne.Position) {
	r := int(m.r)
	m.Circle.Position1 = fyne.NewPos(p.X-r, p.Y-r)
	m.Circle.Position2 = fyne.NewPos(p.X+r, p.Y+r)
}

type selectCircleHueMarker struct {
	*canvas.Circle
	cx, cy float64
	r      float64
}

func newSelectCircleHueMarker(w, h int) *selectCircleHueMarker {
	marker := &selectCircleHueMarker{
		Circle: &canvas.Circle{
			FillColor:   markerFillColor,
			StrokeColor: markerStrokeColor,
			StrokeWidth: 2,
		},
		cx: float64(w) / 2,
		cy: float64(h) / 2,
		r:  (float64(w) / 10) / 2,
	}
	markerCenter := fyne.NewPos(w-int(marker.r), h/2)
	marker.updateMarkerPosition(markerCenter)
	return marker
}

func (m *selectCircleHueMarker) setCircleHueMarkerPosition(p fyne.Position) {
	v := newVectorFromPoints(m.cx, m.cy, float64(p.X), float64(p.Y))
	nv := v.normalize()
	center := newVector(m.cx, m.cy)
	markerCenter := center.add(nv.multiply(0.9 * m.cx)).toPosition()
	m.updateMarkerPosition(markerCenter)
}

func (m *selectCircleHueMarker) setCircleHueMarekerPositionFromHue(hue float64) {
	rad := -2 * math.Pi * hue
	center := newVector(m.cx, m.cy)
	dir := newVector(1, 0).rotate(rad).multiply(m.cx - m.r)
	markerCenter := center.add(dir).toPosition()
	m.updateMarkerPosition(markerCenter)
}

func (m *selectCircleHueMarker) updateMarkerPosition(p fyne.Position) {
	r := int(m.r)
	m.Circle.Position1 = fyne.NewPos(p.X-r, p.Y-r)
	m.Circle.Position2 = fyne.NewPos(p.X+r, p.Y+r)
}

func (m *selectCircleHueMarker) calcHueFromCircleMarker(p fyne.Position) float64 {
	v := newVectorFromPoints(m.cx, m.cy, float64(p.X), float64(p.Y))
	baseV := newVector(1, 0)
	rad := math.Acos(baseV.dot(v) / (v.norm() * baseV.norm()))
	if float64(p.Y)-m.cy >= 0 {
		rad = math.Pi*2 - rad
	}
	hue := rad / (math.Pi * 2)
	return hue
}
