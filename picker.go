package colorpicker

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

var (
	transparent = color.NRGBA{0, 0, 0, 0}
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

	pickerWidth  float32
	pickerHeight float32
	barWidth     float32
	hue          float32
	colorMarker  marker
	hueMarker    barMarker
	*alphaPickerBar
}

func newDefaultHueColorPicker(size float32) ColorPicker {
	pickerSize := fyne.NewSize(size, size)
	barSize := fyne.NewSize(size/10, size)

	picker := &defaultHueColorPicker{
		hue:          0,
		pickerWidth:  pickerSize.Width,
		pickerHeight: pickerSize.Height,
		barWidth:     barSize.Width,
		colorPickerBase: &colorPickerBase{
			changed: func(color.Color) {},
		},
	}

	colorPickerRaster := newTappableRaster(createSaturationValueColorPickerPixelColor(picker.hue))
	colorPickerRaster.SetMinSize(pickerSize)
	colorPickerRaster.tapped = func(p fyne.Position) {
		picker.colorMarker.setPosition(p)
		picker.updatePickerColor()
		colorPickerRaster.Refresh()
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...
	picker.colorPickerRaster = colorPickerRaster

	huePickerRaster := newTappableRaster(hueBarPicker)
	huePickerRaster.SetMinSize(barSize)
	huePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = p.Y / barSize.Height
		colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		setPositionY(picker.hueMarker, p.Y)
		picker.updatePickerColor()
	}
	huePickerRaster.Resize(barSize)

	picker.alphaPickerBar = newAlphaPickerBar(barSize, picker.updatePickerColor)

	picker.colorMarker = newDefaultMarker(5)
	picker.hueMarker = newDefaultBarMarker(picker.barWidth)
	picker.hueMarker.setPosition(fyne.NewPos(picker.hueBarCenter(), 0))

	picker.CanvasObject = newSpaceCenteredLayout(
		fyne.NewContainer(colorPickerRaster, picker.colorMarker.object()),
		fyne.NewContainer(huePickerRaster, picker.hueMarker.object()),
		picker.alphaPickerBar.object(),
	)
	return picker
}

func (p *defaultHueColorPicker) updatePickerColor() {
	x := p.colorMarker.position().X
	y := p.colorMarker.position().Y
	color := fromHSVA(float64(p.hue), float64(x)/float64(p.pickerWidth), 1.0-float64(y)/float64(p.pickerHeight), float64(p.alpha))
	p.changed(color)

	p.alphaPickerBar.setColor(color)
}

func (p *defaultHueColorPicker) SetColor(c color.Color) {
	h, s, v, a := fromColor(c)
	p.hue = float32(h)
	setPositionY(p.hueMarker, p.pickerHeight*float32(h))
	p.colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(p.hue))
	p.colorPickerRaster.Refresh()
	x := float32(math.Round(float64(p.pickerWidth) * s))
	y := float32(math.Round(float64(p.pickerHeight) * (1.0 - v)))
	p.colorMarker.setPosition(fyne.NewPos(x, y))
	p.setAlpha(float32(a))
	p.updatePickerColor()
}

func (p *defaultHueColorPicker) hueBarCenter() float32 {
	return float32(p.barWidth) / 2
}

type circleHueColorPicker struct {
	*colorPickerBase

	pickerWidth    float32
	pickerHeight   float32
	hueCircleWidth float32
	hue            float32
	colorMarker    marker
	hueMarker      barMarker
	*alphaPickerBar
}

func newCircleHueColorPicker(size float32) ColorPicker {
	// pickerAreaWidth < ((areaWidth - (hueBarWidth * 2)) / √2)
	pickerAreaWidth := (size - (size/10)*2) / 1.45
	pickerSize := fyne.NewSize(pickerAreaWidth, pickerAreaWidth)
	hueSize := fyne.NewSize(size, size)
	barSize := fyne.NewSize(size/10, size)

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
		picker.colorMarker.setPosition(p)
		picker.updatePickerColor()
		colorPickerRaster.Refresh()
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...
	picker.colorPickerRaster = colorPickerRaster

	circleHuePickerRaster := newTappableRaster(circleHuePicker)
	circleHuePickerRaster.SetMinSize(hueSize)
	circleHuePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = picker.hueMarker.calcValueFromPosition(p)
		colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.hueMarker.setPosition(p)
		picker.updatePickerColor()
	}
	circleHuePickerRaster.Resize(hueSize)

	picker.alphaPickerBar = newAlphaPickerBar(barSize, picker.updatePickerColor)

	picker.colorMarker = newDefaultMarker(5)
	picker.hueMarker = newCircleBarMarker(hueSize.Width, hueSize.Height, picker.cirlceHueBarWidth())

	picker.CanvasObject = newSpaceCenteredLayout(
		fyne.NewContainerWithLayout(
			layout.NewCenterLayout(),
			fyne.NewContainer(
				circleHuePickerRaster,
				picker.hueMarker.object(),
			),
			fyne.NewContainer(
				colorPickerRaster,
				picker.colorMarker.object(),
			),
		),
		picker.alphaPickerBar.object(),
	)
	return picker
}

