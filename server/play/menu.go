package play

import (
	"errors"
	"github.com/itay2805/mcserver/game"
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/proto/play"
	"log"
)

func HandleHeldItemChange(player *game.Player, reader *minecraft.Reader) {
	slot := reader.ReadShort()

	player.Change(func() {
		player.ChangeHeldItem(int(slot))
	})
}


func HandleCreativeInventoryAction(player *game.Player, reader *minecraft.Reader) {
	// TODO: Check for creative mode

	slot := reader.ReadShort()

	if reader.ReadBoolean() {
		itemId := reader.ReadVarint()
		itemCount := reader.ReadByte()
		var nbt interface{}

		err := minecraft.NewNbtDecoder(reader).Decode(&nbt)
		if errors.Is(err, minecraft.ErrEND) {
			nbt = nil
		} else if err != nil {
			log.Panicln(err)
		}

		newSlot := &play.Slot{
			ItemID:    itemId,
			ItemCount: itemCount,
			NBT:       nbt,
		}

		player.Change(func() {
			player.UpdateInventory(int(slot), newSlot)
		})
	} else {
		player.Change(func() {
			player.UpdateInventory(int(slot), nil)
		})
	}
}