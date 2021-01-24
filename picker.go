package colorpicker

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

var (
	transparent = color.RGBA{0, 0, 0, 0}
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
	hueBarWidth  float32
	hue          float32
	colorMarker  *marker
	hueMarker    *marker
}

func newDefaultHueColorPicker(size float32) ColorPicker {
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
		picker.colorMarker.setPosition(p)
		picker.updatePickerColor()
		colorPickerRaster.Refresh()
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...
	picker.colorPickerRaster = colorPickerRaster

	huePickerRaster := newTappableRaster(hueBarPicker)
	huePickerRaster.SetMinSize(hueSize)
	huePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = p.Y / hueSize.Height
		colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.hueMarker.setPositionY(p.Y)
		picker.updatePickerColor()
	}
	huePickerRaster.Resize(hueSize)

	picker.colorMarker = newMarker(5, 1)
	picker.hueMarker = newMarker(picker.hueBarCenter(), 2)
	picker.hueMarker.setPosition(fyne.NewPos(picker.hueBarCenter(), 0))

	picker.CanvasObject = newSpaceCenteredLayout(
		fyne.NewContainer(colorPickerRaster, picker.colorMarker.Circle),
		fyne.NewContainer(huePickerRaster, picker.hueMarker.Circle),
	)
	return picker
}

func (p *defaultHueColorPicker) updatePickerColor() {
	x := p.colorMarker.center.X
	y := p.colorMarker.center.Y
	color := fromHSV(float64(p.hue), float64(x)/float64(p.pickerWidth), 1.0-float64(y)/float64(p.pickerHeight))
	p.changed(color)
}

func (p *defaultHueColorPicker) SetColor(c color.Color) {
	h, s, v := fromColor(c)
	p.hue = float32(h)
	p.hueMarker.setPositionY(p.pickerHeight * float32(h))
	p.colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(p.hue))
	p.colorPickerRaster.Refresh()
	x := float32(math.Round(float64(p.pickerWidth) * s))
	y := float32(math.Round(float64(p.pickerHeight) * (1.0 - v)))
	p.colorMarker.setPosition(fyne.NewPos(x, y))
	p.updatePickerColor()
}

func (p *defaultHueColorPicker) hueBarCenter() float32 {
	return float32(p.hueBarWidth) / 2
}

type circleHueColorPicker struct {
	*colorPickerBase

	pickerWidth    float32
	pickerHeight   float32
	hueCircleWidth float32
	hue            float32
	colorMarker    *marker
	hueMarker      *circleBarMarker
}

func newCircleHueColorPicker(size float32) ColorPicker {
	// pickerAreaWidth < ((areaWidth - (hueBarWidth * 2)) / √2)
	pickerAreaWidth := (size - (size/10)*2) / 1.45
	pickerSize := fyne.NewSize(pickerAreaWidth, pickerAreaWidth)
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
		picker.colorMarker.setPosition(p)
		picker.updatePickerColor()
		colorPickerRaster.Refresh()
	}
	colorPickerRaster.Resize(pickerSize) // Note: doesn't render if remove this line...
	picker.colorPickerRaster = colorPickerRaster

	circleHuePickerRaster := newTappableRaster(circleHuePicker)
	circleHuePickerRaster.SetMinSize(hueSize)
	circleHuePickerRaster.tapped = func(p fyne.Position) {
		picker.hue = picker.hueMarker.calcHueFromCircleMarker(p)
		colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(picker.hue))
		colorPickerRaster.Refresh()
		picker.hueMarker.setCircleMarkerPosition(p)
		picker.updatePickerColor()
	}
	circleHuePickerRaster.Resize(hueSize)

	picker.colorMarker = newMarker(5, 1)
	picker.hueMarker = newCircleBarMarker(hueSize.Width, hueSize.Height, picker.cirlceHueBarWidth())

	picker.CanvasObject = newSpaceCenteredLayout(
		fyne.NewContainerWithLayout(
			layout.NewCenterLayout(),
			fyne.NewContainer(
				circleHuePickerRaster,
				picker.hueMarker.Circle,
			),
			fyne.NewContainer(
				colorPickerRaster,
				picker.colorMarker.Circle,
			),
		),
	)
	return picker
}

func (p *circleHueColorPicker) cirlceHueBarWidth() float32 {
	return float32(p.hueCircleWidth) / 10
}

func (p *circleHueColorPicker) updatePickerColor() {
	x := p.colorMarker.center.X
	y := p.colorMarker.center.Y
	color := fromHSV(float64(p.hue), float64(x)/float64(p.pickerWidth), 1.0-float64(y)/float64(p.pickerHeight))
	p.changed(color)
}

func (p *circleHueColorPicker) SetColor(c color.Color) {
	h, s, v := fromColor(c)
	p.hue = float32(h)
	p.hueMarker.setCircleMarekerPositionFromHue(p.hue)
	p.colorPickerRaster.setPixelColor(createSaturationValueColorPickerPixelColor(p.hue))
	p.colorPickerRaster.Refresh()
	x := float32(math.Round(float64(p.pickerWidth) * s))
	y := float32(math.Round(float64(p.pickerHeight) * (1.0 - v)))
	p.colorMarker.setPosition(fyne.NewPos(x, y))
	p.updatePickerColor()
}

type valueColorPicker struct {
	*colorPickerBase

	pickerRadius      float32
	pickerCenter      fyne.Position
	valueBarWidth     float32
	value             float32
	colorMarker       *marker
	valueMarker       *marker
	valuePickerRaster *tappableRaster
}

