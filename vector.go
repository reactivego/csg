package csg

import "math"

// Vector represents a 3D vector.
//
// Example usage:
//
//     csg.Vector{1, 2, 3}
//     csg.Vector{X: 1, Y: 2, Z: 3}
//
type Vector struct {
	X, Y, Z float64
}

func (v Vector) Negated() Vector {
	return Vector{-v.X, -v.Y, -v.Z}
}

func (v Vector) Plus(a Vector) Vector {
	return Vector{v.X + a.X, v.Y + a.Y, v.Z + a.Z}
}

func (v Vector) Minus(a Vector) Vector {
	return Vector{v.X - a.X, v.Y - a.Y, v.Z - a.Z}
}

func (v Vector) Times(a float64) Vector {
	return Vector{v.X * a, v.Y * a, v.Z * a}
}

func (v Vector) DividedBy(a float64) Vector {
	return Vector{v.X / a, v.Y / a, v.Z / a}
}

func (v Vector) Dot(a Vector) float64 {
	return v.X*a.X + v.Y*a.Y + v.Z*a.Z
}

func (v Vector) Lerp(a Vector, t float64) Vector {
	return v.Plus(a.Minus(v).Times(t))
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector) Unit() Vector {
	return v.DividedBy(v.Length())
}

func (v Vector) Cross(a Vector) Vector {
	return Vector{
		v.Y*a.Z - v.Z*a.Y,
		v.Z*a.X - v.X*a.Z,
		v.X*a.Y - v.Y*a.X,
	}
}
