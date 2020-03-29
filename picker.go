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
			fyne.NewContainer(colorPickerRaster, picker.Circle),
			fyne.NewContainer(huePickerRaster, picker.selectHueMarker.Line),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	return picker
}

func newCircleColorPicker(size int) *ColorPicker {
	pickerSize := fyne.NewSize(int(float64(size)*0.8/1.4), int(float64(size)*0.8/1.4))
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
		picker.selectCircleHueMarker.Line.Refresh()
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
					picker.selectCircleHueMarker.Line,
				),
				fyne.NewContainer(
					colorPickerRaster,
					picker.Circle,
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
	*canvas.Line
}

func newSelectHueMarker(w int) *selectHueMarker {
	return &selectHueMarker{
		Line: &canvas.Line{
			Position1:   fyne.NewPos(0, 0),
			Position2:   fyne.NewPos(w, 0),
			StrokeColor: color.RGBA{50, 50, 50, 255},
			StrokeWidth: 1,
		},
	}
}

func (m *selectHueMarker) setHueMarkerPosition(h int) {
	m.Position1.Y = h
	m.Position2.Y = h
}

type selectCircleHueMarker struct {
	*canvas.Line
	cx, cy float64
}

func newSelectCircleHueMarker(w, h int) *selectCircleHueMarker {
	return &selectCircleHueMarker{
		Line: &canvas.Line{
			Position1:   fyne.NewPos(w-(w/10), h/2),
			Position2:   fyne.NewPos(w, h/2),
			StrokeColor: color.RGBA{50, 50, 50, 255},
			StrokeWidth: 1,
		},
		cx: float64(w) / 2,
		cy: float64(h) / 2,
	}
}

func (m *selectCircleHueMarker) setCircleHueMarkerPosition(pos fyne.Position) {
	v := newVectorFromPoints(m.cx, m.cy, float64(pos.X), float64(pos.Y))
	nv := v.normalize()
	v1 := nv.multiply(0.8 * m.cx)
	v2 := nv.multiply(m.cx)
	center := newVector(m.cx, m.cy)
	m.Line.Position1 = center.add(v1).toPosition()
	m.Line.Position2 = center.add(v2).toPosition()
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
