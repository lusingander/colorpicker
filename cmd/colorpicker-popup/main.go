package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
	current = defaultColor

	simpleDisplayColor := newSimpleDisplayColor()
	picker := colorpicker.New(200, colorpicker.StyleDefault)
	picker.SetOnChanged(func(c color.Color) {
		current = c
		simpleDisplayColor.setColor(current)
	})
	content := fyne.NewContainer(picker)
	button := widget.NewButton("Open color picker", func() {
		picker.SetColor(current)
		dialog.ShowCustom("Select color", "OK", content, w)
	})
	simpleDisplayColor.setColor(current)

	tappableDisplayColor := newTappableDisplayColor(w)
	tappableDisplayColor.setColor(current)

	return fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewVBoxLayout(),
			layout.NewSpacer(),
			button,
			fyne.NewContainerWithLayout(
				layout.NewHBoxLayout(),
				layout.NewSpacer(),
				simpleDisplayColor.label,
				simpleDisplayColor.rect,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
			widget.NewLabel("Or tap rectangle"),
			fyne.NewContainerWithLayout(
				layout.NewHBoxLayout(),
				layout.NewSpacer(),
				tappableDisplayColor.label,
				tappableDisplayColor.rect,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
}

type simpleDisplayColor struct {
	label *widget.Label
	rect  *canvas.Rectangle
}

func newSimpleDisplayColor() *simpleDisplayColor {
	selectColorCode := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	selectColorRect := &canvas.Rectangle{}
	selectColorRect.SetMinSize(fyne.NewSize(30, 20))
	return &simpleDisplayColor{
		label: selectColorCode,
		rect:  selectColorRect,
	}
}

func (c *simpleDisplayColor) setColor(clr color.Color) {
	c.label.SetText(hexColorString(clr))
	c.rect.FillColor = clr
	c.rect.Refresh()
}

type tappableDisplayColor struct {
	label *widget.Label
	rect  colorpicker.PickerOpenWidget
}

func newTappableDisplayColor(w fyne.Window) *tappableDisplayColor {
	selectColorCode := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	selectColorRect := colorpicker.NewColorSelectModalRect(w, fyne.NewSize(30, 20), defaultColor)
	d := &tappableDisplayColor{
		label: selectColorCode,
		rect:  selectColorRect,
	}
	selectColorRect.SetOnChange(d.setColor)
	return d
}

func (c *tappableDisplayColor) setColor(clr color.Color) {
	c.label.SetText(hexColorString(clr))
}

func hexColorString(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2X%.2X%.2X", rgba.R, rgba.G, rgba.B)
}
