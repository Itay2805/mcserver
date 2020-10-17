package game

import (
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/block"
	"github.com/itay2805/mcserver/minecraft/item"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper functions for transforming special blocks
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
// Transform a torch, just place it on a wall if needed and don't
// allow to place it in the air
//
func (p *Player) transformTorch(face minecraft.Face) (uint16, bool) {
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

//
// Have the block just face to player, this assumes there is no other state
// info other than a 4-way facing
//
func (p *Player) transformFaceToPlayer(blk *block.Block) uint16 {
	return blk.MinStateId + uint16(p.GetFacing().Invert()) - 2
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
// Transform stairs
//
// This handles the half the block is on (based on both the cursor and the face
// the player placed on) and the facing itself which is based on the player.
//
// The connections to other stairs are gonna be handled in the block update
// for simplicity since in placement we don't have the sync yet.
//
// TODO: handle water logged, based on the block of where the block is
// 		 gonna be placed
//
func (p *Player) transformStairs(face minecraft.Face, blk *block.Block) uint16 {
	// NOTE: The shape is handled in the block update instead
	// 		 of here, since it will allow us to update the stairs nearby
	meta := block.StairsMeta{
		Facing:      p.GetFacing(),
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

//
// Do furnace transformation, aka, rotate to where the player is facing
//
func (p *Player) transformFurnace(blk *block.Block) uint16 {
	meta := block.FurnaceMeta{
		Facing:      p.GetFacing(),
		Lit: 		false,
	}
	return blk.MinStateId + meta.ToMeta()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// The top level transform function
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
// Top level transformation function that turns an item to a block, this will also handle
// item uses like flint&steel and buckets.
//
// The reason for that is that it is just easier to handle it on block placement rather than
// item use as the game sends more info
//
func (p *Player) TransformItemToStateId(itm *item.Item, face minecraft.Face) (uint16, bool) {
	switch itm {
	// Either on ground or based on block placed on
	case item.Torch: 						return p.transformTorch(face)

	// Facing based on player look
	case item.CarvedPumpkin:				return p.transformFaceToPlayer(block.CarvedPumpkin), true
	case item.JackOLantern:					return p.transformFaceToPlayer(block.JackOLantern), true
	case item.Anvil:						return p.transformFaceToPlayer(block.Anvil), true
	case item.ChippedAnvil:					return p.transformFaceToPlayer(block.ChippedAnvil), true
	case item.DamagedAnvil:					return p.transformFaceToPlayer(block.DamagedAnvil), true
	case item.WhiteGlazedTerracotta:		return p.transformFaceToPlayer(block.WhiteGlazedTerracotta), true
	case item.OrangeGlazedTerracotta:		return p.transformFaceToPlayer(block.OrangeGlazedTerracotta), true
	case item.MagentaGlazedTerracotta:		return p.transformFaceToPlayer(block.MagentaGlazedTerracotta), true
	case item.LightBlueGlazedTerracotta:	return p.transformFaceToPlayer(block.LightBlueGlazedTerracotta), true
	case item.YellowGlazedTerracotta:		return p.transformFaceToPlayer(block.YellowGlazedTerracotta), true
	case item.LimeGlazedTerracotta:			return p.transformFaceToPlayer(block.LimeGlazedTerracotta), true
	case item.PinkGlazedTerracotta:			return p.transformFaceToPlayer(block.PinkGlazedTerracotta), true
	case item.GrayGlazedTerracotta:			return p.transformFaceToPlayer(block.GrayGlazedTerracotta), true
	case item.LightGrayGlazedTerracotta:	return p.transformFaceToPlayer(block.LightGrayGlazedTerracotta), true
	case item.CyanGlazedTerracotta:			return p.transformFaceToPlayer(block.CyanGlazedTerracotta), true
	case item.PurpleGlazedTerracotta:		return p.transformFaceToPlayer(block.PurpleGlazedTerracotta), true
	case item.BlueGlazedTerracotta:			return p.transformFaceToPlayer(block.BlueGlazedTerracotta), true
	case item.BrownGlazedTerracotta:		return p.transformFaceToPlayer(block.BrownGlazedTerracotta), true
	case item.GreenGlazedTerracotta:		return p.transformFaceToPlayer(block.GreenGlazedTerracotta), true
	case item.RedGlazedTerracotta:			return p.transformFaceToPlayer(block.RedGlazedTerracotta), true
	case item.BlackGlazedTerracotta: 		return p.transformFaceToPlayer(block.BlackGlazedTerracotta), true
	case item.Loom:							return p.transformFaceToPlayer(block.Loom), true
	case item.Stonecutter:					return p.transformFaceToPlayer(block.Stonecutter), true

	// handle stairs
	case item.OakStairs:					return p.transformStairs(face, block.OakStairs), true
	case item.CobblestoneStairs:			return p.transformStairs(face, block.CobblestoneStairs), true
	case item.BrickStairs:					return p.transformStairs(face, block.BrickStairs), true
	case item.StoneBrickStairs:				return p.transformStairs(face, block.StoneBrickStairs), true
	case item.NetherBrickStairs:			return p.transformStairs(face, block.NetherBrickStairs), true
	case item.SandstoneStairs:				return p.transformStairs(face, block.SandstoneStairs), true
	case item.SpruceStairs:					return p.transformStairs(face, block.SpruceStairs), true
	case item.BirchStairs:					return p.transformStairs(face, block.BirchStairs), true
	case item.JungleStairs:					return p.transformStairs(face, block.JungleStairs), true
	case item.QuartzStairs:					return p.transformStairs(face, block.QuartzStairs), true
	case item.AcaciaStairs:					return p.transformStairs(face, block.AcaciaStairs), true
	case item.DarkOakStairs:				return p.transformStairs(face, block.DarkOakStairs), true
	case item.PrismarineStairs:				return p.transformStairs(face, block.PrismarineStairs), true
	case item.PrismarineBrickStairs:		return p.transformStairs(face, block.PrismarineBrickStairs), true
	case item.DarkPrismarineStairs:			return p.transformStairs(face, block.DarkPrismarineStairs), true
	case item.RedSandstoneStairs:			return p.transformStairs(face, block.RedSandstoneStairs), true
	case item.PurpurStairs:					return p.transformStairs(face, block.PurpurStairs), true
	case item.PolishedGraniteStairs:		return p.transformStairs(face, block.PolishedGraniteStairs), true
	case item.SmoothRedSandstoneStairs:		return p.transformStairs(face, block.SmoothRedSandstoneStairs), true
	case item.MossyStoneBrickStairs:		return p.transformStairs(face, block.MossyStoneBrickStairs), true
	case item.PolishedDioriteStairs:		return p.transformStairs(face, block.PolishedDioriteStairs), true
	case item.MossyCobblestoneStairs:		return p.transformStairs(face, block.MossyCobblestoneStairs), true
	case item.EndStoneBrickStairs:			return p.transformStairs(face, block.EndStoneBrickStairs), true
	case item.StoneStairs:					return p.transformStairs(face, block.StoneStairs), true
	case item.SmoothSandstoneStairs:		return p.transformStairs(face, block.SmoothSandstoneStairs), true
	case item.SmoothQuartzStairs:			return p.transformStairs(face, block.SmoothQuartzStairs), true
	case item.GraniteStairs:				return p.transformStairs(face, block.GraniteStairs), true
	case item.AndesiteStairs:				return p.transformStairs(face, block.AndesiteStairs), true
	case item.RedNetherBrickStairs:			return p.transformStairs(face, block.RedNetherBrickStairs), true
	case item.PolishedAndesiteStairs:		return p.transformStairs(face, block.PolishedAndesiteStairs), true
	case item.DioriteStairs:				return p.transformStairs(face, block.DioriteStairs), true

	// handle furnace
	case item.Furnace:						return p.transformFurnace(block.Furnace), true
	case item.BlastFurnace:					return p.transformFurnace(block.BlastFurnace), true

	default:
		// by default just turn the item into the block it represents
		// and use its default state
		blk := block.FromItem(itm)
		if blk == nil {
			return 0, false
		} else {
			return blk.DefaultStateId, true
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Block updates
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
// This is done right after the chunk processes all of the block changes that happened in this tick
// the reason for this is to make sure the map is in the most updated state.
//
// if you want to register future block updates just use the ticker to register a block update, and
// since right after the ticker the block updates happen it will just work
//
func BlockUpdate(position minecraft.Position, blk block.Block) {
}
