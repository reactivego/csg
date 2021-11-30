package csg

import (
	"math"
)

// Cylinder constructs a solid cylinder. Optional parameters are `Start`,
// `End`, `Radius`, and `Slices`, which default to `Start(0, -1, 0)`,
// `End(0, 1, 0)`, `Radius(1)`, and `Slices(16)`. The `Slices` parameter
// controls the tessellation.
//
// Example usage:
//
//     cylinder := csg.Cylinder(
//       Start(0, -1, 0),
//       End(0, 1, 0),
//       Radius(1),
//       Slices(16))
//
func Cylinder(options ...Option) *Solid {
	o := OptionsFrom(options)
	s := o.Start
	e := o.End
	ray := e.Minus(s)
	r := o.Radius
	slices := float64(o.Slices)
	axisZ := ray.Unit()
	axisX := Vector{0, 1, 0}
	if math.Abs(axisZ.Y) > 0.5 {
		axisX = Vector{1, 0, 0}
	}
	axisX = axisX.Cross(axisZ).Unit()
	axisY := axisX.Cross(axisZ).Unit()
	start := Vertex{Pos: s, Normal: axisZ.Negated()}
	end := Vertex{Pos: e, Normal: axisZ.Unit()}
	point := func(stack, slice, normalBlend float64) Vertex {
		angle := slice * math.Pi * 2
		out := axisX.Times(math.Cos(angle)).Plus(axisY.Times(math.Sin(angle)))
		pos := s.Plus(ray.Times(stack)).Plus(out.Times(r))
		normal := out.Times(1 - math.Abs(normalBlend)).Plus(axisZ.Times(normalBlend))
		return Vertex{Pos: pos, Normal: normal}
	}
	var polygons Polygons
	for i := 0.0; i < slices; i++ {
		t0, t1 := i/slices, (i+1)/slices
		polygons = append(polygons, PolygonFromVertices(start, point(0, t0, -1), point(0, t1, -1)))
		polygons = append(polygons, PolygonFromVertices(point(0, t1, 0), point(0, t0, 0), point(1, t0, 0), point(1, t1, 0)))
		polygons = append(polygons, PolygonFromVertices(end, point(1, t1, 1), point(1, t0, 1)))
	}
	return SolidFromPolygons(polygons)
}
