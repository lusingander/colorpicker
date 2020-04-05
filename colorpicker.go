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

	SetOnChanged(func(color.Color))
}

// NewColorPicker returns color picker conrainer.
func NewColorPicker(size int, style PickerStyle) ColorPicker {
	switch style {
	case StyleCircle:
		return newCircleColorPicker(size)
	default:
		return newColorPicker(size)
	}
}
