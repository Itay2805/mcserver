package chunk

type section struct {
	blockCount 	int
	ids        [4096]uint16
	palette 	*palette
}

func (c *section) generatePalette() *palette {
	if c.palette != nil {
		return c.palette
	}

	c.palette = &palette{}

	// set everything to unused
	for i := range c.palette.indexMap {
		c.palette.indexMap[i] = -1
	}

	// populate it
	for _, id := range c.ids {
		if c.palette.indexMap[id] == -1 {
			c.palette.indexMap[id] = len(c.palette.ids)
			c.palette.ids = append(c.palette.ids, id)
		}
	}

	c.palette.computeBitsPerBlock()

	return c.palette
}