func (p *circleHueColorPicker) cirlceHueBarWidth() float32 {
	return float32(p.hueCircleWidth) / 10
}

func (p *circleHueColorPicker) updatePickerColor() {
	x := p.colorMarker.position().X
	y := p.colorMarker.position().Y
	color := fromHSVA(float64(p.hue), float64(x)/float64(p.pickerWidth), 1.0-float64(y)/float64(p.pickerHeight), float64(p.alpha))
	p.changed(color)

	p.alphaPickerBar.setColor(color)
}

func (p *circleHueColorPicker) SetColor(c color.Color) {
	h, s, v, a := fromColor(c)
	p.hue = float32(h)
	p.hueMarker.setPositionFromValue(p.hue)
	p.colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(p.hue))
	p.colorPickerRaster.Refresh()
	x := float32(math.Round(float64(p.pickerWidth) * s))
	y := float32(math.Round(float64(p.pickerHeight) * (1.0 - v)))
	p.colorMarker.setPosition(fyne.NewPos(x, y))
	p.setAlpha(float32(a))
	p.updatePickerColor()
}

type valueColorPicker struct {
	*colorPickerBase

	pickerRadius      float32
	pickerCenter      fyne.Position
	valueBarWidth     float32
	value             float32
	colorMarker       marker
	valueMarker       barMarker
	valuePickerRaster *tappableRaster
	*alphaPickerBar
}

func newValueColorPicker(size float32) ColorPicker {
	pickerSize := fyne.NewSize(size, size)
	barSize := fyne.NewSize(size/10, size)

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
			picker.colorMarker.setPosition(p)
			picker.updatePickerColor()
			colorPickerRaster.Refresh()
		}
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...
	picker.colorPickerRaster = colorPickerRaster

	valuePickerRaster := newTappableRaster(createValueBarPicker(0., 0.))
	valuePickerRaster.SetMinSize(barSize)
	valuePickerRaster.tapped = func(p fyne.Position) {
		picker.value = 1.0 - p.Y/barSize.Height
		colorPickerRaster.setPixelColor(createCircleHueSaturationColorPickerPixelColor(picker.value))
		colorPickerRaster.Refresh()
		setPositionY(picker.valueMarker, p.Y)
		picker.updatePickerColor()
	}
	valuePickerRaster.Resize(barSize)
	picker.valuePickerRaster = valuePickerRaster

	picker.alphaPickerBar = newAlphaPickerBar(barSize, picker.updatePickerColor)

	picker.colorMarker = newDefaultMarker(5)
	picker.colorMarker.setPosition(picker.pickerCenter)
	picker.valueMarker = newDefaultBarMarker(picker.valueBarWidth)
	picker.valueMarker.setPosition(fyne.NewPos(picker.valueBarCenter(), 0))

	picker.CanvasObject = newSpaceCenteredLayout(
		fyne.NewContainer(colorPickerRaster, picker.colorMarker.object()),
		fyne.NewContainer(valuePickerRaster, picker.valueMarker.object()),
		picker.alphaPickerBar.object(),
	)
	return picker
}

func (p *valueColorPicker) SetColor(c color.Color) {
	h, s, v, a := fromColor(c)
	p.value = float32(v)
	areaSize := p.pickerRadius * 2
	setPositionY(p.valueMarker, areaSize*(1.0-p.value))
	p.colorPickerRaster.setPixelColor(createCircleHueSaturationColorPickerPixelColor(p.value))
	p.colorPickerRaster.Refresh()

	baseV := newVector(1, 0)
	rad := -2 * math.Pi * h
	vec := baseV.rotate(rad).multiply(float64(p.pickerRadius) * s)
	center := newVector(float64(p.pickerCenter.X), float64(p.pickerCenter.Y))
	p.colorMarker.setPosition(center.add(vec).toPosition())
	p.setAlpha(float32(a))
	p.updatePickerColor()
}

