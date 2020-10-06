package block

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Base block type
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Material int

const (
	MaterialAir = Material(iota)
	MaterialRock
	MaterialDirt
	MaterialPlant
	MaterialWood
	MaterialWeb
	MaterialWool
)

// the base block type
type Block struct {
	// The identification of this block
	Id				int
	Name			string

	// The flatgrass state id for this block
	DefaultStateId	uint16

	// is this a solid block
	Solid			bool
	Material		Material

	// light related values
	Transparent		bool
	FilterLight		int
	EmitLight		int
}

func GetBlockById(stateId int) *Block {
	return blocks[stateId]
}

func GetBlockByStateId(stateId uint16) *Block {
	return stateIdToBlockId[stateId]
}
