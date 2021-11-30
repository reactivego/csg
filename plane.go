package csg

// Plane represents a plane in 3D space.
type Plane struct {
	Normal Vector
	W      float64
}

// PlaneEPSILON is the tolerance used by `SplitPolygon()` to decide if a
// point is on the plane.
const PlaneEPSILON = 1e-5

func PlaneFromPoints(a, b, c Vector) Plane {
	n := b.Minus(a).Cross(c.Minus(a)).Unit()
	return Plane{Normal: n, W: n.Dot(a)}
}

func (p *Plane) Flip() {
	p.Normal = p.Normal.Negated()
	p.W = -p.W
}

// SplitPolygon splits `polygon` by this plane if needed, then put the polygon
// or polygon fragments in the appropriate lists. Coplanar polygons go into
// either `coplanarFront` or `coplanarBack` depending on their orientation with
// respect to this plane. Polygons in front or in back of this plane go into
// either `front` or `back`.
func (p Plane) SplitPolygon(polygon Polygon, coplanarFront, coplanarBack, front, back *Polygons) {
	type Type int

	const (
		COPLANAR = Type(iota)
		FRONT
		BACK
		SPANNING
	)

	// Classify each point as well as the entire polygon into one of the above
	// four classes.
	var polygonType = COPLANAR
	var types []Type
	for _, v := range polygon.Vertices {
		t := p.Normal.Dot(v.Pos) - p.W
		vType := COPLANAR
		if t < -PlaneEPSILON {
			vType = BACK
		} else if t > PlaneEPSILON {
			vType = FRONT
		}
		polygonType |= vType
		types = append(types, vType)
	}

	// Put the polygon in the correct list, splitting it when necessary.
	switch polygonType {
	case COPLANAR:
		if p.Normal.Dot(polygon.Plane.Normal) > 0 {
			*coplanarFront = append(*coplanarFront, polygon)
		} else {
			*coplanarBack = append(*coplanarBack, polygon)
		}
	case FRONT:
		*front = append(*front, polygon)
	case BACK:
		*back = append(*back, polygon)
	case SPANNING:
		var f, b []Vertex
		for i := range polygon.Vertices {
			j := (i + 1) % len(polygon.Vertices)
			ti, tj := types[i], types[j]
			vi, vj := polygon.Vertices[i], polygon.Vertices[j]
			if ti != BACK {
				f = append(f, vi)
			}
			if ti != FRONT {
				b = append(b, vi)
			}
			if (ti | tj) == SPANNING {
				t := (p.W - p.Normal.Dot(vi.Pos)) / p.Normal.Dot(vj.Pos.Minus(vi.Pos))
				v := vi.Interpolate(vj, t)
				f = append(f, v)
				b = append(b, v)
			}
		}
		if len(f) >= 3 {
			*front = append(*front, Polygon{Vertices: f, Plane: polygon.Plane})
		}
		if len(b) >= 3 {
			*back = append(*back, Polygon{Vertices: b, Plane: polygon.Plane})
		}
	}
}
