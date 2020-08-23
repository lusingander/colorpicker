package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/lusingander/colorpicker"
)

var (
	defaultColor = color.RGBA{0xff, 0x00, 0x00, 0xff}
)

func main() {
	a := app.New()
	w := a.NewWindow("color picker sample")

	w.SetContent(fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			createPickerContainer(200, colorpicker.StyleDefault),
			createPickerContainer(200, colorpicker.StyleCircle),
		),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			createPickerContainer(200, colorpicker.StyleValue),
			createPickerContainer(200, colorpicker.StyleSaturation),
		),
	))

	w.ShowAndRun()
}

func createPickerContainer(height int, style colorpicker.PickerStyle) *fyne.Container {
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
	selectColorRect := &canvas.Rectangle{FillColor: color.RGBA{0, 0, 0, 0}}
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
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2X%.2X%.2X", rgba.R, rgba.G, rgba.B)
}

func styleName(s colorpicker.PickerStyle) string {
	switch s {
	case colorpicker.StyleCircle:
		return "StyleCircle"
	case colorpicker.StyleValue:
		return "StyleValue"
	case colorpicker.StyleSaturation:
		return "StyleSaturation"
	default:
		return "StyleDefault"
	}
}
