package game

import (
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/block"
	"github.com/itay2805/mcserver/minecraft/item"
)

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

func transformFaceToPlayer(player *Player, blk *block.Block) uint16 {
	return blk.MinStateId + uint16(player.GetFacing().Invert()) - 2
}

func TransformItemToStateId(player *Player, itm *item.Item, face minecraft.Face) (uint16, bool) {
	switch itm {
	// Either on ground or based on block placed on
	case item.Torch: 						return transformTorch(face, block.Torch, block.WallTorch)
	case item.RedstoneTorch: 				return transformTorch(face, block.RedstoneTorch, block.RedstoneWallTorch)

	// TODO: can be on all 6 facings
	//	 - end rod
	//   - shulker_box

	// Facing based on player look
	case item.Loom:							return transformFaceToPlayer(player, block.Loom), true
	case item.Stonecutter:					return transformFaceToPlayer(player, block.Stonecutter), true
	case item.Anvil:						return transformFaceToPlayer(player, block.Anvil), true
	case item.ChippedAnvil:					return transformFaceToPlayer(player, block.ChippedAnvil), true
	case item.DamagedAnvil:					return transformFaceToPlayer(player, block.DamagedAnvil), true
	case item.WhiteGlazedTerracotta:		return transformFaceToPlayer(player, block.WhiteGlazedTerracotta), true
	case item.OrangeGlazedTerracotta:		return transformFaceToPlayer(player, block.OrangeGlazedTerracotta), true
	case item.MagentaGlazedTerracotta:		return transformFaceToPlayer(player, block.MagentaGlazedTerracotta), true
	case item.LightBlueGlazedTerracotta:	return transformFaceToPlayer(player, block.LightBlueGlazedTerracotta), true
	case item.YellowGlazedTerracotta:		return transformFaceToPlayer(player, block.YellowGlazedTerracotta), true
	case item.LimeGlazedTerracotta:			return transformFaceToPlayer(player, block.LimeGlazedTerracotta), true
	case item.PinkGlazedTerracotta:			return transformFaceToPlayer(player, block.PinkGlazedTerracotta), true
	case item.GrayGlazedTerracotta:			return transformFaceToPlayer(player, block.GrayGlazedTerracotta), true
	case item.LightGrayGlazedTerracotta:	return transformFaceToPlayer(player, block.LightGrayGlazedTerracotta), true
	case item.CyanGlazedTerracotta:			return transformFaceToPlayer(player, block.CyanGlazedTerracotta), true
	case item.PurpleGlazedTerracotta:		return transformFaceToPlayer(player, block.PurpleGlazedTerracotta), true
	case item.BlueGlazedTerracotta:			return transformFaceToPlayer(player, block.BlueGlazedTerracotta), true
	case item.BrownGlazedTerracotta:		return transformFaceToPlayer(player, block.BrownGlazedTerracotta), true
	case item.GreenGlazedTerracotta:		return transformFaceToPlayer(player, block.GreenGlazedTerracotta), true
	case item.RedGlazedTerracotta:			return transformFaceToPlayer(player, block.RedGlazedTerracotta), true
	case item.BlackGlazedTerracotta: 		return transformFaceToPlayer(player, block.BlackGlazedTerracotta), true

	// TODO: on ground based on player pos,
	//		 on wall based on the facing (like torch)
	// 	- heads
	//  - banners?

	default:
		// by default just turn the item into the block it represents
		// and use its default state
		return block.FromItem(itm).DefaultStateId, true
	}
}
