package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/lusingander/colorpicker"
)

var (
	defaultColor = color.NRGBA{0xff, 0x00, 0x00, 0xff}
)

func main() {
	a := app.New()
	w := a.NewWindow("color picker sample")

	w.SetContent(fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			createPickerContainer(200, colorpicker.StyleHue),
			createPickerContainer(200, colorpicker.StyleHueCircle),
		),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			createPickerContainer(200, colorpicker.StyleValue),
			createPickerContainer(200, colorpicker.StyleSaturation),
		),
	))

	w.ShowAndRun()
}

func createPickerContainer(height float32, style colorpicker.PickerStyle) *fyne.Container {
	displayColor := newDisplayColor()

	// Create picker
	picker := colorpicker.New(height, style)
	picker.SetOnChanged(func(c color.Color) {
		displayColor.setColor(c)
	})
	picker.SetColor(defaultColor)

	return fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		picker, // layout
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			displayColor.label,
			displayColor.rect,
			layout.NewSpacer(),
		),
		widget.NewLabel(styleName(style)),
	)
}

type displayColor struct {
	label *widget.Label
	rect  *canvas.Rectangle
}

func newDisplayColor() *displayColor {
	selectColorCode := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	selectColorRect := &canvas.Rectangle{FillColor: color.NRGBA{0, 0, 0, 0}}
	selectColorRect.SetMinSize(fyne.NewSize(30, 20))
	return &displayColor{
		label: selectColorCode,
		rect:  selectColorRect,
	}
}

func (c *displayColor) setColor(clr color.Color) {
	c.label.SetText(hexColorString(clr))
	c.rect.FillColor = clr
	c.rect.Refresh()
}

func hexColorString(c color.Color) string {
	rgba, _ := c.(color.NRGBA)
	return fmt.Sprintf("#%.2X%.2X%.2X%.2X", rgba.R, rgba.G, rgba.B, rgba.A)
}

func styleName(s colorpicker.PickerStyle) string {
	switch s {
	case colorpicker.StyleHueCircle:
		return "StyleHueCircle"
	case colorpicker.StyleValue:
		return "StyleValue"
	case colorpicker.StyleSaturation:
		return "StyleSaturation"
	default:
		return "StyleHue"
	}
}
