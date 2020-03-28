package colorpicker

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
)

type ColorPicker struct {
	fyne.CanvasObject

	Changed func(color.Color)

	hue float64
	*selectColorMarker
	*selectHueMarker
}

func NewColorPicker(h int) *ColorPicker {
	w := h
	hw := w / 10

	picker := &ColorPicker{
		hue:     0,
		Changed: func(color.Color) {},
	}

	colorPickerRaster := newTappableRaster(createColorPickerPixelColor(picker.hue))
	colorPickerRaster.SetMinSize(fyne.NewSize(w, h))
	colorPickerRaster.tapped = func(p fyne.Position) {
		color := fromHSV(picker.hue, float64(p.X)/float64(w), 1.0-float64(p.Y)/float64(h))
		picker.Changed(color)
		picker.selectColorMarker.setPosition(p)
	}
	colorPickerRaster.Resize(fyne.NewSize(w, h)) // Note: doesn't render if remove this line...

	huePickerRaster := newTappableRaster(huePicker)
	huePickerRaster.SetMinSize(fyne.NewSize(hw, h))
	huePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = float64(p.Y) / float64(h)
		colorPickerRaster.setPixelColor(createColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.selectHueMarker.setPosition(p.Y)
	}
	huePickerRaster.Resize(fyne.NewSize(hw, h))

	picker.selectColorMarker = newSelectColorMarker()
	picker.selectHueMarker = newSelectHueMarker(hw)

	picker.CanvasObject = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			fyne.NewContainer(colorPickerRaster, picker.Circle),
			fyne.NewContainer(huePickerRaster, picker.Line),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	return picker
}

func createColorPickerPixelColor(hue float64) func(int, int, int, int) color.Color {
	return func(x, y, w, h int) color.Color {
		return fromHSV(hue, float64(x)/float64(w), 1.0-float64(y)/float64(h))
	}
}

func huePicker(x, y, w, h int) color.Color {
	return fromHSV(float64(y)/float64(h), 1.0, 1.0)
}

type selectColorMarker struct {
	*canvas.Circle
	r int
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
		r: r,
	}
}

func (m *selectColorMarker) setPosition(pos fyne.Position) {
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

func (m *selectHueMarker) setPosition(h int) {
	m.Position1.Y = h
	m.Position2.Y = h
}
