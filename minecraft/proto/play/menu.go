package play

import (
	"fmt"
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/item"
)

type Slot struct {
	ItemID    	int32
	ItemCount 	byte
	NBT       	interface{}
}

func (s *Slot) String() string {
	if s == nil {
		return fmt.Sprintf("Slot{}")
	} else {
		return fmt.Sprintf("Slot{ Item: %s, Count: %d, NBT: %s }", item.GetById(int(s.ItemID)).Name, s.ItemCount, s.NBT)
	}
}

func (s *Slot) CreateFake() *Slot {
	if s == nil {
		return nil
	} else {
		// TODO: make this handle nbt data properly
		return &Slot{
			ItemID:    s.ItemID,
			ItemCount: 69,
			NBT:       nil,
		}
	}
}

func (s *Slot) Encode(writer *minecraft.Writer) {
	if s == nil {
		writer.WriteBoolean(false)
	} else {
		writer.WriteBoolean(true)
		writer.WriteVarint(s.ItemID)
		writer.WriteByte(s.ItemCount)
		if s.NBT == nil {
			_ = minecraft.NbtMarshal(writer, struct{}{})
		} else {
			_ = minecraft.NbtMarshal(writer, s.NBT)
		}
	}
}

type SetSlot struct {
	WindowID 	int8
	Slot		int16
	SlotData	*Slot
}

func (s SetSlot) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x17)
	writer.WriteByte(byte(s.WindowID))
	writer.WriteShort(s.Slot)
	s.SlotData.Encode(writer)
}
