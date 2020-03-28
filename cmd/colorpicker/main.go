package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	colorpicker "github.com/lusingander/fyne-colorpicker"
)

func main() {
	a := app.New()
	w := a.NewWindow("color picker sample")

	picker := colorpicker.NewColorPicker(200)

	selectColorCode := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	selectColorRect := &canvas.Rectangle{FillColor: color.RGBA{0, 0, 0, 0}}
	selectColorRect.SetMinSize(fyne.NewSize(30, 20))
	picker.Changed = func(c color.Color) {
		selectColorCode.SetText(hexColorString(c))
		selectColorRect.FillColor = c
		selectColorRect.Refresh()
	}

	w.SetContent(fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		picker.CanvasObject,
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			selectColorCode,
			selectColorRect,
			layout.NewSpacer(),
		),
	))

	w.ShowAndRun()
}

func hexColorString(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2X%.2X%.2X", rgba.R, rgba.G, rgba.B)
}
