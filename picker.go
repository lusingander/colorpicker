package colorpicker

import (
	"image/color"
	"math"

	"fyne.io/fyne"
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
	// pickerSize < ((areaWidth - (hueBarWidth * 2)) / √2)
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
	return circleHuePickerFloat(float64(x), float64(y), float64(w), float64(h))
}

func circleHuePickerFloat(x, y, w, h float64) color.Color {
	ir := w/2 - w/10
	or := w / 2
	cx := w / 2
	cy := h / 2

	dist := distance(x, y, cx, cy)
	if dist < ir || or < dist {
		return color.RGBA{0, 0, 0, 0}
	}

	rad := math.Atan((x - cx) / (y - cy))
	rad += (math.Pi / 2)
	if y-cy >= 0 {
		rad += math.Pi
	}
	hue := rad / (2 * math.Pi)

	return fromHSV(hue, 1.0, 1.0)
}
