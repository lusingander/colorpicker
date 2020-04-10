package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
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
	w.Resize(fyne.NewSize(400, 400))
	w.SetContent(createContainer(w))
	w.ShowAndRun()
}

func createContainer(w fyne.Window) fyne.CanvasObject {
	var current color.Color

	displayColor := newDisplayColor()

	picker := colorpicker.New(200, colorpicker.StyleDefault)
	picker.SetOnChanged(func(c color.Color) {
		current = c
		displayColor.setColor(current)
	})
	content := fyne.NewContainer(picker)

	button := widget.NewButton("Open color picker", func() {
		picker.SetColor(current)
		dialog.ShowCustom("Select color", "OK", content, w)
	})

	current = defaultColor
	displayColor.setColor(current)

	return fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewVBoxLayout(),
			layout.NewSpacer(),
			button,
			layout.NewSpacer(),
			fyne.NewContainerWithLayout(
				layout.NewHBoxLayout(),
				layout.NewSpacer(),
				displayColor.label,
				displayColor.rect,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
}

type displayColor struct {
	label *widget.Label
	rect  *canvas.Rectangle
}

func newDisplayColor() *displayColor {
	selectColorCode := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	selectColorRect := &canvas.Rectangle{}
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
