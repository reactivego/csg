package csg

import "math"

// Sphere constructs a solid sphere. Optional parameters are `Center`, `Radius`,
// `Slices`, and `Stacks`, which default to `Center(0, 0, 0)`, `Radius(1)`,
// `Slices(16)`, and `Stacks(8)`.
// The `Slices` and `Stacks` parameters control the tessellation along the
// longitude and latitude directions.
//
// Example usage:
//
//     sphere := csg.Sphere(
//       Center(0, 0, 0),
//       Radius(1),
//       Slices(16),
//       Stacks(8))
//
func Sphere(options ...Option) *Solid {
	o := OptionsFrom(options)
	c := o.Center
	r := o.Radius
	vertices := []Vertex(nil)
	vertex := func(theta, phi float64) {
		theta *= math.Pi * 2
		phi *= math.Pi
		dir := Vector{
			X: math.Cos(theta) * math.Sin(phi),
			Y: math.Cos(phi),
			Z: math.Sin(theta) * math.Sin(phi),
		}
		vertices = append(vertices, Vertex{
			Pos:    c.Plus(dir.Times(r)),
			Normal: dir})
	}
	polygons := make([]Polygon, 0, o.Slices*o.Stacks)
	slices, stacks := float64(o.Slices), float64(o.Stacks)
	for i := 0.0; i < slices; i++ {
		for j := 0.0; j < stacks; j++ {
			vertices = []Vertex{}
			vertex(i/slices, j/stacks)
			if j > 0 {
				vertex((i+1)/slices, j/stacks)
			}
			if j < stacks-1 {
				vertex((i+1)/slices, (j+1)/stacks)
			}
			vertex(i/slices, (j+1)/stacks)
			polygons = append(polygons, PolygonFromVertices(vertices...))
		}
	}
	return SolidFromPolygons(polygons)
}
