package play

import "github.com/itay2805/mcserver/minecraft"

const (
	// TODO: add from https://wiki.vg/Protocol#Effect
	EffectBlockBreak = int32(2001)
)

type Effect struct {
	EffectID				int32
	Location				minecraft.Position
	Data					int32
	DisableRelativeVolume	bool
}

func (r Effect) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x23)
	writer.WriteInt(r.EffectID)
	writer.WritePosition(r.Location)
	writer.WriteInt(r.Data)
	writer.WriteBoolean(r.DisableRelativeVolume)
}
