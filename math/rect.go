package math

import (
	"fmt"
	"log"
	"math"
	"strings"
)

type DistError float64

func (err DistError) Error() string {
	return "rtree: improper distance"
}

// Rect represents a subset of n-dimensional Euclidean space of the form
// [a1, b1] x [a2, b2] x ... x [an, bn], where ai < bi for all 1 <= i <= n.
type Rect struct {
	p, q Point // Enforced by NewRect: p[i] <= q[i] for all i.
}

// PointCoord returns the coordinate of the point of the rectangle at i
func (r *Rect) PointCoord(i int) float64 {
	return r.p[i]
}

// LengthsCoord returns the coordinate of the lengths of the rectangle at i
func (r *Rect) LengthsCoord(i int) float64 {
	return r.q[i] - r.p[i]
}

// Equal returns true if the two rectangles are equal
func (r *Rect) Equal(other *Rect) bool {
	for i, e := range r.p {
		if e != other.p[i] {
			return false
		}
	}
	for i, e := range r.q {
		if e != other.q[i] {
			return false
		}
	}
	return true
}

func (r *Rect) String() string {
	s := make([]string, len(r.p))
	for i, a := range r.p {
		b := r.q[i]
		s[i] = fmt.Sprintf("[%.2f, %.2f]", a, b)
	}
	return strings.Join(s, "x")
}

// NewRect constructs and returns a pointer to a Rect given a corner point and
// the lengths of each dimension.  The point p should be the most-negative point
// on the rectangle (in every dimension) and every length should be positive.
func NewRect(p Point, lengths [3]float64) (r *Rect) {
	r = new(Rect)
	r.p = p
	r.q = [3]float64{}
	for i := range p {
		if lengths[i] <= 0 {
			log.Panicln(DistError(lengths[i]))
		}
		r.q[i] = p[i] + lengths[i]
	}
	return
}

// NewRectFromPoints constructs and returns a pointer to a Rect given a corner points.
func NewRectFromPoints(minPoint, maxPoint Point) *Rect {
	//checking that  min and max points is swapping
	for i, p := range minPoint {
		if minPoint[i] > maxPoint[i] {
			minPoint[i] = maxPoint[i]
			maxPoint[i] = p
		}
	}

	return &Rect{p: minPoint, q: maxPoint}
}

// Size computes the measure of a rectangle (the product of its side lengths).
func (r *Rect) Size() float64 {
	size := 1.0
	for i, a := range r.p {
		b := r.q[i]
		size *= b - a
	}
	return size
}

// Margin computes the sum of the edge lengths of a rectangle.
func (r *Rect) Margin() float64 {
	// The number of edges in an n-dimensional rectangle is n * 2^(n-1)
	// (http://en.wikipedia.org/wiki/Hypercube_graph).  Thus the number
	// of edges of length (ai - bi), where the rectangle is determined
	// by p = (a1, a2, ..., an) and q = (b1, b2, ..., bn), is 2^(n-1).
	//
	// The margin of the rectangle, then, is given by the formula
	// 2^(n-1) * [(b1 - a1) + (b2 - a2) + ... + (bn - an)].
	dim := len(r.p)
	sum := 0.0
	for i, a := range r.p {
		b := r.q[i]
		sum += b - a
	}
	return math.Pow(2, float64(dim-1)) * sum
}

// ContainsPoint tests whether p is located inside or on the boundary of r.
func (r *Rect) ContainsPoint(p Point) bool {
	for i, a := range p {
		// p is contained in (or on) r if and only if p <= a <= q for
		// every dimension.
		if a < r.p[i] || a > r.q[i] {
			return false
		}
	}

	return true
}

// ContainsRect tests whether r2 is is located inside r1.
func (r *Rect) ContainsRect(r2 *Rect) bool {
	for i, a1 := range r.p {
		b1, a2, b2 := r.q[i], r2.p[i], r2.q[i]
		// enforced by constructor: a1 <= b1 and a2 <= b2.
		// so containment holds if and only if a1 <= a2 <= b2 <= b1
		// for every dimension.
		if a1 > a2 || b2 > b1 {
			return false
		}
	}

	return true
}
