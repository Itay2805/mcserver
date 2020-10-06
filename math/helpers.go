package math

// Intersect computes the intersection of two rectangles.  If no intersection
// exists, the intersection is nil.
func Intersect(r1, r2 *Rect) bool {
	// There are four cases of overlap:
	//
	//     1.  a1------------b1
	//              a2------------b2
	//              p--------q
	//
	//     2.       a1------------b1
	//         a2------------b2
	//              p--------q
	//
	//     3.  a1-----------------b1
	//              a2-------b2
	//              p--------q
	//
	//     4.       a1-------b1
	//         a2-----------------b2
	//              p--------q
	//
	// Thus there are only two cases of non-overlap:
	//
	//     1. a1------b1
	//                    a2------b2
	//
	//     2.             a1------b1
	//        a2------b2
	//
	// Enforced by constructor: a1 <= b1 and a2 <= b2.  So we can just
	// check the endpoints.

	for i := range r1.p {
		a1, b1, a2, b2 := r1.p[i], r1.q[i], r2.p[i], r2.q[i]
		if b2 <= a1 || b1 <= a2 {
			return false
		}
	}
	return true
}

// BoundingBox constructs the smallest rectangle containing both r1 and r2.
func BoundingBox(r1, r2 *Rect) (bb *Rect) {
	bb = new(Rect)
	dim := len(r1.p)
	bb.p = [3]float64{}
	bb.q = [3]float64{}
	for i := 0; i < dim; i++ {
		if r1.p[i] <= r2.p[i] {
			bb.p[i] = r1.p[i]
		} else {
			bb.p[i] = r2.p[i]
		}
		if r1.q[i] <= r2.q[i] {
			bb.q[i] = r2.q[i]
		} else {
			bb.q[i] = r1.q[i]
		}
	}
	return
}

// BoundingBoxN constructs the smallest rectangle containing all of r...
func BoundingBoxN(rects ...*Rect) (bb *Rect) {
	if len(rects) == 1 {
		bb = rects[0]
		return
	}
	bb = BoundingBox(rects[0], rects[1])
	for _, rect := range rects[2:] {
		bb = BoundingBox(bb, rect)
	}
	return
}
