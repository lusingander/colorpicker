package colorpicker

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
)

// PickerStyle represents how the picker is displayed.
type PickerStyle int

const (
	// StyleDefault is style to display vertical hue bar.
	StyleDefault PickerStyle = iota
	// StyleCircle is style to display circle hue bar.
	StyleCircle
)

// ColorPicker represents color picker component.
type ColorPicker struct {
	fyne.CanvasObject

	Changed func(color.Color)

	cw, hw, h int
	hue       float64
	*selectColorMarker
	*selectHueMarker
	*selectCircleHueMarker
}

// NewColorPicker returns color picker conrainer.
func NewColorPicker(size int, style PickerStyle) *ColorPicker {
	switch style {
	case StyleCircle:
		return newCircleColorPicker(size)
	default:
		return newColorPicker(size)
	}
}

func newColorPicker(size int) *ColorPicker {
	pickerSize := fyne.NewSize(size, size)
	hueSize := fyne.NewSize(size/10, size)

	picker := &ColorPicker{
		hue:     0,
		Changed: func(color.Color) {},
		cw:      pickerSize.Width,
		hw:      hueSize.Width,
		h:       size,
	}

	colorPickerRaster := newTappableRaster(createColorPickerPixelColor(picker.hue))
	colorPickerRaster.SetMinSize(pickerSize)
	colorPickerRaster.tapped = func(p fyne.Position) {
		picker.setColorMarkerPosition(p)
		picker.updatePickerColor()
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...

	huePickerRaster := newTappableRaster(huePicker)
	huePickerRaster.SetMinSize(hueSize)
	huePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = float64(p.Y) / float64(hueSize.Height)
		colorPickerRaster.setPixelColor(createColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.setHueMarkerPosition(p.Y)
		picker.updatePickerColor()
	}
	huePickerRaster.Resize(hueSize)

	picker.selectColorMarker = newSelectColorMarker()
	picker.selectHueMarker = newSelectHueMarker(hueSize.Width)

	picker.CanvasObject = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			fyne.NewContainer(colorPickerRaster, picker.selectColorMarker.Circle),
			fyne.NewContainer(huePickerRaster, picker.selectHueMarker.Circle),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	return picker
}

func newCircleColorPicker(size int) *ColorPicker {
	pickerSize := fyne.NewSize(int(float64(size)*0.8/1.45), int(float64(size)*0.8/1.45))
	hueSize := fyne.NewSize(size, size)

	picker := &ColorPicker{
		hue:     0,
		Changed: func(color.Color) {},
		cw:      pickerSize.Width,
		hw:      hueSize.Width,
		h:       pickerSize.Height,
	}

	colorPickerRaster := newTappableRaster(createColorPickerPixelColor(picker.hue))
	colorPickerRaster.SetMinSize(pickerSize)
	colorPickerRaster.tapped = func(p fyne.Position) {
		picker.setColorMarkerPosition(p)
		picker.updatePickerColor()
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...

	circleHuePickerRaster := newTappableRaster(circleHuePicker)
	circleHuePickerRaster.SetMinSize(hueSize)
	circleHuePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = picker.selectCircleHueMarker.calcHueFromCircleMarker(p)
		colorPickerRaster.setPixelColor(createColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.setCircleHueMarkerPosition(p)
		picker.updatePickerColor()
	}
	circleHuePickerRaster.Resize(hueSize)

	picker.selectColorMarker = newSelectColorMarker()
	picker.selectCircleHueMarker = newSelectCircleHueMarker(hueSize.Width, hueSize.Height)

	picker.CanvasObject = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			fyne.NewContainerWithLayout(
				layout.NewCenterLayout(),
				fyne.NewContainer(
					circleHuePickerRaster,
					picker.selectCircleHueMarker.Circle,
				),
				fyne.NewContainer(
					colorPickerRaster,
					picker.selectColorMarker.Circle,
				),
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	return picker
}

func (p *ColorPicker) updatePickerColor() {
	x := p.selectColorMarker.center.X
	y := p.selectColorMarker.center.Y
	color := fromHSV(p.hue, float64(x)/float64(p.cw), 1.0-float64(y)/float64(p.h))
	p.Changed(color)
}

func createColorPickerPixelColor(hue float64) func(int, int, int, int) color.Color {
	return func(x, y, w, h int) color.Color {
		return fromHSV(hue, float64(x)/float64(w), 1.0-float64(y)/float64(h))
	}
}

func huePicker(x, y, w, h int) color.Color {
	return fromHSV(float64(y)/float64(h), 1.0, 1.0)
}

func circleHuePicker(x, y, w, h int) color.Color {
	fx := float64(x)
	fy := float64(y)
	fw := float64(w)
	fh := float64(h)

	ir := fw/2 - fw/10
	or := fw / 2
	cx := fw / 2
	cy := fh / 2

	dist := math.Sqrt(math.Pow(fx-cx, 2) + math.Pow(fy-cy, 2))
	if dist < ir || or < dist {
		return color.RGBA{0, 0, 0, 0}
	}

	rad := math.Atan((fx - cx) / (fy - cy))
	rad += (math.Pi / 2)
	if fy-cy >= 0 {
		rad += math.Pi
	}
	rad /= 2 * math.Pi

	return fromHSV(rad, 1.0, 1.0)
}

type selectColorMarker struct {
	*canvas.Circle
	center fyne.Position
	r      int
}

func newSelectColorMarker() *selectColorMarker {
	p := fyne.NewPos(0, 0)
	r := 5
	return &selectColorMarker{
		Circle: &canvas.Circle{
			Position1:   fyne.NewPos(p.X-r, p.Y-r),
			Position2:   fyne.NewPos(p.X+r, p.Y+r),
			StrokeColor: color.RGBA{50, 50, 50, 255},
			StrokeWidth: 1,
		},
		center: p,
		r:      r,
	}
}

func (m *selectColorMarker) setColorMarkerPosition(pos fyne.Position) {
	m.center = pos
	m.Move(fyne.NewPos(pos.X-m.r, pos.Y-m.r))
}

type selectHueMarker struct {
	*canvas.Circle
	r float64
}

func newSelectHueMarker(w int) *selectHueMarker {
	r := float64(w) / 2
	marker := &selectHueMarker{
		Circle: &canvas.Circle{
			FillColor:   color.RGBA{50, 50, 50, 120},
			StrokeColor: color.RGBA{50, 50, 50, 200},
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
			FillColor:   color.RGBA{50, 50, 50, 120},
			StrokeColor: color.RGBA{50, 50, 50, 200},
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

func (m *selectCircleHueMarker) setCircleHueMarkerPosition(pos fyne.Position) {
	v := newVectorFromPoints(m.cx, m.cy, float64(pos.X), float64(pos.Y))
	nv := v.normalize()
	center := newVector(m.cx, m.cy)
	markerCenter := center.add(nv.multiply(0.9 * m.cx)).toPosition()
	m.updateMarkerPosition(markerCenter)
}

func (m *selectCircleHueMarker) updateMarkerPosition(p fyne.Position) {
	r := int(m.r)
	m.Circle.Position1 = fyne.NewPos(p.X-r, p.Y-r)
	m.Circle.Position2 = fyne.NewPos(p.X+r, p.Y+r)
}

func (m *selectCircleHueMarker) calcHueFromCircleMarker(pos fyne.Position) float64 {
	v := newVectorFromPoints(m.cx, m.cy, float64(pos.X), float64(pos.Y))
	baseV := newVector(1, 0)
	rad := math.Acos(baseV.dot(v) / (v.norm() * baseV.norm()))
	if float64(pos.Y)-m.cy >= 0 {
		rad = math.Pi*2 - rad
	}
	rad /= (math.Pi * 2)
	return rad
}
