package colorpicker

import (
	"image/color"

	"fyne.io/fyne"
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
type ColorPicker interface {
	fyne.CanvasObject

	SetColor(color.Color)
	SetOnChanged(func(color.Color))
}

// New returns color picker conrainer.
func New(size int, style PickerStyle) ColorPicker {
	switch style {
	case StyleCircle:
		return newCircleHueColorPicker(size)
	default:
		return newDefaultHueColorPicker(size)
	}
}
