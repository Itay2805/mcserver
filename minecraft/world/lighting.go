package world

import (
	"github.com/eapache/queue"
	"github.com/itay2805/mcserver/minecraft/block"
	"github.com/itay2805/mcserver/minecraft/chunk"
)

// the light updates that need to be processed
var blockLightUpdates = queue.New()
var skyLightUpdates = queue.New()

type lightUpdate struct {
	x, y, z 	int
	chunk 		*chunk.Chunk
	world		*World
}

// queue a light update to happen later on
func QueueLightUpdate(w *World, c *chunk.Chunk, x, y, z int) {
	blockLightUpdates.Add(lightUpdate{
		x:     x,
		y:     y,
		z:     z,
		chunk: c,
		world: w,
	})
}

// process a light update
func ProcessLightUpdates() {
	for blockLightUpdates.Length() > 0 {
		update := blockLightUpdates.Remove().(lightUpdate)
		updateBlockLight(blockLightUpdates, update.world, update.chunk, update.x, update.y, update.z)
	}

	for skyLightUpdates.Length() > 0 {
		update := skyLightUpdates.Remove().(lightUpdate)
		updateSkylight(skyLightUpdates, update.world, update.chunk, update.x, update.y, update.z)
	}
}

// light up a single chunk, used when generating new chunks
func (w *World) lightChunk(c *chunk.Chunk) {
	blockLightUpdates := queue.New()
	skyLightUpdates := queue.New()

	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			c.SetSkyLight(x, 255, z, 15)

			currFilter := 15
			for y := 254; y >= 0; y-- {
				typ := block.GetByStateId(c.GetBlockState(x, y, z))

				if currFilter > 0 {
					currFilter -= typ.FilterLight
					if currFilter < 0 {
						currFilter = 0
					}
				}

				if !typ.Transparent {
					if typ.EmitLight > 0 {
						blockLightUpdates.Add(lightUpdate{
							x:     x,
							y:     y,
							z:     z,
							chunk: c,
							world: w,
						})
					}
				} else {
					skyLightUpdates.Add(lightUpdate{
						x:     x,
						y:     y,
						z:     z,
						chunk: c,
						world: w,
					})
				}

				c.SetSkyLight(x, y, z, currFilter)
			}
		}
	}

	for blockLightUpdates.Length() > 0 {
		update := blockLightUpdates.Remove().(lightUpdate)
		updateBlockLight(blockLightUpdates, update.world, update.chunk, update.x, update.y, update.z)
	}

	for skyLightUpdates.Length() > 0 {
		update := skyLightUpdates.Remove().(lightUpdate)
		updateSkylight(skyLightUpdates, update.world, update.chunk, update.x, update.y, update.z)
	}
}

// return the max of all nums
func max(nums ...int) int {
	largest := nums[0]
	for _, num := range nums {
		if num > largest {
			largest = num
		}
	}
	return largest
}

// This function will do a sky light update from the given block in
// a flood fill manner
func updateSkylight(q *queue.Queue, w *World, c *chunk.Chunk, x, y, z int) {
	typ := block.GetByStateId(c.GetBlockState(x, y, z))
	light := 0

	if typ.FilterLight == 15 {
		light = 0
	} else {
		var sle int
		var slw int
		var slu int
		var sld int
		var sls int
		var sln int

		if x < 15 {
			sle = c.GetSkyLight(x + 1, y, z)
		} else {
			eastC := w.GetChunk(c.X + 1, c.Z)
			if eastC == nil {
				sle = 0
			} else {
				sle = eastC.GetSkyLight(0, y, z)
			}
		}

		if x > 0 {
			slw = c.GetSkyLight(x - 1, y, z)
		} else {
			westC := w.GetChunk(c.X - 1, c.Z)
			if westC == nil {
				slw = 0
			} else {
				slw = westC.GetSkyLight(15, y, z)
			}
		}

		if y < 255 {
			slu = c.GetSkyLight(x, y + 1, z)
		} else {
			slu = 15
		}

		if y > 0 {
			sld = c.GetSkyLight(x, y - 1, z)
		} else {
			sld = 0
		}

		if z < 15 {
			sls = c.GetSkyLight(x, y, z + 1)
		} else {
			southC := w.GetChunk(c.X, c.Z + 1)
			if southC == nil {
				sls = 0
			} else {
				sls = southC.GetSkyLight(x, y, 0)
			}
		}

		if z > 0 {
			sln = c.GetSkyLight(x, y, z - 1)
		} else {
			northC := w.GetChunk(c.X, c.Z - 1)
			if northC == nil {
				sln = 0
			} else {
				sln = northC.GetSkyLight(x, y, 15)
			}
		}

		brightest := max(sle, slw, slu, sld, sls, sln, 0)
		light = brightest - typ.FilterLight - 1
		if light < 0 {
			light = 0
		}
	}

	if c.GetSkyLight(x, y, z) != light {
		c.SetSkyLight(x, y, z, light)

		if x < 15 {
			q.Add(lightUpdate{
				x:     x + 1,
				y:     y,
				z:     z,
				chunk: c,
				world: w,
			})
		}

		if x > 0 {
			q.Add(lightUpdate{
				x:     x - 1,
				y:     y,
				z:     z,
				chunk: c,
				world: w,
			})
		}

		if y < 255 {
			q.Add(lightUpdate{
				x:     x,
				y:     y + 1,
				z:     z,
				chunk: c,
				world: w,
			})
		}

		if y > 0 {
			q.Add(lightUpdate{
				x:     x,
				y:     y - 1,
				z:     z,
				chunk: c,
				world: w,
			})

		}

		if z < 15 {
			q.Add(lightUpdate{
				x:     x,
				y:     y,
				z:     z + 1,
				chunk: c,
				world: w,
			})

		}

		if z > 0 {
			q.Add(lightUpdate{
				x:     x,
				y:     y,
				z:     z - 1,
				chunk: c,
				world: w,
			})
		}
	}
}

