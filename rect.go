package colorpicker

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	colorSelectModalPickerDefaultSize = 200
)

// PickerOpenWidget represents a widget that can open a color picker.
type PickerOpenWidget interface {
	fyne.CanvasObject

	SetColor(color.Color)
	SetOnChange(f func(color.Color))
	SetPickerStyle(s PickerStyle)
}

type colorSelectModalRect struct {
	*tappableRect
	parent      fyne.Window
	onChange    func(color.Color)
	pickerStyle PickerStyle
}

// NewColorSelectModalRect returns a rectangle that can be tapped to open a color picker modal.
func NewColorSelectModalRect(parent fyne.Window, minSize fyne.Size, defalutColor color.Color) PickerOpenWidget {
	rect := &colorSelectModalRect{
		tappableRect: newTappableRect(defalutColor),
		parent:       parent,
		pickerStyle:  StyleHue,
	}
	rect.tappableRect.tapped = rect.tapped
	rect.tappableRect.SetMinSize(minSize)
	return rect
}

func (r *colorSelectModalRect) SetOnChange(f func(color.Color)) {
	r.onChange = f
}

func (r *colorSelectModalRect) SetPickerStyle(s PickerStyle) {
	r.pickerStyle = s
}

func (r *colorSelectModalRect) SetColor(c color.Color) {
	r.setColor(c)
}

func (r *colorSelectModalRect) tapped(e *fyne.PointEvent) {
	picker := New(colorSelectModalPickerDefaultSize, r.pickerStyle)
	picker.SetColor(r.color())
	picker.SetOnChanged(func(c color.Color) {
		if r.onChange != nil {
			r.onChange(c)
		}
		r.setColor(c)
	})

	dialog.ShowCustom("Select color", "OK", fyne.NewContainer(picker), r.parent)
}

func (r *colorSelectModalRect) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type tappableRect struct {
	widget.BaseWidget
	rect   *canvas.Rectangle
	tapped func(*fyne.PointEvent)
}

func newTappableRect(fillColor color.Color) *tappableRect {
	r := &tappableRect{
		rect: &canvas.Rectangle{
			StrokeColor: color.NRGBA{255, 255, 255, 255},
			StrokeWidth: 1,
			FillColor:   fillColor,
		},
	}
	r.ExtendBaseWidget(r)
	return r
}

func (r *tappableRect) color() color.Color {
	return r.rect.FillColor
}

func (r *tappableRect) setColor(c color.Color) {
	r.rect.FillColor = c
	r.Refresh()
}

func (r *tappableRect) CreateRenderer() fyne.WidgetRenderer {
	return &tappableRectRenderer{rect: r.rect}
}

func (r *tappableRect) SetMinSize(size fyne.Size) {
	r.rect.SetMinSize(size)
}

func (r *tappableRect) MinSize() fyne.Size {
	return r.rect.MinSize()
}

func (r *tappableRect) Tapped(e *fyne.PointEvent) {
	if r.tapped != nil {
		r.tapped(e)
	}
}

func (r *tappableRect) TappedSecondary(*fyne.PointEvent) {}

type tappableRectRenderer struct {
	rect *canvas.Rectangle
}

func (r *tappableRectRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
}

func (r *tappableRectRenderer) MinSize() fyne.Size {
	return r.rect.MinSize()
}

func (r *tappableRectRenderer) Refresh() {
	canvas.Refresh(r.rect)
}

func (r *tappableRectRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *tappableRectRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rect}
}

func (r *tappableRectRenderer) Destroy() {}
