package common

func VarintSize(val int32) int {
	v := uint16(val)

	// TODO: faster implementation for this
	count := 0
	for {
		v >>= 7
		count++
		if v == 0 {
			break
		}
	}
	return count
}
