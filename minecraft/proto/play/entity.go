package play

import (
	"github.com/google/uuid"
	"github.com/itay2805/mcserver/minecraft"
)

// TODO: Spawn Entity

// TODO: Spawn Experience Orn

// TODO: Spawn Weather Entity

// TODO: Spawn Living Entity

// TODO: Spawn Painting

type SpawnPlayer struct {
	EntityID int32
	UUID     uuid.UUID
	X        float64
	Y        float64
	Z        float64
	Yaw      minecraft.Angle
	Pitch    minecraft.Angle
}

func (r SpawnPlayer) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x05)
	writer.WriteVarint(r.EntityID)
	writer.WriteUUID(r.UUID)
	writer.WriteDouble(r.X)
	writer.WriteDouble(r.Y)
	writer.WriteDouble(r.Z)
	writer.WriteAngle(r.Yaw)
	writer.WriteAngle(r.Pitch)
}

const (
	AnimationNone = -1
	AnimationSwingMainHand = 0
	AnimationTakeDamage = 1
	AnimationLeaveBed = 2
	AnimationSwingOffhand = 3
	AnimationCriticalEffect = 4
	AnimationMagicCriticalEffect = 5
)

type EntityAnimation struct {
	EntityId		int32
	Animation		byte
}

func (r EntityAnimation) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x06)
	writer.WriteVarint(r.EntityId)
	writer.WriteByte(r.Animation)
}

type EntityPosition struct {
	EntityId		int32
	DeltaX			int16
	DeltaY			int16
	DeltaZ			int16
	OnGround		bool
}

func (r EntityPosition) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x29)
	writer.WriteVarint(r.EntityId)
	writer.WriteShort(r.DeltaX)
	writer.WriteShort(r.DeltaY)
	writer.WriteShort(r.DeltaZ)
	writer.WriteBoolean(r.OnGround)
}

type EntityPositionAndRotation struct {
	EntityId		int32
	DeltaX			int16
	DeltaY			int16
	DeltaZ			int16
	Yaw				minecraft.Angle
	Pitch			minecraft.Angle
	OnGround		bool
}

func (r EntityPositionAndRotation) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x2A)
	writer.WriteVarint(r.EntityId)
	writer.WriteShort(r.DeltaX)
	writer.WriteShort(r.DeltaY)
	writer.WriteShort(r.DeltaZ)
	writer.WriteAngle(r.Yaw)
	writer.WriteAngle(r.Pitch)
	writer.WriteBoolean(r.OnGround)
}

type EntityRotation struct {
	EntityID int32
	Yaw      minecraft.Angle
	Pitch    minecraft.Angle
	OnGround bool
}

func (r EntityRotation) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x2B)
	writer.WriteVarint(r.EntityID)
	writer.WriteAngle(r.Yaw)
	writer.WriteAngle(r.Pitch)
	writer.WriteBoolean(r.OnGround)
}

type EntityMovement struct {
	EntityID int32
}

func (r EntityMovement) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x2C)
	writer.WriteVarint(r.EntityID)
}

type DestroyEntities struct {
	EntityIDs		[]int32
}

func (r DestroyEntities) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x38)
	writer.WriteVarint(int32(len(r.EntityIDs)))
	for _, eid := range r.EntityIDs {
		writer.WriteVarint(eid)
	}
}

// TODO: Remove Entity Effect

type EntityHeadLook struct {
	EntityID		int32
	HeadYaw			minecraft.Angle
}

func (r EntityHeadLook) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x3C)
	writer.WriteVarint(r.EntityID)
	writer.WriteAngle(r.HeadYaw)
}

type IMetadata interface {
	WriteMetadata(writer *minecraft.EntityMetadataWriter)
}

type EntityMetadata struct {
	EntityID 		int32
	Metadata		IMetadata
}

func (r EntityMetadata) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x44)
	writer.WriteVarint(r.EntityID)
	ew := writer.StartEntityMetadata()
	r.Metadata.WriteMetadata(ew)
	ew.Done()
}

// TODO: Attach Entity

// TODO: Entity Velocity

// TODO: Entity Equipment

// TODO: Entity Sound Effect

type EntityTeleport struct {
	EntityID		int32
	X				float64
	Y				float64
	Z				float64
	Yaw				minecraft.Angle
	Pitch			minecraft.Angle
	OnGround		bool
}

func (r EntityTeleport) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x57)
	writer.WriteVarint(r.EntityID)
	writer.WriteDouble(r.X)
	writer.WriteDouble(r.Y)
	writer.WriteDouble(r.Z)
	writer.WriteAngle(r.Yaw)
	writer.WriteAngle(r.Pitch)
	writer.WriteBoolean(r.OnGround)
}

// TODO: Entity Properties

// TODO: Entity Effect