package colorpicker

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

type colorPicker struct {
	fyne.CanvasObject

	changed func(color.Color)

	cw, hw, h int
	hue       float64
	*selectColorMarker
	*selectHueMarker
	*selectCircleHueMarker

	colorPickerRaster *tappableRaster
}

type colorPickerBase struct {
	fyne.CanvasObject
	colorPickerRaster *tappableRaster
	changed           func(color.Color)
}

type defaultHueColorPicker struct {
	*colorPickerBase

	pickerWidth  int
	pickerHeight int
	hueBarWidth  int
	hue          float64
	*selectColorMarker
	*selectHueMarker
}

func (p *colorPicker) SetColor(c color.Color) {
	h, s, v := fromColor(c)
	p.hue = h
	if p.selectHueMarker != nil {
		p.setHueMarkerPosition(int(float64(p.h) * h))
	} else if p.selectCircleHueMarker != nil {
		p.setCircleHueMarekerPositionFromHue(p.hue)
	}
	p.colorPickerRaster.setPixelColor(createColorPickerPixelColor(p.hue))
	p.colorPickerRaster.Refresh()
	x := int(round(float64(p.cw) * s))
	y := int(round(float64(p.h) * (1.0 - v)))
	p.setColorMarkerPosition(fyne.NewPos(x, y))
	p.updatePickerColor()
}

func (p *colorPicker) SetOnChanged(f func(color.Color)) {
	p.changed = f
}

func (p *colorPicker) updatePickerColor() {
	x := p.selectColorMarker.center.X
	y := p.selectColorMarker.center.Y
	color := fromHSV(p.hue, float64(x)/float64(p.cw), 1.0-float64(y)/float64(p.h))
	p.changed(color)
}

func (p *colorPicker) CreateRenderer() fyne.WidgetRenderer {
	return &colorPickerWidgetRender{picker: p}
}

