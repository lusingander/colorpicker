package colorpicker

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

type colorPickerBase struct {
	fyne.CanvasObject
	colorPickerRaster *tappableRaster
	changed           func(color.Color)
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

type defaultHueColorPicker struct {
	*colorPickerBase

	pickerWidth  int
	pickerHeight int
	hueBarWidth  int
	hue          float64
	*selectColorMarker
	*selectVerticalBarMarker
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

	colorPickerRaster := newTappableRaster(createSaturationValueColorPickerPixelColor(picker.hue))
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
		colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.setVerticalBarMarkerPosition(p.Y)
		picker.updatePickerColor()
	}
	huePickerRaster.Resize(hueSize)

	picker.selectColorMarker = newSelectColorMarker()
	picker.selectVerticalBarMarker = newSelectVerticalBarMarker(hueSize.Width)

	picker.CanvasObject = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			fyne.NewContainer(colorPickerRaster, picker.selectColorMarker.Circle),
			fyne.NewContainer(huePickerRaster, picker.selectVerticalBarMarker.Circle),
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
	p.setVerticalBarMarkerPosition(int(float64(p.pickerHeight) * h))
	p.colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(p.hue))
	p.colorPickerRaster.Refresh()
	x := int(round(float64(p.pickerWidth) * s))
	y := int(round(float64(p.pickerHeight) * (1.0 - v)))
	p.setColorMarkerPosition(fyne.NewPos(x, y))
	p.updatePickerColor()
}

type circleHueColorPicker struct {
	*colorPickerBase

	pickerWidth    int
	pickerHeight   int
	hueCircleWidth int
	hue            float64
	*selectColorMarker
	*selectCircleMarker
}

func newCircleHueColorPicker(size int) ColorPicker {
	// pickerSize < ((areaWidth - (hueBarWidth * 2)) / √2)
	pickerSize := fyne.NewSize(int(float64(size)*0.8/1.45), int(float64(size)*0.8/1.45))
	hueSize := fyne.NewSize(size, size)

	picker := &circleHueColorPicker{
		hue:            0,
		pickerWidth:    pickerSize.Width,
		pickerHeight:   pickerSize.Height,
		hueCircleWidth: hueSize.Width,
		colorPickerBase: &colorPickerBase{
			changed: func(color.Color) {},
		},
	}

	colorPickerRaster := newTappableRaster(createSaturationValueColorPickerPixelColor(picker.hue))
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
		picker.hue = picker.selectCircleMarker.calcHueFromCircleMarker(p)
		colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.setCircleMarkerPosition(p)
		picker.updatePickerColor()
	}
	circleHuePickerRaster.Resize(hueSize)

	picker.selectColorMarker = newSelectColorMarker()
	picker.selectCircleMarker = newSelectCircleMarker(hueSize.Width, hueSize.Height)

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
					picker.selectCircleMarker.Circle,
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

func (p *circleHueColorPicker) updatePickerColor() {
	x := p.selectColorMarker.center.X
	y := p.selectColorMarker.center.Y
	color := fromHSV(p.hue, float64(x)/float64(p.pickerWidth), 1.0-float64(y)/float64(p.pickerHeight))
	p.changed(color)
}

func (p *circleHueColorPicker) SetColor(c color.Color) {
	h, s, v := fromColor(c)
	p.hue = h
	p.setCircleMarekerPositionFromHue(p.hue)
	p.colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(p.hue))
	p.colorPickerRaster.Refresh()
	x := int(round(float64(p.pickerWidth) * s))
	y := int(round(float64(p.pickerHeight) * (1.0 - v)))
	p.setColorMarkerPosition(fyne.NewPos(x, y))
	p.updatePickerColor()
}

type valueColorPicker struct {
	*colorPickerBase

	pickerRadius  int
	pickerCenter  fyne.Position
	valueBarWidth int
	value         float64
	*selectColorMarker
	*selectVerticalBarMarker
}

func newValueColorPicker(size int) ColorPicker {
	pickerSize := fyne.NewSize(size, size)
	picker := &valueColorPicker{
		value:         1.,
		pickerRadius:  size / 2,
		pickerCenter:  fyne.NewPos(size/2, size/2),
		valueBarWidth: size / 10,
		colorPickerBase: &colorPickerBase{
			changed: func(color.Color) {},
		},
	}

	colorPickerRaster := newTappableRaster(createCircleHueSaturationColorPickerPixelColor(picker.value))
	colorPickerRaster.SetMinSize(pickerSize)
	colorPickerRaster.tapped = func(p fyne.Position) {
		if picker.isInPickerArea(p) {
			picker.setColorMarkerPosition(p)
			colorPickerRaster.Refresh()
		}
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...
	picker.colorPickerRaster = colorPickerRaster

	picker.selectColorMarker = newSelectColorMarker()
	picker.setColorMarkerPosition(picker.pickerCenter)

	picker.CanvasObject = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			fyne.NewContainer(colorPickerRaster, picker.selectColorMarker.Circle),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	return picker
}

func (p *valueColorPicker) SetColor(c color.Color) {
	// TODO: impl
}

func (p *valueColorPicker) isInPickerArea(pos fyne.Position) bool {
	d := distance(float64(pos.X), float64(pos.Y), float64(p.pickerCenter.X), float64(p.pickerCenter.Y))
	return d <= float64(p.pickerRadius)
}

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

func createSaturationValueColorPickerPixelColor(hue float64) func(int, int, int, int) color.Color {
	return func(x, y, w, h int) color.Color {
		return fromHSV(hue, float64(x)/float64(w), 1.0-float64(y)/float64(h))
	}
}

func hueBarPicker(x, y, w, h int) color.Color {
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

func createCircleHueSaturationColorPickerPixelColor(value float64) func(int, int, int, int) color.Color {
	return func(x, y, w, h int) color.Color {
		fx := float64(x)
		fy := float64(y)
		fw := float64(w)
		fh := float64(h)
		cx := fw / 2
		cy := fh / 2

		dist := distance(fx, fy, cx, cy)
		if cx < dist {
			return color.RGBA{0, 0, 0, 0}
		}

		rad := math.Atan((fx - cx) / (fy - cy))
		rad += (math.Pi / 2)
		if fy-cy >= 0 {
			rad += math.Pi
		}
		hue := rad / (2 * math.Pi)

		return fromHSV(hue, dist/cx, value)
	}
}
