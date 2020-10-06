package game

func forEachChunkInRange(x0, y0, radius int, cb func(x, z int)) {
	radius++
	x := radius
	y := 0
	xChange := 1 - (radius << 1)
	yChange := 0
	radiusError := 0

	for x >= y {
		for i := x0 - x; i <= x0 + x; i++ {
			cb(i, y0 + y)
			cb(i, y0 - y)
		}

		for i := x0 - y; i <= x0 + y; i++ {
			cb(i, y0 + x)
			cb(i, y0 - x)
		}

		y++
		radiusError += yChange
		yChange += 2
		if ((radiusError << 1) + xChange) > 0 {
			x--
			radiusError += xChange
			xChange += 2
		}
	}
}
