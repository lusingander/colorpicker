package colorpicker

import (
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type tappableRaster struct {
	widget.BaseWidget

	r   *canvas.Raster
	img draw.Image

	tapped func(fyne.Position)
}

func newTappableRaster(pixelColor func(x, y, w, h int) color.Color) *tappableRaster {
	r := &tappableRaster{
		r: &canvas.Raster{},
	}
	r.setPixelColor(pixelColor)
	r.ExtendBaseWidget(r)
	return r
}

func (r *tappableRaster) setPixelColor(pixelColor func(x, y, w, h int) color.Color) {
	r.r.Generator = func(w, h int) image.Image {
		if r.img == nil || r.img.Bounds().Size().X != w || r.img.Bounds().Size().Y != h {
			rect := image.Rect(0, 0, w, h)
			r.img = image.NewNRGBA(rect)
		}
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r.img.Set(x, y, pixelColor(x, y, w, h))
			}
		}
		return r.img
	}
}

func (r *tappableRaster) CreateRenderer() fyne.WidgetRenderer {
	return &rasterWidgetRender{raster: r}
}

func (r *tappableRaster) SetMinSize(size fyne.Size) {
	r.r.SetMinSize(size)
}

func (r *tappableRaster) MinSize() fyne.Size {
	return r.r.MinSize()
}

func (r *tappableRaster) Tapped(e *fyne.PointEvent) {
	if r.tapped != nil && r.isOnRaster(e.Position) {
		r.tapped(e.Position)
	}
}

func (r *tappableRaster) TappedSecondary(*fyne.PointEvent) {}

func (r *tappableRaster) Dragged(e *fyne.DragEvent) {
	if r.tapped != nil && r.isOnRaster(e.Position) {
		r.tapped(e.Position)
	}
}

func (r *tappableRaster) DragEnd() {}

func (r *tappableRaster) Cursor() desktop.Cursor {
	return desktop.CrosshairCursor
}

func (r *tappableRaster) isOnRaster(p fyne.Position) bool {
	return 0 <= p.X && p.X <= r.Size().Width && 0 <= p.Y && p.Y <= r.Size().Height
}

type rasterWidgetRender struct {
	raster *tappableRaster
}

func (r *rasterWidgetRender) Layout(size fyne.Size) {
	r.raster.r.Resize(size)
}

func (r *rasterWidgetRender) MinSize() fyne.Size {
	return r.raster.r.MinSize()
}

func (r *rasterWidgetRender) Refresh() {
	canvas.Refresh(r.raster.r)
}

func (r *rasterWidgetRender) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *rasterWidgetRender) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.raster.r}
}

func (r *rasterWidgetRender) Destroy() {}
