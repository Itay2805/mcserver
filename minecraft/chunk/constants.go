package chunk

const (
	// height in blocks of a chunk column
	ChunkHeight = 256

	// width in blocks of a chunk column
	ChunkWidth = 16

	// height in blocks of a chunk section
	SectionHeight = 16

	// width in blocks of a chunk section
	SectionWidth = 16

	// volume in blocks of a chunk section
	SectionVolume = SectionHeight * SectionWidth * SectionWidth

	// number of chunk sections in a chunk column
	NumSections = 16

	// number of light sections in a chunk column
	NumLightSections = NumSections + 2

	// number of elements in the light byte array
	LightVolume	= (SectionHeight * SectionWidth * SectionWidth) / 2

	// maximum number of bits per block allowed when using the section palette.
	// values above will switch to global palette
	MaxBitsPerBlock = 8

	// number of bits used for each block in the global palette.
	// this value should not be hardcoded according to wiki.vg
	GlobalBitsPerBlock = 14
)
