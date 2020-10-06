package play

import "github.com/itay2805/mcserver/minecraft"

type Slot struct {
	ItemID    int16
	ItemCount byte
	ItemMeta  int16
	NBT       interface{}
}

func (s Slot) CreateFake() Slot {
	// TODO: make this handle enchantment effect
	return Slot{
		ItemID:    s.ItemID,
		ItemCount: 69,
		ItemMeta:  s.ItemMeta,
		NBT:       nil,
	}
}

func (s Slot) Encode(writer *minecraft.Writer) {
	writer.WriteShort(s.ItemID)
	if s.ItemID != -1 {
		writer.WriteByte(s.ItemCount)
		writer.WriteShort(s.ItemMeta)
		if s.NBT == nil {
			_ = minecraft.NbtMarshal(writer, struct{}{})
		} else {
			_ = minecraft.NbtMarshal(writer, s.NBT)
		}
	}
}

type SetSlot struct {
	WindowID 	byte
	Slot		int16
	SlotData	Slot
}


func (s SetSlot) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x16)
	writer.WriteByte(s.WindowID)
	writer.WriteShort(s.Slot)
	s.SlotData.Encode(writer)
}
