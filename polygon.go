package csg

// Polygon represents a convex polygon. The vertices used to initialize a
// polygon must be coplanar and form a convex loop.
//
// Each convex polygon has a `Shared` property, which is shared between all
// polygons that are clones of each other or were split from the same polygon.
// This can be used to define per-polygon properties (such as surface color).
type Polygon struct {
	Vertices []Vertex
	Plane    Plane
}

func PolygonFromVertices(vertices ...Vertex) Polygon {
	return Polygon{
		Vertices: vertices,
		Plane:    PlaneFromPoints(vertices[0].Pos, vertices[1].Pos, vertices[2].Pos),
	}
}

func (p *Polygon) Flip() {
	for i, j := 0, len(p.Vertices)-1; i < j; i, j = i+1, j-1 {
		p.Vertices[i], p.Vertices[j] = p.Vertices[j], p.Vertices[i]
	}
	for i := range p.Vertices {
		p.Vertices[i].Flip()
	}
	p.Plane.Flip()
}

type Polygons []Polygon

func (p Polygons) Clone() Polygons {
	p = append(Polygons(nil), p...)
	for i := range p {
		p[i].Vertices = append([]Vertex(nil), p[i].Vertices...)
	}
	return p
}
