package csg

import (
	"fmt"
	"strings"
)

// BSP holds a node in a BSP tree. A BSP tree is built from a collection of
// polygons by picking a polygon to split along. That polygon (and all other
// coplanar polygons) are added directly to that node and the other polygons
// are added to the front and/or back subtrees. This is not a leafy BSP tree
// since there is no distinction between internal and leaf nodes.
type BSP struct {
	Plane    *Plane
	Polygons Polygons
	Front    *BSP
	Back     *BSP
}

// Invert converts solid space to empty space and empty space to solid space.
func (n *BSP) Invert() {
	for i := range n.Polygons {
		n.Polygons[i].Flip()
	}
	n.Plane.Flip()
	if n.Front != nil {
		n.Front.Invert()
	}
	if n.Back != nil {
		n.Back.Invert()
	}
	n.Front, n.Back = n.Back, n.Front
}

// ClipPolygons recursively removes all polygons in `polygons` that are inside
// this BSP tree.
func (n BSP) ClipPolygons(polygons Polygons) Polygons {
	if n.Plane == nil {
		return append(Polygons(nil), polygons...)
	}
	var front, back Polygons
	for _, p := range polygons {
		n.Plane.SplitPolygon(p, &front, &back, &front, &back)
	}
	if n.Front != nil {
		front = n.Front.ClipPolygons(front)
	}
	if n.Back != nil {
		back = n.Back.ClipPolygons(back)
	} else {
		back = nil
	}
	return append(front, back...)
}

// ClipTo removes all polygons in this BSP tree that are inside the other BSP
// tree `bsp`.
func (n *BSP) ClipTo(bsp *BSP) {
	n.Polygons = bsp.ClipPolygons(n.Polygons)
	if n.Front != nil {
		n.Front.ClipTo(bsp)
	}
	if n.Back != nil {
		n.Back.ClipTo(bsp)
	}
}

// AllPolygons returns a list of all polygons in this BSP tree.
func (n BSP) AllPolygons() Polygons {
	polygons := append(Polygons(nil), n.Polygons...)
	if n.Front != nil {
		polygons = append(polygons, n.Front.AllPolygons()...)
	}
	if n.Back != nil {
		polygons = append(polygons, n.Back.AllPolygons()...)
	}
	return polygons
}

// AddPolygons builds a BSP tree out of `polygons`. When called on an existing
// tree, the new polygons are filtered down to the bottom of the tree and become
// new nodes there. Each set of polygons is partitioned using the first polygon
// (no heuristic is used to pick a good split).
func (n *BSP) AddPolygons(polygons Polygons) {
	if len(polygons) == 0 {
		return
	}
	if n.Plane == nil {
		p := polygons[0].Plane
		n.Plane = &p
	}
	var front, back Polygons
	for _, p := range polygons {
		n.Plane.SplitPolygon(p, &n.Polygons, &n.Polygons, &front, &back)
	}
	if len(front) > 0 {
		if n.Front == nil {
			n.Front = &BSP{}
		}
		n.Front.AddPolygons(front)
	}
	if len(back) > 0 {
		if n.Back == nil {
			n.Back = &BSP{}
		}
		n.Back.AddPolygons(back)
	}
}

func (n *BSP) print(level int, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("%*s%s:%+v\n", level*2, "", "plane", n.Plane.Normal))
	if n.Front != nil {
		n.Front.print(level+1, sb)
	}
	if n.Back != nil {
		n.Back.print(level+1, sb)
	}
}

func (n *BSP) String() string {
	var sb strings.Builder
	n.print(0, &sb)
	return sb.String()
}
