# csg

    import "github.com/reactivego/csg"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/csg.svg)](https://pkg.go.dev/github.com/reactivego/csg#section-documentation)

![Example](../assets/example.svg)

Constructive Solid Geometry (CSG) is a modeling technique that uses Boolean operations like union and intersection to combine 3D solids. This library implements CSG operations on meshes elegantly and concisely using BSP trees, and is meant to serve as an easily understandable implementation of the algorithm. All edge cases involving overlapping coplanar polygons in both solids are correctly handled.

Example usage:
``` go
import "github.com/reactivego/csg"

cube := csg.Cube()
sphere := csg.Sphere(csg.Radius(1.3))
polygons := cube.Subtract(sphere).Polygons
```

## Operations

This library provides three CSG operations: union, subtract, and intersect.
The operations are rendered below.

| ![Cube](../assets/cube.svg) | ![Sphere](../assets/sphere.svg) |
|:---:|:---:|
| `a` | `b` |


The solids `a` and `b` were generated with the following code:

```go
import . "github.com/reactivego/csg"

a := Cube(Center(-0.25, -0.25, -0.25))
b := Sphere(Center(0.25, 0.25, 0.25), Radius(1.3))
```

| ![Union](../assets/union.svg) | ![Subtract](../assets/subtract.svg) | ![Intersect](../assets/intersect.svg) |
|:---:|:---:|:---:|
| `a.Union(b)`| `a.Subtract(b)` | `a.Intersect(b)` |

## Combined CSG Example

Below is a solid constructed from a combination of operations:

| ![Cube](../assets/a.svg) | ![Sphere](../assets/b.svg) | ![Cylinder X](../assets/c.svg) | ![Cylinder Y](../assets/d.svg) | ![Cylinder Z](../assets/e.svg) |
|:---:|:---:|:---:|:---:|:---:|
| `a` | `b` | `c` | `d` | `e` |

The solids above were generated with the following code:

```go
import . "github.com/reactivego/csg"

a := Cube()
b := Sphere(Radius(1.35), Stacks(12))
c := Cylinder(Radius(0.7), Start(-1, 0, 0), End(1, 0, 0))
d := Cylinder(Radius(0.7), Start(0, -1, 0), End(0, 1, 0))
e := Cylinder(Radius(0.7), Start(0, 0, -1), End(0, 0, 1))
```

| ![Combined](../assets/combined.svg) |
|:---:|
| `a.Intersect(b).Subtract(c.Union(d).Union(e))` |

The combined solid was generated with the code:
```go
a.Intersect(b).Subtract(c.Union(d).Union(e))
```

## Implementation Details

All CSG operations are implemented in terms of two functions, `ClipTo()` and
`Invert()`, which remove parts of a BSP tree inside another BSP tree and swap
solid and empty space, respectively. To find the union of `a` and `b`, we
want to remove everything in `a` inside `b` and everything in `b` inside `a`,
then combine polygons from `a` and `b` into one solid:

``` go
a.ClipTo(b)
b.ClipTo(a)
a.AddPolygons(b.AllPolygons())
```

The only tricky part is handling overlapping coplanar polygons in both trees.
The code above keeps both copies, but we need to keep them in one tree and
remove them in the other tree. To remove them from `b` we can clip the
inverse of `b` against `a`. The code for union now looks like this:

``` go
a.ClipTo(b)
b.ClipTo(a)
b.Invert()
b.ClipTo(a)
b.Invert()
a.AddPolygons(b.AllPolygons())
```

Subtraction and intersection naturally follow from set operations. If
union is `A | B`, subtraction is `A - B = ~(~A | B)` and intersection is
`A & B = ~(~A | ~B)` where `~` is the complement operator.

## Acknowledgments

This package is a direct conversion of [cgs.js](http://evanw.github.io/csg.js/) to Go.
The original JavaScript code was written by Evan Wallace and committed to git on November 30, 2011.
As an odd coincidence, I am writing this a full 10 years later on Nov 30, 2021.

## License

This library is licensed under the terms of the MIT License.
See [LICENSE](LICENSE) file for copyright notice and exact wording.