func newDefaultHueColorPicker(size int) ColorPicker {
	pickerSize := fyne.NewSize(size, size)
	hueSize := fyne.NewSize(size/10, size)

	picker := &defaultHueColorPicker{
		hue:          0,
		pickerWidth:  pickerSize.Width,
		pickerHeight: pickerSize.Height,
		hueBarWidth:  hueSize.Width,
		colorPickerBase: &colorPickerBase{
			changed: func(color.Color) {},
		},
	}

	colorPickerRaster := newTappableRaster(createSatularionValueColorPickerPixelColor(picker.hue))
	colorPickerRaster.SetMinSize(pickerSize)
	colorPickerRaster.tapped = func(p fyne.Position) {
		picker.setColorMarkerPosition(p)
		picker.updatePickerColor()
		colorPickerRaster.Refresh()
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...
	picker.colorPickerRaster = colorPickerRaster

	huePickerRaster := newTappableRaster(hueBarPicker)
	huePickerRaster.SetMinSize(hueSize)
	huePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = float64(p.Y) / float64(hueSize.Height)
		colorPickerRaster.setPixelColor(createSatularionValueColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.setHueMarkerPosition(p.Y)
		picker.updatePickerColor()
	}
	huePickerRaster.Resize(hueSize)

	picker.selectColorMarker = newSelectColorMarker()
	picker.selectHueMarker = newSelectHueMarker(hueSize.Width)

	picker.CanvasObject = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			fyne.NewContainer(colorPickerRaster, picker.selectColorMarker.Circle),
			fyne.NewContainer(huePickerRaster, picker.selectHueMarker.Circle),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	return picker
}

func (p *defaultHueColorPicker) updatePickerColor() {
	x := p.selectColorMarker.center.X
	y := p.selectColorMarker.center.Y
	color := fromHSV(p.hue, float64(x)/float64(p.pickerWidth), 1.0-float64(y)/float64(p.pickerHeight))
	p.changed(color)
}

func (p *defaultHueColorPicker) SetColor(c color.Color) {
	h, s, v := fromColor(c)
	p.hue = h
	p.setHueMarkerPosition(int(float64(p.pickerHeight) * h))
	p.colorPickerRaster.setPixelColor(createColorPickerPixelColor(p.hue))
	p.colorPickerRaster.Refresh()
	x := int(round(float64(p.pickerWidth) * s))
	y := int(round(float64(p.pickerHeight) * (1.0 - v)))
	p.setColorMarkerPosition(fyne.NewPos(x, y))
	p.updatePickerColor()
}

func (p *colorPickerBase) SetOnChanged(f func(color.Color)) {
	p.changed = f
}

func (p *colorPickerBase) CreateRenderer() fyne.WidgetRenderer {
	return &colorPickerBaseWidgetRender{picker: p}
}

func (p *colorPickerBase) Refresh() {
	p.CanvasObject.Refresh()
}

func (p *colorPickerBase) Position() fyne.Position {
	return p.CanvasObject.Position()
}

func (p *colorPickerBase) Move(pos fyne.Position) {
	p.CanvasObject.Move(pos)
}

func (p *colorPickerBase) Size() fyne.Size {
	return p.CanvasObject.Size()
}

func (p *colorPickerBase) MinSize() fyne.Size {
	return p.CanvasObject.MinSize()
}

func (p *colorPickerBase) Resize(size fyne.Size) {
	p.CanvasObject.Resize(size)
}

func (p *colorPickerBase) Show() {
	p.CanvasObject.Show()
}

func (p *colorPickerBase) Hide() {
	p.CanvasObject.Hide()
}

func (p *colorPickerBase) Visible() bool {
	return p.CanvasObject.Visible()
}

func newColorPicker(size int) ColorPicker {
	pickerSize := fyne.NewSize(size, size)
	hueSize := fyne.NewSize(size/10, size)

	picker := &colorPicker{
		hue:     0,
		changed: func(color.Color) {},
		cw:      pickerSize.Width,
		hw:      hueSize.Width,
		h:       size,
	}

	colorPickerRaster := newTappableRaster(createColorPickerPixelColor(picker.hue))
	colorPickerRaster.SetMinSize(pickerSize)
	colorPickerRaster.tapped = func(p fyne.Position) {
		picker.setColorMarkerPosition(p)
		picker.updatePickerColor()
		colorPickerRaster.Refresh()
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...
	picker.colorPickerRaster = colorPickerRaster

	huePickerRaster := newTappableRaster(huePicker)
	huePickerRaster.SetMinSize(hueSize)
	huePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = float64(p.Y) / float64(hueSize.Height)
		colorPickerRaster.setPixelColor(createColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.setHueMarkerPosition(p.Y)
		picker.updatePickerColor()
	}
	huePickerRaster.Resize(hueSize)

	picker.selectColorMarker = newSelectColorMarker()
	picker.selectHueMarker = newSelectHueMarker(hueSize.Width)

	picker.CanvasObject = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			fyne.NewContainer(colorPickerRaster, picker.selectColorMarker.Circle),
			fyne.NewContainer(huePickerRaster, picker.selectHueMarker.Circle),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	return picker
}

func newCircleColorPicker(size int) ColorPicker {
	// pickerSize < ((areaWidth - (hueBarWidth * 2)) / √2)
	pickerSize := fyne.NewSize(int(float64(size)*0.8/1.45), int(float64(size)*0.8/1.45))
	hueSize := fyne.NewSize(size, size)

	picker := &colorPicker{
		hue:     0,
		changed: func(color.Color) {},
		cw:      pickerSize.Width,
		hw:      hueSize.Width,
		h:       pickerSize.Height,
	}

	colorPickerRaster := newTappableRaster(createColorPickerPixelColor(picker.hue))
	colorPickerRaster.SetMinSize(pickerSize)
	colorPickerRaster.tapped = func(p fyne.Position) {
		picker.setColorMarkerPosition(p)
		picker.updatePickerColor()
		colorPickerRaster.Refresh()
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...
	picker.colorPickerRaster = colorPickerRaster

	circleHuePickerRaster := newTappableRaster(circleHuePicker)
	circleHuePickerRaster.SetMinSize(hueSize)
	circleHuePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = picker.selectCircleHueMarker.calcHueFromCircleMarker(p)
		colorPickerRaster.setPixelColor(createColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.setCircleHueMarkerPosition(p)
		picker.updatePickerColor()
	}
	circleHuePickerRaster.Resize(hueSize)

	picker.selectColorMarker = newSelectColorMarker()
	picker.selectCircleHueMarker = newSelectCircleHueMarker(hueSize.Width, hueSize.Height)

	picker.CanvasObject = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			fyne.NewContainerWithLayout(
				layout.NewCenterLayout(),
				fyne.NewContainer(
					circleHuePickerRaster,
					picker.selectCircleHueMarker.Circle,
				),
				fyne.NewContainer(
					colorPickerRaster,
					picker.selectColorMarker.Circle,
				),
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	return picker
}

type colorPickerWidgetRender struct {
	picker *colorPicker
}

func (r *colorPickerWidgetRender) Layout(size fyne.Size) {
	r.picker.CanvasObject.Resize(size)
}

func (r *colorPickerWidgetRender) MinSize() fyne.Size {
	return r.picker.CanvasObject.MinSize()
}

func (r *colorPickerWidgetRender) Refresh() {
	r.picker.CanvasObject.Refresh()
}

func (r *colorPickerWidgetRender) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *colorPickerWidgetRender) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.picker.CanvasObject}
}

func (r *colorPickerWidgetRender) Destroy() {}

type colorPickerBaseWidgetRender struct {
	picker *colorPickerBase
}

func (r *colorPickerBaseWidgetRender) Layout(size fyne.Size) {
	r.picker.CanvasObject.Resize(size)
}

func (r *colorPickerBaseWidgetRender) MinSize() fyne.Size {
	return r.picker.CanvasObject.MinSize()
}

func (r *colorPickerBaseWidgetRender) Refresh() {
	r.picker.CanvasObject.Refresh()
}

func (r *colorPickerBaseWidgetRender) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *colorPickerBaseWidgetRender) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.picker.CanvasObject}
}

