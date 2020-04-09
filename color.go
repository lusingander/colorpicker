package colorpicker

import (
	"image/color"
	"math"
)

func fromHSV(h, s, v float64) *color.RGBA {
	if s == 0 {
		return &color.RGBA{
			R: round(v * 255),
			G: round(v * 255),
			B: round(v * 255),
			A: 0xff,
		}
	}

	h = h * 6
	if h == 6 {
		h = 0
	}
	i := math.Floor(h)
	v1 := v * (1 - s)
	v2 := v * (1 - s*(h-i))
	v3 := v * (1 - s*(1-(h-i)))

	switch int(i) {
	case 0:
		return &color.RGBA{
			R: round(v * 255),
			G: round(v3 * 255),
			B: round(v1 * 255),
			A: 0xff,
		}
	case 1:
		return &color.RGBA{
			R: round(v2 * 255),
			G: round(v * 255),
			B: round(v1 * 255),
			A: 0xff,
		}
	case 2:
		return &color.RGBA{
			R: round(v1 * 255),
			G: round(v * 255),
			B: round(v3 * 255),
			A: 0xff,
		}
	case 3:
		return &color.RGBA{
			R: round(v1 * 255),
			G: round(v2 * 255),
			B: round(v * 255),
			A: 0xff,
		}
	case 4:
		return &color.RGBA{
			R: round(v3 * 255),
			G: round(v1 * 255),
			B: round(v * 255),
			A: 0xff,
		}
	default:
		return &color.RGBA{
			R: round(v * 255),
			G: round(v1 * 255),
			B: round(v2 * 255),
			A: 0xff,
		}
	}
}

func fromColor(c color.Color) (h, s, v float64) {
	r, g, b := toFloatRGB(c)
	min := math.Min(r, math.Min(g, b))
	max := math.Max(r, math.Max(g, b))
	v = max

	d := max - min
	if d == 0 {
		h = 0
		s = 0
		return
	}
	s = d / max
	if r == max {
		h = (1. / 6.) * ((g - b) / d)
	} else if g == max {
		h = (1./6.)*((b-r)/d) + (1. / 3.)
	} else { // b == max
		h = (1./6.)*((r-g)/d) + (2. / 3.)
	}
	if h < 0 {
		h++
	} else if h > 1 {
		h--
	}
	return
}

func toFloatRGB(c color.Color) (float64, float64, float64) {
	max := 65535.
	r, g, b, _ := c.RGBA()
	return float64(r) / max, float64(g) / max, float64(b) / max
}

func round(v float64) uint8 {
	return uint8(math.Round(v))
}