func updateBlockLight(q *queue.Queue, w *World, c *chunk.Chunk, x, y, z int) {
	typ := block.GetByStateId(c.GetBlockState(x, y, z))

	light := 0

	if typ.FilterLight != 0 {
		light = typ.EmitLight
	} else {
		var ble int
		var blw int
		var blu int
		var bld int
		var bls int
		var bln int

		if x < 15 {
			ble = c.GetBlockLight(x + 1, y, z)
		} else {
			eastC := w.GetChunk(c.X + 1, c.Z)
			if eastC == nil {
				ble = 0
			} else {
				ble = eastC.GetBlockLight(0, y, z)
			}
		}

		if x > 0 {
			blw = c.GetBlockLight(x - 1, y, z)
		} else {
			westC := w.GetChunk(c.X - 1, c.Z)
			if westC == nil {
				blw = 0
			} else {
				blw = westC.GetBlockLight(15, y, z)
			}
		}

		if y < 255 {
			blu = c.GetBlockLight(x, y + 1, z)
		} else {
			blu = 15
		}

		if y > 0 {
			bld = c.GetBlockLight(x, y - 1, z)
		} else {
			bld = 0
		}

		if z < 15 {
			bls = c.GetBlockLight(x, y, z + 1)
		} else {
			southC := w.GetChunk(c.X, c.Z + 1)
			if southC == nil {
				bls = 0
			} else {
				bls = southC.GetBlockLight(x, y, 0)
			}
		}

		if z > 0 {
			bln = c.GetBlockLight(x, y, z - 1)
		} else {
			northC := w.GetChunk(c.X, c.Z - 1)
			if northC == nil {
				bln = 0
			} else {
				bln = northC.GetBlockLight(x, y, 15)
			}
		}

		brightest := max(ble, blw, blu, bld, bls, bln, 0)
		light = brightest - 1 + typ.EmitLight
		if light < 0 {
			light = 0
		} else if light > 15 {
			light = 15
		}
	}

	if c.GetBlockLight(x, y, z) != light {
		c.SetBlockLight(x, y, z, light)

		if x < 15 {
			q.Add(lightUpdate{
				x:     x + 1,
				y:     y,
				z:     z,
				chunk: c,
				world: w,
			})
		}

		if x > 0 {
			q.Add(lightUpdate{
				x:     x - 1,
				y:     y,
				z:     z,
				chunk: c,
				world: w,
			})
		}

		if y < 255 {
			q.Add(lightUpdate{
				x:     x,
				y:     y + 1,
				z:     z,
				chunk: c,
				world: w,
			})
		}

		if y > 0 {
			q.Add(lightUpdate{
				x:     x,
				y:     y - 1,
				z:     z,
				chunk: c,
				world: w,
			})

		}

		if z < 15 {
			q.Add(lightUpdate{
				x:     x,
				y:     y,
				z:     z + 1,
				chunk: c,
				world: w,
			})

		}

		if z > 0 {
			q.Add(lightUpdate{
				x:     x,
				y:     y,
				z:     z - 1,
				chunk: c,
				world: w,
			})
		}
	}
}
