package colorpicker

import (
	"image/color"
	"math"
	"testing"
)

const (
	floatThreshold = 0.001
)

func TestFromHSV(t *testing.T) {
	tests := []struct {
		h, s, v float64
		want    color.Color
	}{
		{0 / 360., 1., 1., color.NRGBA{0xff, 0x00, 0x00, 0xff}},
		{30 / 360., 1., 1., color.NRGBA{0xff, 0x80, 0x00, 0xff}},
		{60 / 360., 1., 1., color.NRGBA{0xff, 0xff, 0x00, 0xff}},
		{90 / 360., 1., 1., color.NRGBA{0x80, 0xff, 0x00, 0xff}},
		{120 / 360., 1., 1., color.NRGBA{0x00, 0xff, 0x00, 0xff}},
		{150 / 360., 1., 1., color.NRGBA{0x00, 0xff, 0x80, 0xff}},
		{180 / 360., 1., 1., color.NRGBA{0x00, 0xff, 0xff, 0xff}},
		{210 / 360., 1., 1., color.NRGBA{0x00, 0x80, 0xff, 0xff}},
		{240 / 360., 1., 1., color.NRGBA{0x00, 0x00, 0xff, 0xff}},
		{270 / 360., 1., 1., color.NRGBA{0x80, 0x00, 0xff, 0xff}},
		{300 / 360., 1., 1., color.NRGBA{0xff, 0x00, 0xff, 0xff}},
		{330 / 360., 1., 1., color.NRGBA{0xff, 0x00, 0x80, 0xff}},
		{360 / 360., 1., 1., color.NRGBA{0xff, 0x00, 0x00, 0xff}},
		{0 / 360., 1., 0., color.NRGBA{0x00, 0x00, 0x00, 0xff}},
		{0 / 360., 1., 0.2, color.NRGBA{0x33, 0x00, 0x00, 0xff}},
		{0 / 360., 1., 0.4, color.NRGBA{0x66, 0x00, 0x00, 0xff}},
		{0 / 360., 1., 0.6, color.NRGBA{0x99, 0x00, 0x00, 0xff}},
		{0 / 360., 1., 0.8, color.NRGBA{0xcc, 0x00, 0x00, 0xff}},
		{180 / 360., 0., 1., color.NRGBA{0xff, 0xff, 0xff, 0xff}},
		{180 / 360., 0.2, 1., color.NRGBA{0xcc, 0xff, 0xff, 0xff}},
		{180 / 360., 0.4, 1., color.NRGBA{0x99, 0xff, 0xff, 0xff}},
		{180 / 360., 0.6, 1., color.NRGBA{0x66, 0xff, 0xff, 0xff}},
		{180 / 360., 0.8, 1., color.NRGBA{0x33, 0xff, 0xff, 0xff}},
	}
	for _, test := range tests {
		got := fromHSV(test.h, test.s, test.v)
		if got != test.want {
			t.Errorf("fromHSV(%f, %f, %f) = %v; want %v",
				test.h, test.s, test.v, got, test.want)
		}
	}
}

func TestFromColor(t *testing.T) {
	tests := []struct {
		c       color.Color
		h, s, v float64
	}{
		{color.NRGBA{0xff, 0x00, 0x00, 0xff}, 0 / 360., 1., 1.},
		{color.NRGBA{0xff, 0x80, 0x00, 0xff}, 30 / 360., 1., 1.},
		{color.NRGBA{0xff, 0xff, 0x00, 0xff}, 60 / 360., 1., 1.},
		{color.NRGBA{0x80, 0xff, 0x00, 0xff}, 90 / 360., 1., 1.},
		{color.NRGBA{0x00, 0xff, 0x00, 0xff}, 120 / 360., 1., 1.},
		{color.NRGBA{0x00, 0xff, 0x80, 0xff}, 150 / 360., 1., 1.},
		{color.NRGBA{0x00, 0xff, 0xff, 0xff}, 180 / 360., 1., 1.},
		{color.NRGBA{0x00, 0x80, 0xff, 0xff}, 210 / 360., 1., 1.},
		{color.NRGBA{0x00, 0x00, 0xff, 0xff}, 240 / 360., 1., 1.},
		{color.NRGBA{0x80, 0x00, 0xff, 0xff}, 270 / 360., 1., 1.},
		{color.NRGBA{0xff, 0x00, 0xff, 0xff}, 300 / 360., 1., 1.},
		{color.NRGBA{0xff, 0x00, 0x80, 0xff}, 330 / 360., 1., 1.},
		{color.NRGBA{0x33, 0x00, 0x00, 0xff}, 0 / 360., 1., 0.2},
		{color.NRGBA{0x66, 0x00, 0x00, 0xff}, 0 / 360., 1., 0.4},
		{color.NRGBA{0x99, 0x00, 0x00, 0xff}, 0 / 360., 1., 0.6},
		{color.NRGBA{0xcc, 0x00, 0x00, 0xff}, 0 / 360., 1., 0.8},
		{color.NRGBA{0xcc, 0xff, 0xff, 0xff}, 180 / 360., 0.2, 1.},
		{color.NRGBA{0x99, 0xff, 0xff, 0xff}, 180 / 360., 0.4, 1.},
		{color.NRGBA{0x66, 0xff, 0xff, 0xff}, 180 / 360., 0.6, 1.},
		{color.NRGBA{0x33, 0xff, 0xff, 0xff}, 180 / 360., 0.8, 1.},
	}
	for _, test := range tests {
		h, s, v, _ := fromColor(test.c)
		if notEquals(test.h, h) || notEquals(test.s, s) || notEquals(test.v, v) {
			t.Errorf("fromColor(%v) = %f, %f, %f; want %f, %f, %f",
				test.c, h, s, v, test.h, test.s, test.v)
		}
	}
}

func TestFromColorAndFromHSV(t *testing.T) {
	wants := make([]color.Color, 0, 256*256*256)
	for r := 0; r <= 255; r++ {
		for g := 0; g <= 255; g++ {
			for b := 0; b <= 255; b++ {
				wants = append(wants, color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff})
			}
		}
	}
	for _, want := range wants {
		h, s, v, _ := fromColor(want)
		got := fromHSV(h, s, v)
		if want != got {
			t.Errorf("fromHSV(fromColor(%v)) = %v", want, got)
		}
	}
}

func notEquals(f1, f2 float64) bool {
	return math.Abs(f1-f2) > floatThreshold
}
