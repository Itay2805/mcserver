package block

import "github.com/itay2805/mcserver/minecraft/item"

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Base block type
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// the base block type
type Block struct {
	// The identification of this block
	Id				int
	Name			string

	// The flatgrass state id for this block
	DefaultStateId	uint16
	Item			*item.Item

	// is this a solid block
	Solid			bool

	// light related values
	Transparent		bool
	FilterLight		int
	EmitLight		int
}

func FromItem(item *item.Item) *Block {
	return blockByItemId[item.ID]
}

func GetById(stateId int) *Block {
	return blocks[stateId]
}

func GetByStateId(stateId uint16) *Block {
	return stateIdToBlockId[stateId]
}