func (r *colorPickerBaseWidgetRender) Destroy() {}

func createColorPickerPixelColor(hue float64) func(int, int, int, int) color.Color {
	return func(x, y, w, h int) color.Color {
		return fromHSV(hue, float64(x)/float64(w), 1.0-float64(y)/float64(h))
	}
}

func createSatularionValueColorPickerPixelColor(hue float64) func(int, int, int, int) color.Color {
	return func(x, y, w, h int) color.Color {
		return fromHSV(hue, float64(x)/float64(w), 1.0-float64(y)/float64(h))
	}
}

func hueBarPicker(x, y, w, h int) color.Color {
	return fromHSV(float64(y)/float64(h), 1.0, 1.0)
}

func huePicker(x, y, w, h int) color.Color {
	return fromHSV(float64(y)/float64(h), 1.0, 1.0)
}

func circleHuePicker(x, y, w, h int) color.Color {
	return circleHuePickerFloat(float64(x), float64(y), float64(w), float64(h))
}

func circleHuePickerFloat(x, y, w, h float64) color.Color {
	ir := w/2 - w/10
	or := w / 2
	cx := w / 2
	cy := h / 2

	dist := distance(x, y, cx, cy)
	if dist < ir || or < dist {
		return color.RGBA{0, 0, 0, 0}
	}

	rad := math.Atan((x - cx) / (y - cy))
	rad += (math.Pi / 2)
	if y-cy >= 0 {
		rad += math.Pi
	}
	hue := rad / (2 * math.Pi)

	return fromHSV(hue, 1.0, 1.0)
}
