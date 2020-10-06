package math

import "math"

// Point represents a point in n-dimensional Euclidean space.
type Point [3]float64

func NewPoint(x, y, z float64) Point {
	return Point{x, y, z}
}

func (p Point) X() float64 {
	return p[0]
}

func (p Point) Y() float64 {
	return p[1]
}

func (p Point) Z() float64 {
	return p[2]
}

// Dist computes the Euclidean distance between two points p and q.
func (p Point) Dist(q Point) float64 {
	sum := 0.0
	for i := range p {
		dx := p[i] - q[i]
		sum += dx * dx
	}
	return math.Sqrt(sum)
}

// MinDist computes the square of the distance from a point to a rectangle.
// If the point is contained in the rectangle then the distance is zero.
//
// Implemented per Definition 2 of "Nearest Neighbor Queries" by
// N. Roussopoulos, S. Kelley and F. Vincent, ACM SIGMOD, pages 71-79, 1995.
func (p Point) MinDist(r *Rect) float64 {
	sum := 0.0
	for i, pi := range p {
		if pi < r.p[i] {
			d := pi - r.p[i]
			sum += d * d
		} else if pi > r.q[i] {
			d := pi - r.q[i]
			sum += d * d
		} else {
			sum += 0
		}
	}
	return sum
}


// MinMaxDist computes the minimum of the maximum distances from p to points
// on r.  If r is the bounding box of some geometric objects, then there is
// at least one object contained in r within MinMaxDist(p, r) of p.
//
// Implemented per Definition 4 of "Nearest Neighbor Queries" by
// N. Roussopoulos, S. Kelley and F. Vincent, ACM SIGMOD, pages 71-79, 1995.
func (p Point) MinMaxDist(r *Rect) float64 {
	// by definition, MinMaxDist(p, r) =
	// min{1<=k<=n}(|pk - rmk|^2 + sum{1<=i<=n, i != k}(|pi - rMi|^2))
	// where rmk and rMk are defined as follows:

	rm := func(k int) float64 {
		if p[k] <= (r.p[k]+r.q[k])/2 {
			return r.p[k]
		}
		return r.q[k]
	}

	rM := func(k int) float64 {
		if p[k] >= (r.p[k]+r.q[k])/2 {
			return r.p[k]
		}
		return r.q[k]
	}

	// This formula can be computed in linear time by precomputing
	// S = sum{1<=i<=n}(|pi - rMi|^2).

	S := 0.0
	for i := range p {
		d := p[i] - rM(i)
		S += d * d
	}

	// Compute MinMaxDist using the precomputed S.
	min := math.MaxFloat64
	for k := range p {
		d1 := p[k] - rM(k)
		d2 := p[k] - rm(k)
		d := S - d1*d1 + d2*d2
		if d < min {
			min = d
		}
	}

	return min
}

// ToRect constructs a rectangle containing p with side lengths 2*tol.
func (p Point) ToRect(tol float64) *Rect {
	a, b := [3]float64{}, [3]float64{}
	for i := range p {
		a[i] = p[i] - tol
		b[i] = p[i] + tol
	}
	return &Rect{a, b}
}
