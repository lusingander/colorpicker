package colorpicker

import (
	"image/color"

	"fyne.io/fyne"
)

// PickerStyle represents how the picker is displayed.
type PickerStyle int

const (
	// StyleDefault is style to display saturation-value area and vertical hue bar.
	StyleDefault PickerStyle = iota
	// StyleCircle is style to display saturation-value area and circle hue bar.
	StyleCircle
	// StyleValue is style to display hue-saturation area and vertical value bar.
	StyleValue
)

// ColorPicker represents color picker component.
type ColorPicker interface {
	fyne.CanvasObject

	SetColor(color.Color)
	SetOnChanged(func(color.Color))
}

// New returns color picker container.
func New(size int, style PickerStyle) ColorPicker {
	switch style {
	case StyleCircle:
		return newCircleHueColorPicker(size)
	case StyleValue:
		return newValueColorPicker(size)
	default:
		return newDefaultHueColorPicker(size)
	}
}
