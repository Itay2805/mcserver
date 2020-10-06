package nullprovider

import "github.com/itay2805/mcserver/minecraft/chunk"

type NullProvider struct {
}

func (n *NullProvider) SaveChunk(chunk *chunk.Chunk) {

}

func (n *NullProvider) LoadChunk(x, z int) *chunk.Chunk {
	return nil
}

