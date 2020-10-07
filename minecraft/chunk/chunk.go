package chunk

import (
	"github.com/itay2805/mcserver/common"
	"github.com/itay2805/mcserver/minecraft"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper functions to get the index to certain arrays in the chunk
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func getSectionIndex(y int) int {
	return y >> 4
}

func getLightSection(y int) int {
	return (y >> 4) + 1
}

func getBiomeIndex(x, y, z int) int {
	return ((y >> 2) & 63) << 4 | ((z >> 2) & 3) << 2 | ((x >> 2) & 3)
}

func getSectionBlockIndex(x, y, z int) int {
	return ((y & 0xf) << 8) | (z << 4) | x
}

func getHeightMapIndex(x, z int) int {
	return x * 16 + z
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// The chunk itself
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Chunk struct {
	X, Z 	int

	// block related data
	sections			[NumSections]*section

	// biome related data
	biomes				[4 * 4 * 64]int32

	// light related data
	skyLightSections	[NumLightSections]*[LightVolume]uint8
	blockLightSections	[NumLightSections]*[LightVolume]uint8
}

func NewChunk(x, z int) *Chunk {
	return &Chunk{
		X:                  x,
		Z:                  z,
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helpers to get values from the chunk
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Chunk) GetBlockState(x, y, z int) uint16 {
	sec := c.sections[getSectionIndex(y)]
	if sec == nil {
		return 0
	}

	return sec.ids[getSectionBlockIndex(x, y, z)]
}

func (c *Chunk) SetBlockState(x, y, z int, state uint16) {
	sec := c.sections[getSectionIndex(y)]
	if sec == nil {
		// don't insert a new section if this is an air block
		if state == 0 {
			return
		}
		sec = &section{}
		c.sections[getSectionIndex(y)] = sec
	}

	// get the old state for updating the block count
	// NOTE: This assumes air is always 0
	oldState := sec.ids[getSectionBlockIndex(x, y, z)]

	// set the block count
	if oldState != 0 {
		if state == 0 {
			// was a block, not anymore
			sec.blockCount--
		}
	} else if state != 0 {
		// was air, not anymore
		sec.blockCount++
	}

	// set the block as needed
	if sec.blockCount == 0 {
		// no more blocks, remove section
		c.sections[getSectionIndex(y)] = nil
	} else {
		// set the block
		sec.ids[getSectionBlockIndex(x, y, z)] = state
		sec.palette = nil
	}
}

func (c *Chunk) GetSkyLight(x, y, z int) int {
	sec := c.skyLightSections[getLightSection(y)]
	if sec == nil {
		return 0
	}

	idx := getSectionBlockIndex(x, y, z)
	half := idx >> 1

	if (idx & 1) == 1 {
		return int(sec[half] >> 4)
	} else {
		return int(sec[half] & 0xf)
	}
}

func (c *Chunk) SetSkyLight(x, y, z int, light int) {
	sec := c.skyLightSections[getLightSection(y)]
	if sec == nil {
		if light == 0 {
			// ignore if setting light to 0
			return
		}
		sec = &[LightVolume]uint8{}
		c.skyLightSections[getLightSection(y)] = sec
	}

	val := uint8(light & 0xf)
	idx := getSectionBlockIndex(x, y, z)
	half := idx >> 1

	if (idx & 1) == 1 {
		sec[half] = (sec[half] & 0x0f) | (val << 4)
	} else {
		sec[half] = (sec[half] & 0xf0) | val
	}
}

func (c *Chunk) GetBlockLight(x, y, z int) int {
	sec := c.blockLightSections[getLightSection(y)]
	if sec == nil {
		return 0
	}

	idx := getSectionBlockIndex(x, y, z)
	half := idx >> 1

	if (idx & 1) == 1 {
		return int(sec[half] >> 4)
	} else {
		return int(sec[half] & 0xf)
	}
}

func (c *Chunk) SetBlockLight(x, y, z int, light int) {
	sec := c.blockLightSections[getLightSection(y)]
	if sec == nil {
		if light == 0 {
			// ignore if setting light to 0
			return
		}
		sec = &[LightVolume]uint8{}
		c.blockLightSections[getLightSection(y)] = sec
	}

	val := uint8(light & 0xf)
	idx := getSectionBlockIndex(x, y, z)
	half := idx >> 1

	if (idx & 1) == 1 {
		sec[half] = (sec[half] & 0x0f) | (val << 4)
	} else {
		sec[half] = (sec[half] & 0xf0) | val
	}
}

//
// Write the chunk data into the given writer in the Chunk Data format
//
// NOTE: This does not include block entities! these are managed by the world
//		 so they need to be written by the world as well!
//
func (c *Chunk) MakeChunkDataPacket(writer *minecraft.Writer) {
	// prepare the data
	primaryBitMask := int32(0)
	dataSize := 0
	for i, sec := range c.sections {
		if sec == nil {
			continue
		}

		// we have data
		primaryBitMask |= 1 << i

		// generate the palette
		// no need to save because this is cached
		palette := sec.generatePalette()

		dataSize += 2 // block count
		dataSize += 1 // bits per block

		// if this is less or equals to 8 then
		// account for the palette
		if palette.bitsPerBlock <= MaxBitsPerBlock {
			// the size of the number of blocks in the palette
			dataSize += common.VarintSize(int32(len(palette.ids)))

			// the size of each if in the palette
			for _, id := range palette.ids {
				dataSize += common.VarintSize(int32(id))
			}
		}

		// data array length + data sizes
		arrayLen := common.CompactArrayLength(int(palette.bitsPerBlock), 4096)
		dataSize += common.VarintSize(int32(arrayLen))
		dataSize += arrayLen * 8
	}

	// grow to fit the whole packet
	writer.Grow(common.VarintSize(0x22) +
				4 + 4 + 1 + 4096 +
				common.VarintSize(primaryBitMask) +
				common.VarintSize(int32(dataSize)) +
				dataSize)

	// write the header
	writer.WriteVarint(0x22)
	writer.WriteInt(int32(c.X))
	writer.WriteInt(int32(c.Z))
	writer.WriteBoolean(true)

	// write the headers
	writer.WriteVarint(primaryBitMask)

	// we don't send height maps to the client
	nbt := writer.StartNBT()
	nbt.StartCompound("")
	nbt.EndCompound()

	// Write the biomes
	for _, biome := range c.biomes {
		writer.WriteInt(biome)
	}

	// write the chunk data
	writer.WriteVarint(int32(dataSize))
	for _, sec := range c.sections {
		if sec == nil {
			continue
		}

		// generate the palette
		// no need to save because this is cached
		palette := sec.generatePalette()

		writer.WriteShort(int16(sec.blockCount))
		writer.WriteByte(palette.bitsPerBlock)

		// if this is less or equals to 8 then
		// account for the palette
		direct := palette.bitsPerBlock > MaxBitsPerBlock
		if !direct {
			// the size of the number of blocks in the palette
			writer.WriteVarint(int32(len(palette.ids)))

			// the size of each if in the palette
			for _, id := range palette.ids {
				writer.WriteVarint(int32(id))
			}
		}

		// prepare the array
		packed := common.NewCompactArray(int(palette.bitsPerBlock), 4096)
		if direct {
			// for direct just pack all the items
			for i, id := range sec.ids {
				packed.Set(i, int(id))
			}
		} else {
			// for indirect use the map to figure
			// the id to put
			for i, id := range sec.ids {
				packed.Set(i, palette.indexMap[id])
			}
		}

		// data array length + data
		writer.WriteVarint(int32(len(packed.Values)))
		for _, l := range packed.Values {
			writer.WriteLong(l)
		}
	}
}

// TODO: maybe cache this
func (c *Chunk) MakeUpdateLightPacket(writer *minecraft.Writer) {
	// prepare the masks
	// NOTE: This assumes the amount of sky light sections is
	// 		 the same as block light sections, might need to be
	//		 changed in the future
	skyLightMask := int32(0)
	blockLightMask := int32(0)
	count := 0
	for i := 0; i < NumLightSections; i++ {
		if c.skyLightSections[i] != nil {
			skyLightMask |= 1 << i
			count++
		}

		if c.blockLightSections[i] != nil {
			blockLightMask |= 1 << i
			count++
		}
	}

	// we already know everything so just grow it properly
	writer.Grow(
		common.VarintSize(0x25) +
			common.VarintSize(int32(c.X)) +
			common.VarintSize(int32(c.Z)) +
			common.VarintSize(skyLightMask) + common.VarintSize(blockLightMask) +
			common.VarintSize(0) + common.VarintSize(0) +
			(common.VarintSize(2048) + 2048) * 18 * 2)

	// write the basic data
	writer.WriteVarint(0x25)
	writer.WriteVarint(int32(c.X))
	writer.WriteVarint(int32(c.Z))

	// write the masks
	writer.WriteVarint(skyLightMask)
	writer.WriteVarint(blockLightMask)
	writer.WriteVarint(0)
	writer.WriteVarint(0)

	// write out all of the sky light arrays
	for _, sec := range c.skyLightSections {
		if sec == nil {
			continue
		}

		writer.WriteVarint(int32(len(sec)))
		writer.WriteBytes(sec[:])
	}

	// write out all of the block light arrays
	for _, sec := range c.blockLightSections {
		if sec == nil {
			continue
		}

		writer.WriteVarint(int32(len(sec)))
		writer.WriteBytes(sec[:])
	}
}
