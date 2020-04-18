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

func main() {
	a := app.New()
	w := a.NewWindow("color picker sample")

	w.SetContent(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		createPickerContainer(200, colorpicker.StyleDefault),
		createPickerContainer(200, colorpicker.StyleCircle),
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
	if s == colorpicker.StyleCircle {
		return "StyleCircle"
	}
	return "StyleDefault"
}
