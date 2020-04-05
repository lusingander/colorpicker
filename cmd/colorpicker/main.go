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

	w.SetContent(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		createPickerContainer(200, colorpicker.StyleDefault),
		createPickerContainer(200, colorpicker.StyleCircle),
	))

	w.ShowAndRun()
}

func createPickerContainer(height int, style colorpicker.PickerStyle) *fyne.Container {
	selectColorCode := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	selectColorRect := &canvas.Rectangle{FillColor: color.RGBA{0, 0, 0, 0}}
	selectColorRect.SetMinSize(fyne.NewSize(30, 20))

	// Create picker
	picker := colorpicker.NewColorPicker(height, style)
	picker.SetOnChanged(func(c color.Color) {
		selectColorCode.SetText(hexColorString(c))
		selectColorRect.FillColor = c
		selectColorRect.Refresh()
	})

	return fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		picker, // layout
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			selectColorCode,
			selectColorRect,
			layout.NewSpacer(),
		),
		widget.NewLabel(styleName(style)),
	)
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
