package world

import (
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/chunk"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// World generation
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Generator interface {
	//
	// Return the world's spawn location
	//
	GetSpawn() minecraft.Position

	//
	// Generate a new chunk at the given
	// coordinates
	//
	GenerateChunk(x, z int) *chunk.Chunk
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// World providers
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Provider interface {
	//
	// Save the chunk to persistent storage
	//
	SaveChunk(chunk *chunk.Chunk)

	//
	// Load the chunk from persistent storage
	// if does not exists then will return nil
	//
	LoadChunk(x, z int) *chunk.Chunk
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// The world itself
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ChunkPos struct {
	X, Z int
}

type World struct {
	Provider 	Provider
	Generator 	Generator
	chunks 		sync.Map
}

func NewWorld(provider Provider, generator Generator) *World {
	return &World{
		Provider:  provider,
		Generator: generator,
		chunks:    sync.Map{},
	}
}

// get a chunk that is already loaded in memory
func (w *World) GetChunk(x, z int) *chunk.Chunk {
	if c, ok := w.chunks.Load(ChunkPos{ x, z }); ok {
		return c.(*chunk.Chunk)
	}
	return nil
}

// will load a new channel if not loaded already or generate
// a new one if does not exist at all
func (w *World) LoadChunk(x, z int) *chunk.Chunk {
	// try to get from memory
	c := w.GetChunk(x, z)
	if c != nil {
		return c
	}

	// try to load from file
	c = w.Provider.LoadChunk(x, z)
	if c != nil {
		w.chunks.Store(ChunkPos{ x, z }, c)
		return c
	}

	// generate a new one
	c = w.Generator.GenerateChunk(x, z)
	w.chunks.Store(ChunkPos{ x, z }, c)
	return c
}

func (w *World) GetBlockState(x, y, z int) uint16 {
	chunkX := x >> 4
	chunkZ := z >> 4
	c := w.LoadChunk(chunkX, chunkZ)
	return c.GetBlockState(x & 0xf, y, z & 0xf)
}

func (w *World) SetBlockState(x, y, z int, state uint16) {
	chunkX := x >> 4
	chunkZ := z >> 4
	c := w.LoadChunk(chunkX, chunkZ)
	c.SetBlockState(x & 0xf, y, z & 0xf, state)
}

func (w *World) GetSkyLight(x, y, z int) int {
	chunkX := x >> 4
	chunkZ := z >> 4
	c := w.LoadChunk(chunkX, chunkZ)
	return c.GetSkyLight(x & 0xf, y, z & 0xf)
}

func (w *World) SetSkyLight(x, y, z int, light int) {
	chunkX := x >> 4
	chunkZ := z >> 4
	c := w.LoadChunk(chunkX, chunkZ)
	c.SetSkyLight(x & 0xf, y, z & 0xf, light)
}

func (w *World) GetBlockLight(x, y, z int) int {
	chunkX := x >> 4
	chunkZ := z >> 4
	c := w.LoadChunk(chunkX, chunkZ)
	return c.GetBlockLight(x & 0xf, y, z & 0xf)
}

func (w *World) SetBlockLight(x, y, z int, light int) {
	chunkX := x >> 4
	chunkZ := z >> 4
	c := w.LoadChunk(chunkX, chunkZ)
	c.SetBlockLight(x & 0xf, y, z & 0xf, light)
}
