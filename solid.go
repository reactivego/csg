package csg

// Solid holds a binary space partition tree representing a 3D solid.
// Two solids can be combined using the `Union()`, `Subtract()` and
// `Intersect()` methods.
type Solid struct {
	Polygons Polygons
}

// SolidFromPolygons constructs a CSG solid from a list of `csg.Polygon`
// instances.
func SolidFromPolygons(polygons Polygons) *Solid {
	return &Solid{Polygons: polygons}
}

// Union returns a new CSG solid representing space in either this solid or
// in the solid `csg`. Neither this solid nor the solid `csg` are modified.
//
//     A.Union(B)
//
//     +-------+            +-------+
//     |       |            |       |
//     |   A   |            |       |
//     |    +--+----+   =   |       +----+
//     +----+--+    |       +----+       |
//          |   B   |            |       |
//          |       |            |       |
//          +-------+            +-------+
//
func (s *Solid) Union(other *Solid) *Solid {
	a, b := &BSP{}, &BSP{}
	a.AddPolygons(s.Polygons.Clone())
	b.AddPolygons(other.Polygons.Clone())
	a.ClipTo(b)
	b.ClipTo(a)
	b.Invert()
	b.ClipTo(a)
	b.Invert()
	a.AddPolygons(b.AllPolygons())
	return SolidFromPolygons(a.AllPolygons())
}

// Subtract returns a new CSG solid representing space in this solid but not
// in the solid `csg`. Neither this solid nor the solid `csg` are modified.
//
//     A.Subtract(B)
//
//     +-------+            +-------+
//     |       |            |       |
//     |   A   |            |       |
//     |    +--+----+   =   |    +--+
//     +----+--+    |       +----+
//          |   B   |
//          |       |
//          +-------+
//
func (s *Solid) Subtract(other *Solid) *Solid {
	a, b := &BSP{}, &BSP{}
	a.AddPolygons(s.Polygons.Clone())
	b.AddPolygons(other.Polygons.Clone())
	a.Invert()
	a.ClipTo(b)
	b.ClipTo(a)
	b.Invert()
	b.ClipTo(a)
	b.Invert()
	a.AddPolygons(b.AllPolygons())
	a.Invert()
	return SolidFromPolygons(a.AllPolygons())
}

// Intersect returns a new CSG solid representing space both this solid and in
// the solid `csg`. Neither this solid nor the solid `csg` are modified.
//
//     A.intersect(B)
//
//     +-------+
//     |       |
//     |   A   |
//     |    +--+----+   =   +--+
//     +----+--+    |       +--+
//          |   B   |
//          |       |
//          +-------+
//
func (s *Solid) Intersect(other *Solid) *Solid {
	a, b := &BSP{}, &BSP{}
	a.AddPolygons(s.Polygons.Clone())
	b.AddPolygons(other.Polygons.Clone())
	a.Invert()
	b.ClipTo(a)
	b.Invert()
	a.ClipTo(b)
	b.ClipTo(a)
	a.AddPolygons(b.AllPolygons())
	a.Invert()
	return SolidFromPolygons(a.AllPolygons())
}

// Inverse returns a new CSG solid with solid and empty space switched.
// This solid is not modified.
func (s *Solid) Inverse() *Solid {
	polygons := s.Polygons.Clone()
	for i := range polygons {
		polygons[i].Flip()
	}
	return SolidFromPolygons(polygons)
}
