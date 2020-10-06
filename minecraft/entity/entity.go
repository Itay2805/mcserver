package entity

import (
	"github.com/google/uuid"
	"github.com/itay2805/mcserver/math"
	"github.com/itay2805/mcserver/minecraft"
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
	IsOnfire = 0x1
	IsCrouching = 0x2
	IsSprinting = 0x8
	IsSwimming = 0x10
	IsInvisible = 0x20
	HasGlowingEffect = 0x40
	Is = 0x80
)

type Pose int32

const (
	PoseStanding = Pose(iota)
	PoseFallFlying
	PoseSleeping
	PoseSwimming
	PoseSpinAttack
	PoseSneaking
	PoseDying
)

type Entity struct {
	// the entity type info
	Type		*Type

	// Identifiers
	// per-runtime
	EID					int
	// persistent
	UUID				uuid.UUID

	// the entity position
	Moved				bool
	Rotated				bool
	OnGroundChanged 	bool

	Position			math.Point
	// TODO: PositionDelta
	Velocity 			math.Point
	Yaw, Pitch, HeadYaw	minecraft.Angle
	OnGround			bool

	// flags
	OnFire				bool
	Sprinting			bool
	Glowing				bool

	// the current pose
	Pose 				Pose

	// The metadata of the entity
	// has changed
	MetadataChanged		bool

	// the bounds of the entity
	bounds				*math.Rect
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

	if e.OnFire {
		val |= 0x01
	}

	if e.Sprinting {
		val |= 0x08
	}

	if e.Glowing {
		val |= 0x40
	}

	if val != 0 {
		writer.WriteByte(0, val)
	}

	writer.WriteVarint(6, int32(e.Pose))
}

func GetEntityTypeByName(name string) {

}
