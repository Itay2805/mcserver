package chunk

import "math/bits"

type palette struct {
	ids          []uint16
	indexMap     [4096]int
	bitsPerBlock byte
}

func (p *palette) computeBitsPerBlock() {
	num := uint(len(p.ids))
	if num == 0 {
		p.bitsPerBlock = 4
	} else {
		res := bits.Len(num)
		if (num & (num - 1)) != 0 {
			res++
		}

		if res < 4 {
			res = 4
		}

		p.bitsPerBlock = byte(res)
	}
}


