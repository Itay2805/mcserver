package flatgrass

import (
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/block"
	"github.com/itay2805/mcserver/minecraft/chunk"
)

type FlatgraassGenerator struct {
	stone		*block.Block
	grass		*block.Block
	dirt		*block.Block
}

func (d *FlatgraassGenerator) GetSpawn() minecraft.Position {
	return minecraft.Position{ X: 8, Y: 62, Z: 8 }
}

func (d *FlatgraassGenerator) GenerateChunk(x, z int) *chunk.Chunk {
	c := chunk.NewChunk(x, z)

	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			c.SetBlockState(x, 0, z, block.Bedrock.DefaultStateId)

			for y := 1; y < 48; y++ {
				c.SetBlockState(x, y, z, block.Stone.DefaultStateId)
			}

			for y := 48; y < 61; y++ {
				c.SetBlockState(x, y, z, block.Dirt.DefaultStateId)
			}

			c.SetBlockState(x, 61, z, block.GrassBlock.DefaultStateId)
		}
	}

	return c
}