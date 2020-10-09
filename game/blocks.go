package game

import (
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/block"
	"github.com/itay2805/mcserver/minecraft/item"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper functions for transforming special blocks
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func transformTorch(face minecraft.Face) (uint16, bool) {
	if face == minecraft.FaceBottom {
		// can't place on the ceiling
		return 0, false
	} else if face == minecraft.FaceTop {
		// just return the default torch
		return block.Torch.DefaultStateId, true
	} else {
		// turn the facing to an axis
		return block.WallTorch.MinStateId + uint16(face) - 2, true
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func transformFaceToPlayer(player *Player, blk *block.Block) uint16 {
	return blk.MinStateId + uint16(player.GetFacing().Invert()) - 2
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func transformStairs(player *Player, face minecraft.Face, blk *block.Block) uint16 {
	// NOTE: The shape is handled in the block update instead
	// 		 of here, since it will allow us to update the stairs nearby
	meta := block.StairsMeta{
		Facing:      player.GetFacing(),
		Half:        minecraft.FaceTop,
		Shape: 		 minecraft.ShapeStraight,
		Waterlogged: false,
	}

	// if placing on ceiling then turn upside down
	if face == minecraft.FaceTop {
		meta.Half = minecraft.FaceBottom
	}

	// TODO: if placing with cursor on top half then FaceBottom

	// TODO: if placing inside water block then turn waterlogged

	return blk.MinStateId + meta.ToMeta()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func transformFurnace(player *Player, blk *block.Block) uint16 {
	meta := block.FurnaceMeta{
		Facing:      player.GetFacing(),
		Lit: 		false,
	}
	return blk.MinStateId + meta.ToMeta()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// The top level transform function
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TransformItemToStateId(player *Player, itm *item.Item, face minecraft.Face) (uint16, bool) {
	switch itm {
	// Either on ground or based on block placed on
	case item.Torch: 						return transformTorch(face)

	// Facing based on player look
	case item.CarvedPumpkin:				return transformFaceToPlayer(player, block.CarvedPumpkin), true
	case item.JackOLantern:					return transformFaceToPlayer(player, block.JackOLantern), true
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
	case item.Loom:							return transformFaceToPlayer(player, block.Loom), true
	case item.Stonecutter:					return transformFaceToPlayer(player, block.Stonecutter), true

	// handle stairs
	case item.OakStairs:					return transformStairs(player, face, block.OakStairs), true
	case item.CobblestoneStairs:			return transformStairs(player, face, block.CobblestoneStairs), true
	case item.BrickStairs:					return transformStairs(player, face, block.BrickStairs), true
	case item.StoneBrickStairs:				return transformStairs(player, face, block.StoneBrickStairs), true
	case item.NetherBrickStairs:			return transformStairs(player, face, block.NetherBrickStairs), true
	case item.SandstoneStairs:				return transformStairs(player, face, block.SandstoneStairs), true
	case item.SpruceStairs:					return transformStairs(player, face, block.SpruceStairs), true
	case item.BirchStairs:					return transformStairs(player, face, block.BirchStairs), true
	case item.JungleStairs:					return transformStairs(player, face, block.JungleStairs), true
	case item.QuartzStairs:					return transformStairs(player, face, block.QuartzStairs), true
	case item.AcaciaStairs:					return transformStairs(player, face, block.AcaciaStairs), true
	case item.DarkOakStairs:				return transformStairs(player, face, block.DarkOakStairs), true
	case item.PrismarineStairs:				return transformStairs(player, face, block.PrismarineStairs), true
	case item.PrismarineBrickStairs:		return transformStairs(player, face, block.PrismarineBrickStairs), true
	case item.DarkPrismarineStairs:			return transformStairs(player, face, block.DarkPrismarineStairs), true
	case item.RedSandstoneStairs:			return transformStairs(player, face, block.RedSandstoneStairs), true
	case item.PurpurStairs:					return transformStairs(player, face, block.PurpurStairs), true
	case item.PolishedGraniteStairs:		return transformStairs(player, face, block.PolishedGraniteStairs), true
	case item.SmoothRedSandstoneStairs:		return transformStairs(player, face, block.SmoothRedSandstoneStairs), true
	case item.MossyStoneBrickStairs:		return transformStairs(player, face, block.MossyStoneBrickStairs), true
	case item.PolishedDioriteStairs:		return transformStairs(player, face, block.PolishedDioriteStairs), true
	case item.MossyCobblestoneStairs:		return transformStairs(player, face, block.MossyCobblestoneStairs), true
	case item.EndStoneBrickStairs:			return transformStairs(player, face, block.EndStoneBrickStairs), true
	case item.StoneStairs:					return transformStairs(player, face, block.StoneStairs), true
	case item.SmoothSandstoneStairs:		return transformStairs(player, face, block.SmoothSandstoneStairs), true
	case item.SmoothQuartzStairs:			return transformStairs(player, face, block.SmoothQuartzStairs), true
	case item.GraniteStairs:				return transformStairs(player, face, block.GraniteStairs), true
	case item.AndesiteStairs:				return transformStairs(player, face, block.AndesiteStairs), true
	case item.RedNetherBrickStairs:			return transformStairs(player, face, block.RedNetherBrickStairs), true
	case item.PolishedAndesiteStairs:		return transformStairs(player, face, block.PolishedAndesiteStairs), true
	case item.DioriteStairs:				return transformStairs(player, face, block.DioriteStairs), true

	// handle furnace
	case item.Furnace:						return transformFurnace(player, block.Furnace), true
	case item.BlastFurnace:					return transformFurnace(player, block.BlastFurnace), true

	default:
		// by default just turn the item into the block it represents
		// and use its default state
		return block.FromItem(itm).DefaultStateId, true
	}
}
