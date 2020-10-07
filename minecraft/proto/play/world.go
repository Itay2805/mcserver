package play

import (
	"github.com/itay2805/mcserver/minecraft"
)

type BlockBreakAnimation struct {
	EntityId		int32
	Location		minecraft.Position
	DestroyStage	byte
}

func (r BlockBreakAnimation) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x09)
	writer.WriteVarint(r.EntityId)
	writer.WritePosition(r.Location)
	writer.WriteByte(r.DestroyStage)
}

// TODO: Block Entity Data

// TODO: Block Action

type BlockChange struct {
	Location 		minecraft.Position
	BlockID			uint16
}

func (r BlockChange) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x0C)
	writer.WritePosition(r.Location)
	writer.WriteVarint(int32(r.BlockID))
}

type BlockRecord struct {
	BlockX 		byte
	BlockZ 		byte
	BlockY		byte
	BlockState 	uint16
}

type MultiBlockChange struct {
	ChunkX  int32
	ChunkZ  int32
	Records []BlockRecord
}

func (r MultiBlockChange) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x10)
	writer.WriteInt(r.ChunkX)
	writer.WriteInt(r.ChunkZ)
	writer.WriteVarint(int32(len(r.Records)))
	for _, rec := range r.Records {
		writer.WriteByte(rec.BlockX << 4 | rec.BlockZ)
		writer.WriteByte(rec.BlockY)
		writer.WriteVarint(int32(rec.BlockState))
	}
}

type UnloadChunk struct {
	ChunkX 	int32
	ChunkZ 	int32
}

func (p UnloadChunk) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x1E)
	writer.WriteInt(p.ChunkX)
	writer.WriteInt(p.ChunkZ)
}

type UpdateViewPosition struct {
	ChunkX		int32
	ChunkZ		int32
}

func (p UpdateViewPosition) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x41)
	writer.WriteVarint(p.ChunkX)
	writer.WriteVarint(p.ChunkZ)
}

type TimeUpdate struct {
	WorldAge	int64
	TimeOfDay	int64
}

func (p TimeUpdate) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x4F)
	writer.WriteLong(p.WorldAge)
	writer.WriteLong(p.TimeOfDay)
}