func newValueColorPicker(size float32) ColorPicker {
	pickerSize := fyne.NewSize(size, size)
	valueSize := fyne.NewSize(size/10, size)

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
	valuePickerRaster.SetMinSize(valueSize)
	valuePickerRaster.tapped = func(p fyne.Position) {
		picker.value = 1.0 - p.Y/valueSize.Height
		colorPickerRaster.setPixelColor(createCircleHueSaturationColorPickerPixelColor(picker.value))
		colorPickerRaster.Refresh()
		picker.valueMarker.setPositionY(p.Y)
		picker.updatePickerColor()
	}
	valuePickerRaster.Resize(valueSize)
	picker.valuePickerRaster = valuePickerRaster

	picker.colorMarker = newMarker(5, 1)
	picker.colorMarker.setPosition(picker.pickerCenter)
	picker.valueMarker = newMarker(picker.valueBarCenter(), 2)
	picker.valueMarker.setPosition(fyne.NewPos(picker.valueBarCenter(), 0))

	picker.CanvasObject = newSpaceCenteredLayout(
		fyne.NewContainer(colorPickerRaster, picker.colorMarker.Circle),
		fyne.NewContainer(valuePickerRaster, picker.valueMarker.Circle),
	)
	return picker
}

func (p *valueColorPicker) SetColor(c color.Color) {
	h, s, v := fromColor(c)
	p.value = float32(v)
	areaSize := p.pickerRadius * 2
	p.valueMarker.setPositionY(areaSize * (1.0 - p.value))
	p.colorPickerRaster.setPixelColor(createCircleHueSaturationColorPickerPixelColor(p.value))
	p.colorPickerRaster.Refresh()

	baseV := newVector(1, 0)
	rad := -2 * math.Pi * h
	vec := baseV.rotate(rad).multiply(float64(p.pickerRadius) * s)
	center := newVector(float64(p.pickerCenter.X), float64(p.pickerCenter.Y))
	p.colorMarker.setPosition(center.add(vec).toPosition())
	p.updatePickerColor()
}

func (p *valueColorPicker) updatePickerColor() {
	color := calcColorFromCirclePointAndValue(
		float64(p.colorMarker.center.X),
		float64(p.colorMarker.center.Y),
		float64(p.pickerCenter.X),
		float64(p.pickerCenter.Y),
		float64(p.value),
	)
	p.changed(color)

	// TODO: should not recalculate...
	h, s, v := fromColor(color)
	if v > 0 {
		p.valuePickerRaster.setPixelColor(createValueBarPicker(float32(h), float32(s)))
		p.valuePickerRaster.Refresh()
	}
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
	colorMarker            *marker
	saturationMarker       *marker
	saturationPickerRaster *tappableRaster
}

func newSaturationColorPicker(size float32) ColorPicker {
	pickerSize := fyne.NewSize(size, size)
	saturationSize := fyne.NewSize(size/10, size)

	picker := &saturationColorPicker{
		saturation:         0,
		pickerWidth:        pickerSize.Width,
		pickerHeight:       pickerSize.Height,
		saturationBarWidth: saturationSize.Width,
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
	saturationPickerRaster.SetMinSize(saturationSize)
	saturationPickerRaster.tapped = func(p fyne.Position) {
		picker.saturation = 1.0 - p.Y/saturationSize.Height
		colorPickerRaster.setPixelColor(createHueValueColorPickerPixelColor(picker.saturation))
		colorPickerRaster.Refresh()
		picker.saturationMarker.setPositionY(p.Y)
		picker.updatePickerColor()
	}
	saturationPickerRaster.Resize(saturationSize)
	picker.saturationPickerRaster = saturationPickerRaster

	picker.colorMarker = newMarker(5, 1)
	picker.saturationMarker = newMarker(picker.saturationBarCenter(), 2)
	picker.saturationMarker.setPosition(fyne.NewPos(picker.saturationBarCenter(), 0))

	picker.CanvasObject = newSpaceCenteredLayout(
		fyne.NewContainer(colorPickerRaster, picker.colorMarker.Circle),
		fyne.NewContainer(saturationPickerRaster, picker.saturationMarker.Circle),
	)
	return picker
}

func (p *saturationColorPicker) updatePickerColor() {
	x := p.colorMarker.center.X
	y := p.colorMarker.center.Y
	color := fromHSV(float64(x)/float64(p.pickerWidth), float64(p.saturation), 1.0-float64(y)/float64(p.pickerHeight))
	p.changed(color)

	// TODO: should not recalculate...
	h, s, v := fromColor(color)
	if s > 0 {
		p.saturationPickerRaster.setPixelColor(createSaturationBarPicker(h, v))
		p.saturationPickerRaster.Refresh()
	}
}

func (p *saturationColorPicker) SetColor(c color.Color) {
	h, s, v := fromColor(c)
	p.saturation = float32(s)
	p.saturationMarker.setPositionY(p.pickerHeight * (1.0 - float32(s)))
	p.colorPickerRaster.setPixelColor(createHueValueColorPickerPixelColor(p.saturation))
	p.colorPickerRaster.Refresh()
	x := float32(math.Round(float64(p.pickerWidth) * h))
	y := float32(math.Round(float64(p.pickerHeight) * (1.0 - v)))
	p.colorMarker.setPosition(fyne.NewPos(x, y))
	p.updatePickerColor()
}

func (p *saturationColorPicker) saturationBarCenter() float32 {
	return float32(p.saturationBarWidth) / 2
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
