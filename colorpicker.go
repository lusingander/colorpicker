package colorpicker

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// PickerStyle represents how the picker is displayed.
type PickerStyle int

const (
	// StyleHue is style to display saturation-value area and vertical hue bar.
	StyleHue PickerStyle = iota
	// StyleHueCircle is style to display saturation-value area and circle hue bar.
	StyleHueCircle
	// StyleValue is style to display hue-saturation area and vertical value bar.
	StyleValue
	// StyleSaturation is style to display hue-value area and vertical saturation bar.
	StyleSaturation
)

// ColorPicker represents color picker component.
type ColorPicker interface {
	fyne.CanvasObject

	SetColor(color.Color)
	SetOnChanged(func(color.Color))
}

// New returns color picker container.
func New(size float32, style PickerStyle) ColorPicker {
	switch style {
	case StyleHueCircle:
		return newCircleHueColorPicker(size)
	case StyleValue:
		return newValueColorPicker(size)
	case StyleSaturation:
		return newSaturationColorPicker(size)
	default:
		return newDefaultHueColorPicker(size)
	}
}
