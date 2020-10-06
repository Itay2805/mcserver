package entity

import "github.com/itay2805/mcserver/minecraft"

type Living struct {
	Entity

	// active hand
	IsHandActive	bool
	OffhandActive	bool
}

func (p *Living) GetEntity() *Entity {
	return &p.Entity
}

func (p *Living) WriteMetadata(writer *minecraft.EntityMetadataWriter) {
	p.Entity.WriteMetadata(writer)

	val := byte(0)

	if p.IsHandActive {
		val |= 0x1
	}

	if p.OffhandActive {
		val |= 0x2
	}

	writer.WriteByte(7, val)
}
