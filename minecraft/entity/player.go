package entity

import (
	"github.com/google/uuid"
	"github.com/itay2805/mcserver/math"
	"github.com/itay2805/mcserver/minecraft"
)

const (
	HandLeft = 0
	HandRight = 1
)

type Player struct {
	Living

	// visual
	SkinMask	byte
	MainHand	byte

	Username	string
	Health		float32
	Food		float32
}

// The player entity type
var playerType = byName["player"]

func NewPlayer(username string, uuid uuid.UUID) *Player {
	return &Player{
		Living:   Living{
			Entity:        Entity{
				Type:      playerType,
				EID:       0, // TODO: generate
				UUID:      uuid,
				Position:  math.Point{ 0, 0, 0 },
				Velocity:  math.Point{ 0, 0, 0},
				Yaw:       0,
				Pitch:     0,
				HeadYaw:   0,
				OnGround:  true,
				OnFire:    false,
				Sprinting: false,
				Glowing:   false,
				Pose:      PoseStanding,
			},
			IsHandActive:  false,
			OffhandActive: false,
		},
		SkinMask: 0,
		MainHand: 0,
		Username: username,
		Health:   20,
		Food:     20,
	}
}

func (p *Player) WriteMetadata(writer *minecraft.EntityMetadataWriter) {
	p.Living.WriteMetadata(writer)
	writer.WriteByte(16, p.SkinMask)
	writer.WriteByte(17, p.MainHand)
}
