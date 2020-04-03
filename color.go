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

func round(v float64) uint8 {
	return uint8(math.Round(v))
}
