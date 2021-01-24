package colorpicker

import (
	"math"

	"fyne.io/fyne/v2"
)

type vector struct {
	x, y float64
}

func newVector(x, y float64) *vector {
	return &vector{x, y}
}

func newVectorFromPoints(x1, y1, x2, y2 float64) *vector {
	return &vector{
		x: x2 - x1,
		y: y2 - y1,
	}
}

func (v *vector) add(u *vector) *vector {
	return &vector{v.x + u.x, v.y + u.y}
}

func (v *vector) multiply(s float64) *vector {
	return &vector{v.x * s, v.y * s}
}

func (v *vector) normalize() *vector {
	n := v.norm()
	return &vector{v.x / n, v.y / n}
}

func (v *vector) norm() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *vector) dot(u *vector) float64 {
	return v.x*u.x + v.y*u.y
}

func (v *vector) toPosition() fyne.Position {
	return fyne.NewPos(float32(v.x), float32(v.y))
}

func (v *vector) rotate(rad float64) *vector {
	cos := math.Cos(rad)
	sin := math.Sin(rad)
	return &vector{cos*v.x - sin*v.y, sin*v.x + cos*v.y}
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(square(x1-x2) + square(y1-y2))
}

func square(v float64) float64 {
	return v * v
}
