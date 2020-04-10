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

func main() {
	a := app.New()
	w := a.NewWindow("color picker sample")
	w.Resize(fyne.NewSize(400, 400))
	w.SetContent(createContainer(w))
	w.ShowAndRun()
}

func createContainer(w fyne.Window) fyne.CanvasObject {
	var current color.Color

	selectColorCode := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	selectColorRect := &canvas.Rectangle{FillColor: current}
	selectColorRect.SetMinSize(fyne.NewSize(30, 20))

	picker := colorpicker.New(200, colorpicker.StyleDefault)
	picker.SetOnChanged(func(c color.Color) {
		current = c
		selectColorCode.SetText(hexColorString(current))
		selectColorRect.FillColor = current
		selectColorRect.Refresh()
	})
	content := fyne.NewContainer(picker)

	button := widget.NewButton("Open color picker", func() {
		picker.SetColor(current)
		dialog.ShowCustom("Select color", "OK", content, w)
	})

	current = color.RGBA{0xff, 0x00, 0x00, 0xff}
	selectColorCode.SetText(hexColorString(current))
	selectColorRect.FillColor = current
	selectColorRect.Refresh()

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
				selectColorCode,
				selectColorRect,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
}

func hexColorString(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2X%.2X%.2X", rgba.R, rgba.G, rgba.B)
}
