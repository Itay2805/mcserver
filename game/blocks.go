package game

import (
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/block"
	"github.com/itay2805/mcserver/minecraft/item"
)

// for torches we need to simply choose either the standing type or the wall type
// this code is common for both a normal and redstone torch
func transformTorch(face minecraft.Face, normal, wall *block.Block) (uint16, bool) {
	if face == minecraft.FaceBottom {
		// can't place on the ceiling
		return 0, false
	} else if face == minecraft.FaceTop {
		// just return the default torch
		return normal.DefaultStateId, true
	} else {
		// turn the facing to an axis
		return wall.MinStateId + uint16(face) - 2, true
	}
}

func TransformItemToStateId(itm *item.Item, face minecraft.Face) (uint16, bool) {
	switch itm {
	case item.Torch: 			return transformTorch(face, block.Torch, block.WallTorch)
	case item.RedstoneTorch: 	return transformTorch(face, block.RedstoneTorch, block.RedstoneWallTorch)
	default:
		// by default just turn the item into the block it represents
		// and use its default state
		return block.FromItem(itm).DefaultStateId, true
	}
}
