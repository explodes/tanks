package tanksutil

import (
	"math"
)

type Vec struct {
	X, Y float64
}

func V(x, y float64) Vec {
	return Vec{
		X: x,
		Y: y,
	}
}

// Angle returns the angle between the vector u and the x-axis. The result is in range [-Pi, Pi].
func (u Vec) Angle() float64 {
	return math.Atan2(u.Y, u.X)
}

// Scaled returns the vector u multiplied by c.
func (u Vec) Scaled(p float64) Vec {
	return Vec{
		X: u.X * p,
		Y: u.Y * p,
	}
}

// Add returns the sum of vectors u and v.
func (u Vec) Add(v Vec) Vec {
	return Vec{
		u.X + v.X,
		u.Y + v.Y,
	}
}

// Rotated returns the vector u rotated by the given angle in radians.
func (u Vec) Rotated(angle float64) Vec {
	sin, cos := math.Sincos(angle)
	return Vec{
		u.X*cos - u.Y*sin,
		u.X*sin + u.Y*cos,
	}
}

type Rect struct {
	Min, Max Vec
}

func R(x1, y1, x2, y2 float64) Rect {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	return Rect{
		Min: V(x1, y1),
		Max: V(x2, y2),
	}
}

func (r Rect) W() float64 {
	return r.Max.X - r.Min.X
}

func (r Rect) H() float64 {
	return r.Max.Y - r.Min.Y
}

func (r Rect) Center() Vec {
	return V(
		r.Min.X+r.W()/2,
		r.Min.Y+r.H()/2,
	)
}