func (p *valueColorPicker) updatePickerColor() {
	c := calcColorFromCirclePointAndValue(
		float64(p.colorMarker.position().X),
		float64(p.colorMarker.position().Y),
		float64(p.pickerCenter.X),
		float64(p.pickerCenter.Y),
		float64(p.value),
	)
	rgba, _ := c.(color.NRGBA)
	rgba.A = roundUint8(float64(p.alpha) * 255)
	p.changed(rgba)

	// TODO: should not recalculate...
	h, s, v, _ := fromColor(c)
	if v > 0 {
		p.valuePickerRaster.setPixelColor(createValueBarPicker(float32(h), float32(s)))
		p.valuePickerRaster.Refresh()
	}

	p.alphaPickerBar.setColor(rgba)
}

func (p *valueColorPicker) isInPickerArea(pos fyne.Position) bool {
	d := distance(float64(pos.X), float64(pos.Y), float64(p.pickerCenter.X), float64(p.pickerCenter.Y))
	return d <= float64(p.pickerRadius)
}

func (p *valueColorPicker) valueBarCenter() float32 {
	return float32(p.valueBarWidth) / 2
}

type saturationColorPicker struct {
	*colorPickerBase

	pickerWidth            float32
	pickerHeight           float32
	saturationBarWidth     float32
	saturation             float32
	colorMarker            marker
	saturationMarker       barMarker
	saturationPickerRaster *tappableRaster
	*alphaPickerBar
}

func newSaturationColorPicker(size float32) ColorPicker {
	pickerSize := fyne.NewSize(size, size)
	barSize := fyne.NewSize(size/10, size)

	picker := &saturationColorPicker{
		saturation:         0,
		pickerWidth:        pickerSize.Width,
		pickerHeight:       pickerSize.Height,
		saturationBarWidth: barSize.Width,
		colorPickerBase: &colorPickerBase{
			changed: func(color.Color) {},
		},
	}

	colorPickerRaster := newTappableRaster(createHueValueColorPickerPixelColor(picker.saturation))
	colorPickerRaster.SetMinSize(pickerSize)
	colorPickerRaster.tapped = func(p fyne.Position) {
		picker.colorMarker.setPosition(p)
		picker.updatePickerColor()
		colorPickerRaster.Refresh()
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...
	picker.colorPickerRaster = colorPickerRaster

	saturationPickerRaster := newTappableRaster(createSaturationBarPicker(0., 1.))
	saturationPickerRaster.SetMinSize(barSize)
	saturationPickerRaster.tapped = func(p fyne.Position) {
		picker.saturation = 1.0 - p.Y/barSize.Height
		colorPickerRaster.setPixelColor(createHueValueColorPickerPixelColor(picker.saturation))
		colorPickerRaster.Refresh()
		setPositionY(picker.saturationMarker, p.Y)
		picker.updatePickerColor()
	}
	saturationPickerRaster.Resize(barSize)
	picker.saturationPickerRaster = saturationPickerRaster

	picker.alphaPickerBar = newAlphaPickerBar(barSize, picker.updatePickerColor)

	picker.colorMarker = newDefaultMarker(5)
	picker.saturationMarker = newDefaultBarMarker(picker.saturationBarWidth)
	picker.saturationMarker.setPosition(fyne.NewPos(picker.saturationBarCenter(), 0))

	picker.CanvasObject = newSpaceCenteredLayout(
		fyne.NewContainer(colorPickerRaster, picker.colorMarker.object()),
		fyne.NewContainer(saturationPickerRaster, picker.saturationMarker.object()),
		picker.alphaPickerBar.object(),
	)
	return picker
}

func (p *saturationColorPicker) updatePickerColor() {
	x := p.colorMarker.position().X
	y := p.colorMarker.position().Y
	color := fromHSVA(float64(x)/float64(p.pickerWidth), float64(p.saturation), 1.0-float64(y)/float64(p.pickerHeight), float64(p.alpha))
	p.changed(color)

	// TODO: should not recalculate...
	h, s, v, _ := fromColor(color)
	if s > 0 {
		p.saturationPickerRaster.setPixelColor(createSaturationBarPicker(h, v))
		p.saturationPickerRaster.Refresh()
	}

	p.alphaPickerBar.setColor(color)
}

func (p *saturationColorPicker) SetColor(c color.Color) {
	h, s, v, a := fromColor(c)
	p.saturation = float32(s)
	setPositionY(p.saturationMarker, p.pickerHeight*(1.0-float32(s)))
	p.colorPickerRaster.setPixelColor(createHueValueColorPickerPixelColor(p.saturation))
	p.colorPickerRaster.Refresh()
	x := float32(math.Round(float64(p.pickerWidth) * h))
	y := float32(math.Round(float64(p.pickerHeight) * (1.0 - v)))
	p.colorMarker.setPosition(fyne.NewPos(x, y))
	p.setAlpha(float32(a))
	p.updatePickerColor()
}

func (p *saturationColorPicker) saturationBarCenter() float32 {
	return float32(p.saturationBarWidth) / 2
}

type alphaPickerBar struct {
	alpha  float32
	marker barMarker
	raster *tappableRaster

	barHeight float32
}

func newAlphaPickerBar(size fyne.Size, tapped func()) *alphaPickerBar {
	bar := &alphaPickerBar{
		alpha:     1.,
		barHeight: size.Height,
	}

	alphaPickerRaster := newTappableRaster(createAlphaBarPickerPixelColor(transparent))
	alphaPickerRaster.SetMinSize(size)
	alphaPickerRaster.tapped = func(p fyne.Position) {
		bar.alpha = 1. - (p.Y / size.Height)
		setPositionY(bar.marker, p.Y)
		tapped()
	}
	alphaPickerRaster.Resize(size)
	bar.raster = alphaPickerRaster

	bar.marker = newDefaultBarMarker(size.Width)
	bar.marker.setPosition(fyne.NewPos(float32(size.Width)/2, 0))

	return bar
}

func (b *alphaPickerBar) object() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(
		layout.NewMaxLayout(),
		newCheckeredBackground(),
		fyne.NewContainer(b.raster, b.marker.object()),
	)
}

