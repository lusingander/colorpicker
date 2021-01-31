package colorpicker

import (
	"image/color"
	"math"
)

func fromHSV(h, s, v float64) color.NRGBA {
	if s == 0 {
		return color.NRGBA{
			R: roundUint8(v * 255),
			G: roundUint8(v * 255),
			B: roundUint8(v * 255),
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
		return color.NRGBA{
			R: roundUint8(v * 255),
			G: roundUint8(v3 * 255),
			B: roundUint8(v1 * 255),
			A: 0xff,
		}
	case 1:
		return color.NRGBA{
			R: roundUint8(v2 * 255),
			G: roundUint8(v * 255),
			B: roundUint8(v1 * 255),
			A: 0xff,
		}
	case 2:
		return color.NRGBA{
			R: roundUint8(v1 * 255),
			G: roundUint8(v * 255),
			B: roundUint8(v3 * 255),
			A: 0xff,
		}
	case 3:
		return color.NRGBA{
			R: roundUint8(v1 * 255),
			G: roundUint8(v2 * 255),
			B: roundUint8(v * 255),
			A: 0xff,
		}
	case 4:
		return color.NRGBA{
			R: roundUint8(v3 * 255),
			G: roundUint8(v1 * 255),
			B: roundUint8(v * 255),
			A: 0xff,
		}
	default:
		return color.NRGBA{
			R: roundUint8(v * 255),
			G: roundUint8(v1 * 255),
			B: roundUint8(v2 * 255),
			A: 0xff,
		}
	}
}

func fromHSVA(h, s, v, a float64) color.NRGBA {
	rgba := fromHSV(h, s, v)
	rgba.A = roundUint8(a * 255)
	return rgba
}

func fromColor(c color.Color) (h, s, v, a float64) {
	r, g, b, a := toFloatRGBA(c)
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

func toFloatRGBA(c color.Color) (float64, float64, float64, float64) {
	rgba, _ := c.(color.NRGBA)
	max := 255.
	return float64(rgba.R) / max, float64(rgba.G) / max, float64(rgba.B) / max, float64(rgba.A) / max
}

func roundUint8(v float64) uint8 {
	return uint8(math.Round(v))
}

func fromFloatNRGBA(r, g, b, a float64) color.Color {
	return color.NRGBA{
		R: roundUint8(r * 255),
		G: roundUint8(g * 255),
		B: roundUint8(b * 255),
		A: roundUint8(a * 255),
	}
}
