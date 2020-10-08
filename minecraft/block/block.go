package block

import (
	"fmt"
	"github.com/itay2805/mcserver/minecraft/item"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Base block type
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// the base block type
type Block struct {
	// The identification of this block
	Id				int
	Name			string

	// The item corresponding to this block
	// if any
	Item			*item.Item

	// state id
	MinStateId		uint16
	MaxStateId		uint16
	DefaultStateId	uint16

	// is this a solid block
	Solid			bool

	// light related values
	Transparent		bool
	FilterLight		int
	EmitLight		int
}

func (b *Block) String() string {
	return fmt.Sprintf("Block{ Name: \"%s\" }", b.Name)
}

func FromItem(item *item.Item) *Block {
	if item.ID > len(blockByItemId) {
		return nil
	}
	return blockByItemId[item.ID]
}

func GetById(stateId int) *Block {
	return blocks[stateId]
}

func GetByStateId(stateId uint16) *Block {
	return stateIdToBlockId[stateId]
}