func (b *alphaPickerBar) setColor(c color.Color) {
	b.raster.setPixelColor(createAlphaBarPickerPixelColor(c))
	b.raster.Refresh()
}

func (b *alphaPickerBar) setAlpha(a float32) {
	b.alpha = a
	setPositionY(b.marker, b.barHeight*(1.-a))
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

func createAlphaBarPickerPixelColor(c color.Color) func(int, int, int, int) color.Color {
	return func(x, y, w, h int) color.Color {
		r, g, b, _ := toFloatRGBA(c)
		a := 1. - (float64(y) / float64(h))
		return fromFloatNRGBA(r, g, b, a)
	}
}

func createSaturationValueColorPickerPixelColor(hue float32) func(int, int, int, int) color.Color {
	return func(x, y, w, h int) color.Color {
		return fromHSV(float64(hue), float64(x)/float64(w), 1.0-float64(y)/float64(h))
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
		return transparent
	}

	rad := math.Atan2(y-cy, cx-x)
	rad += math.Pi
	hue := rad / (2 * math.Pi)

	return fromHSV(hue, 1.0, 1.0)
}

func createCircleHueSaturationColorPickerPixelColor(value float32) func(int, int, int, int) color.Color {
	return func(x, y, w, h int) color.Color {
		return calcColorFromCirclePointAndValue(float64(x), float64(y), float64(w)/2., float64(h)/2., float64(value))
	}
}

func calcColorFromCirclePointAndValue(x, y, cx, cy, value float64) color.Color {
	dist := distance(x, y, cx, cy)
	if cx < dist {
		return transparent
	}

	rad := math.Atan2(y-cy, cx-x)
	rad += math.Pi
	hue := rad / (2 * math.Pi)

	return fromHSV(hue, dist/cx, value)
}

func createValueBarPicker(hue, saturation float32) func(x, y, w, h int) color.Color {
	return func(x, y, w, h int) color.Color {
		return fromHSV(float64(hue), float64(saturation), 1.0-float64(y)/float64(h))
	}
}

func createHueValueColorPickerPixelColor(saturation float32) func(int, int, int, int) color.Color {
	return func(x, y, w, h int) color.Color {
		return fromHSV(float64(x)/float64(w), float64(saturation), 1.0-float64(y)/float64(h))
	}
}

func createSaturationBarPicker(hue, value float64) func(x, y, w, h int) color.Color {
	return func(x, y, w, h int) color.Color {
		return fromHSV(hue, 1.0-float64(y)/float64(h), value)
	}
}

func newSpaceCenteredLayout(objects ...fyne.CanvasObject) *fyne.Container {
	l := newSpacedLayout(
		layout.NewVBoxLayout(),
		newSpacedLayout(
			layout.NewHBoxLayout(),
			objects...,
		),
	)
	l.Resize(l.MinSize())
	return l
}

func newSpacedLayout(l fyne.Layout, objects ...fyne.CanvasObject) *fyne.Container {
	c := fyne.NewContainerWithLayout(l)
	c.AddObject(layout.NewSpacer())
	for _, o := range objects {
		c.AddObject(o)
	}
	c.AddObject(layout.NewSpacer())
	return c
}

func newCheckeredBackground() *canvas.Raster {
	return canvas.NewRasterWithPixels(func(x, y, _, _ int) color.Color {
		const boxSize = 10

		if (x/boxSize)%2 == (y/boxSize)%2 {
			return color.Gray{Y: 58}
		}

		return color.Gray{Y: 84}
	})
}
