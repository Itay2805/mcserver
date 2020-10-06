package common

type CompactArray struct {
	Values         []int64
	BitsPerElement int
	Max            int
}

func alignUp(x, alignTo int) int {
	return (x + (alignTo - 1)) &^ (alignTo - 1)
}

func CompactArrayLength(bits, count int) int {
	return alignUp(bits * count, 64) / 64
}

func NewCompactArray(bits, count int) *CompactArray {
	return &CompactArray{
		Values:         make([]int64, CompactArrayLength(bits, count)),
		BitsPerElement: bits,
		Max:            (1 << bits) - 1,
	}
}

func (c *CompactArray) Set(index, value int) int {
	// calculate the different indexes
	bIndex := index * c.BitsPerElement
	sIndex := bIndex >> 0x06
	eIndex := (((index + 1) * c.BitsPerElement) - 1) >> 0x06
	uIndex := bIndex ^ (sIndex << 0x06)

	// take the prev value
	previousValue := int64(uint64(c.Values[sIndex]) >> uIndex) & int64(c.Max)

	// set the new value
	c.Values[sIndex] = c.Values[sIndex] & int64(^(c.Max << uIndex)) | int64((value&c.Max) << uIndex)

	// if we cross boundries then handle it
	if sIndex != eIndex {
		zIndex := 64 - uIndex
		pIndex := c.BitsPerElement - 1

		// prev value...
		previousValue |= (c.Values[eIndex] << zIndex) & int64(c.Max)

		// new value...
		c.Values[eIndex] = int64(((uint64(c.Values[eIndex]) >> pIndex) << pIndex) | uint64((value&c.Max) >> zIndex))
	}

	return int(previousValue)
}

func (c *CompactArray) Get(index int) int {
	bIndex := index * c.BitsPerElement
	sIndex := bIndex >> 0x06
	eIndex := (((index + 1) * c.BitsPerElement) - 1) >> 0x06
	uIndex := bIndex ^ (sIndex << 0x06)

	if sIndex == eIndex {
		return int((uint64(c.Values[sIndex]) >> uIndex) & uint64(c.Max))
	}

	zIndex := 64 - uIndex

	return int(uint64(c.Values[sIndex] >> uIndex) | uint64(c.Values[eIndex] << zIndex) & uint64(c.Max))
}
