package colorpicker

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

var (
	markerFillColor   = color.RGBA{50, 50, 50, 120}
	markerStrokeColor = color.RGBA{50, 50, 50, 200}
)

type marker struct {
	*canvas.Circle
	center fyne.Position
	radius float32
}

func newMarker(radius float32, strokeWidth int) *marker {
	marker := &marker{
		Circle: &canvas.Circle{
			FillColor:   markerFillColor,
			StrokeColor: markerStrokeColor,
			StrokeWidth: float32(strokeWidth),
		},
		radius: radius,
	}
	marker.setPosition(fyne.NewPos(0, 0))
	return marker
}

func (m *marker) setPosition(p fyne.Position) {
	m.center = p
	m.Position1 = fyne.NewPos(p.X-float32(m.radius), p.Y-float32(m.radius))
	m.Position2 = fyne.NewPos(p.X+float32(m.radius), p.Y+float32(m.radius))
}

func (m *marker) setPositionY(y float32) {
	m.setPosition(fyne.NewPos(m.center.X, y))
}

type circleBarMarker struct {
	*marker
	cx, cy float32
}

func newCircleBarMarker(w, h float32, hueBarWidth float32) *circleBarMarker {
	fw := float64(w)
	fh := float64(h)
	fr := hueBarWidth / 2
	marker := &circleBarMarker{
		marker: newMarker(fr, 2),
		cx:     w / 2,
		cy:     h / 2,
	}
	markerCenter := fyne.NewPos(float32(math.Round(fw-float64(fr))), float32(math.Round(fh/2)))
	marker.marker.setPosition(markerCenter)
	return marker
}

func (m *circleBarMarker) setCircleMarkerPosition(p fyne.Position) {
	v := newVectorFromPoints(float64(m.cx), float64(m.cy), float64(p.X), float64(p.Y))
	nv := v.normalize()
	center := newVector(float64(m.cx), float64(m.cy))
	markerCenter := center.add(nv.multiply(float64(m.cx - m.radius))).toPosition()
	m.marker.setPosition(markerCenter)
}

func (m *circleBarMarker) setCircleMarekerPositionFromHue(hue float32) {
	rad := float64(-2 * math.Pi * hue)
	center := newVector(float64(m.cx), float64(m.cy))
	dir := newVector(1, 0).rotate(rad).multiply(float64(m.cx - m.radius))
	markerCenter := center.add(dir).toPosition()
	m.marker.setPosition(markerCenter)
}

func (m *circleBarMarker) calcHueFromCircleMarker(p fyne.Position) float32 {
	v := newVectorFromPoints(float64(m.cx), float64(m.cy), float64(p.X), float64(p.Y))
	baseV := newVector(1, 0)
	rad := math.Acos(baseV.dot(v) / (v.norm() * baseV.norm()))
	if float64(p.Y-m.cy) >= 0 {
		rad = math.Pi*2 - rad
	}
	hue := rad / (math.Pi * 2)
	return float32(hue)
}
