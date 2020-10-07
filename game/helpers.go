package game

func forEachChunkInRange(x0, y0, radius int, cb func(x, z int)) {
	for x := x0 - radius - 2; x < x0 + radius + 2; x++ {
		for z := y0 - radius - 2; z < y0 + radius + 2; z++ {
			cb(x, z)
		}
	}
}
