package biome

type Type int

const (
	TypeOcean = Type(iota)
	TypePlains
	TypeDesert
	TypeHills
	TypeForest
	TypeTaiga
	TypeSwamp
	TypeRiver
	TypeNether
	TypeSky
	TypeSnow
	TypeMushroomIsland
	TypeBeach
	TypeJungle
	TypeStoneBeach
	TypeSavanna
	TypeMesa
	TypeVoid
)

type Biome struct {
	Id			int
	Name		string
	Type		Type
	Rainfall	float32
	Temperature	float32
}

func AreSimilar(a, b *Biome) bool {
	if a == b {
		return true
	}

	if a == WoodedBadlandsPlateau || a == BadlandsPlateau {
		return b == WoodedBadlandsPlateau || b == BadlandsPlateau
	}

	return a.Type == b.Type
}

var shallowOceanBits = 	(1 << Ocean.Id) |
						(1 << FrozenOcean.Id) |
						(1 << WarmOcean.Id) |
						(1 << LukewarmOcean.Id) |
						(1 << ColdOcean.Id)

func IsShallowOcean(b *Biome) bool {
	return ((1 << b.Id) & shallowOceanBits) != 0
}

var deepOceanBits = (1 << DeepOcean.Id) |
				(1 << DeepWarmOcean.Id) |
				(1 << DeepLukewarmOcean.Id) |
				(1 << DeepColdOcean.Id) |
				(1 << DeepFrozenOcean.Id)

func IsDeepOcean(b *Biome) bool {
	return ((1 << b.Id) & deepOceanBits) != 0
}


var oceanBits = (1 << Ocean.Id) |
				(1 << FrozenOcean.Id) |
				(1 << WarmOcean.Id) |
				(1 << LukewarmOcean.Id) |
				(1 << ColdOcean.Id) |
				(1 << DeepOcean.Id) |
				(1 << DeepWarmOcean.Id) |
				(1 << DeepLukewarmOcean.Id) |
				(1 << DeepColdOcean.Id) |
				(1 << DeepFrozenOcean.Id)

func IsOceanic(b *Biome) bool {
	return ((1 << b.Id) & oceanBits) != 0
}
