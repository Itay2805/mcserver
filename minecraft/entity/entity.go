package entity

import (
	"github.com/google/uuid"
	"github.com/itay2805/mcserver/math"
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/proto/play"
)

type Type struct {
	Id			int
	Name		string
	Width		float64
	Height		float64
}


type IEntity interface {
	math.Spatial
	GetEntity() *Entity
	UpdateBounds()
}

const (
	isOnfire = 0x1
	isCrouching = 0x2
	isSprinting = 0x8
	isSwimming = 0x10
	isInvisible = 0x20
	hasGlowingEffect = 0x40
	isFlyingWithElytra = 0x80
)

type Entity struct {
	// the entity type info
	Type		*Type

	// Identifiers
	// per-runtime
	EID					int32
	// persistent
	UUID				uuid.UUID

	// the entity position
	Moved				bool
	Rotated				bool
	OnGroundChanged 	bool

	PrevPosition		math.Point
	Position			math.Point
	// TODO: PositionDelta
	Velocity 			math.Point
	Yaw, Pitch, HeadYaw	minecraft.Angle
	OnGround			bool

	// flags
	OnFire				bool
	Sprinting			bool
	Invisible			bool
	Glowing				bool
	Pose				minecraft.Pose

	// The metadata of the entity
	// has changed
	MetadataChanged		bool

	// the bounds of the entity
	bounds				*math.Rect

	// The entity equipment
	Equipment			[6]*play.Slot
	EquipmentChanged	int

	// The animation to play
	Animation			byte
}

func (e *Entity) GetFacing() minecraft.Face {
	// 225-256 or 0-31
	if 225 <= e.HeadYaw || e.HeadYaw <= 31 {
		return minecraft.FaceSouth
	} else if 161 <= e.HeadYaw && e.HeadYaw <= 224 {
		return minecraft.FaceEast
	} else if 97 <= e.HeadYaw && e.HeadYaw <= 160 {
		return minecraft.FaceNorth
	} else {
		return minecraft.FaceWest
	}
}

func (e *Entity) UpdateBounds() {
	e.bounds = math.NewRect(
		math.NewPoint(e.Position.X() - (e.Type.Width / 2), e.Position.Y(), e.Position.Z() - (e.Type.Width / 2)),
		[3]float64{
			e.Type.Width,
			e.Type.Height,
			e.Type.Width,
		},
	)
}

func (e *Entity) Bounds() *math.Rect {
	return e.bounds
}

func (e *Entity) GetEntity() *Entity {
	return e
}

func (e *Entity) WriteMetadata(writer *minecraft.EntityMetadataWriter) {
	val := byte(0)

	// flags
	if e.OnFire {
		val |= isOnfire
	}
	if e.Sprinting {
		val |= isSprinting
	}
	if e.Glowing {
		val |= hasGlowingEffect
	}

	writer.WriteByte(0, val)
	writer.WritePose(6, e.Pose)
}

func GetEntityTypeByName(name string) {

}
