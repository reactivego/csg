package csg

// Vertex represents a vertex of a polygon. Use your own vertex class instead of this
// one to provide additional features like texture coordinates and vertex
// colors. Custom vertex classes need to provide a `Pos` property and `Flip()`, and
// `Interpolate()` methods that behave analogous to the ones defined by `csg.Vertex`.
// This struct provides `Normal` so convenience functions like `csg.Sphere()` can
// return a smooth vertex normal, but `Normal` is not used anywhere else.
type Vertex struct{ Pos, Normal Vector }

// Flip inverts all orientation-specific data (e.g. vertex normal). Called when the
// orientation of a polygon is flipped.
func (v *Vertex) Flip() {
	v.Normal = v.Normal.Negated()
}

// Interpolate creates a new vertex between this vertex and `other` by linearly
// interpolating all properties using a parameter of `t`.
func (v Vertex) Interpolate(other Vertex, t float64) Vertex {
	return Vertex{
		Pos:    v.Pos.Lerp(other.Pos, t),
		Normal: v.Normal.Lerp(other.Normal, t),
	}
}
