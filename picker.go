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
}

func NewColorPicker(h int) *ColorPicker {
	w := h

	picker := &ColorPicker{
		hue:     0,
		Changed: func(color.Color) {},
	}

	colorPickerRaster := newTappableRaster(createColorPickerPixelColor(picker.hue))
	colorPickerRaster.SetMinSize(fyne.NewSize(w, h))
	colorPickerRaster.tapped = func(p fyne.Position) {
		color := fromHSV(picker.hue, float64(p.X)/float64(w), 1.0-float64(p.Y)/float64(h))
		picker.Changed(color)
		picker.setPosition(p)
	}
	colorPickerRaster.Resize(fyne.NewSize(w, h)) // Note: doesn't render if remove this line...

	huePickerRaster := newTappableRaster(huePicker)
	huePickerRaster.SetMinSize(fyne.NewSize(w/10, h))
	huePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = float64(p.Y) / float64(h)
		colorPickerRaster.setPixelColor(createColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
	}
	huePickerRaster.Resize(fyne.NewSize(w, h))

	picker.selectColorMarker = newSelectColorMarker()

	picker.CanvasObject = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			fyne.NewContainer(colorPickerRaster, picker.Circle),
			huePickerRaster,
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
