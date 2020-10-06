package play

import (
	"github.com/google/uuid"
	"github.com/itay2805/mcserver/minecraft"
)

const (
	ChatBox 		= 0
	SystemMessage 	= 1
	GameInfo		= 2
)

type ChatMessage struct {
	Data		minecraft.Chat
	Position	byte
}

func (r ChatMessage) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x0F)
	writer.WriteChat(r.Data)
	writer.WriteByte(r.Position)
}

type Disconnect struct {
	Reason 	minecraft.Chat
}

func (r Disconnect) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x1B)
	writer.WriteChat(r.Reason)
}

type KeepAlive struct {
	KeepAliveId int64
}

func (r KeepAlive) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x21)
	writer.WriteLong(r.KeepAliveId)
}

type JoinGame struct {
	EntityId			int32
	Gamemode			byte
	Dimension 			int32
	HashedSeed			int64
	LevelType			string
	ViewDistance		int32
	ReducedDebugInfo	bool
	EnableRespawnScreen	bool
}

func (r JoinGame) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x26)
	writer.WriteInt(r.EntityId)
	writer.WriteByte(r.Gamemode)
	writer.WriteInt(r.Dimension)
	writer.WriteLong(r.HashedSeed)
	// Max Players, Was once used by the client to draw the player list, but now is ignored
	// source: https://wiki.vg/Protocol#Join_Game
	writer.WriteByte(0)
	writer.WriteString(r.LevelType)
	writer.WriteVarint(r.ViewDistance)
	writer.WriteBoolean(r.ReducedDebugInfo)
	writer.WriteBoolean(r.EnableRespawnScreen)
}

type PlayerInfo struct {
	AddPlayer			[]PIAddPlayer
	UpdateGamemode		[]PIUpdateGamemode
	UpdateLatency		[]PIUpdateLatency
	UpdateDisplayName	[]PIDisplayName
	RemovePlayer		[]PIRemovePlayer
}

func (p PlayerInfo) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x34)

	if len(p.AddPlayer) != 0 {
		writer.WriteVarint(0)
		writer.WriteVarint(int32(len(p.AddPlayer)))
		for _, u := range p.AddPlayer {
			u.Encode(writer)
		}

	} else if len(p.UpdateGamemode) != 0 {
		writer.WriteVarint(1)
		writer.WriteVarint(int32(len(p.UpdateGamemode)))
		for _, u := range p.UpdateGamemode {
			u.Encode(writer)
		}

	} else if len(p.UpdateLatency) != 0 {
		writer.WriteVarint(2)
		writer.WriteVarint(int32(len(p.UpdateLatency)))
		for _, u := range p.UpdateLatency {
			u.Encode(writer)
		}

	} else if len(p.UpdateDisplayName) != 0 {
		writer.WriteVarint(3)
		writer.WriteVarint(int32(len(p.UpdateDisplayName)))
		for _, u := range p.UpdateDisplayName {
			u.Encode(writer)
		}

	} else if len(p.RemovePlayer) != 0 {
		writer.WriteVarint(4)
		writer.WriteVarint(int32(len(p.RemovePlayer)))
		for _, u := range p.RemovePlayer {
			u.Encode(writer)
		}

	}
}

type PIAddPlayer struct {
	UUID			uuid.UUID
	Name 				string
	//Props				[]mojang.PlayerProperty
	Gamemode			int32
	Ping				int32
	DisplayName			*minecraft.Chat
}

func (p PIAddPlayer) Encode(writer *minecraft.Writer) {
	writer.WriteUUID(p.UUID)
	writer.WriteString(p.Name)

	writer.WriteVarint(0)
	//writer.WriteVarint(int32(len(p.Props)))
	//for _, prop := range p.Props {
	//	writer.WriteString(prop.Name)
	//	writer.WriteString(prop.Value)
	//	if len(prop.Signature) != 0 {
	//		writer.WriteBoolean(true)
	//		writer.WriteString(prop.Signature)
	//	} else {
	//		writer.WriteBoolean(false)
	//	}
	//}

	writer.WriteVarint(p.Gamemode)
	writer.WriteVarint(p.Ping)

	if p.DisplayName != nil {
		writer.WriteBoolean(true)
		writer.WriteChat(*p.DisplayName)
	} else {
		writer.WriteBoolean(false)
	}
}

type PIUpdateGamemode struct {
	UUID			uuid.UUID
	Gamemode		int32
}

func (p PIUpdateGamemode) Encode(writer *minecraft.Writer) {
	writer.WriteUUID(p.UUID)
	writer.WriteVarint(p.Gamemode)
}

type PIUpdateLatency struct {
	UUID			uuid.UUID
	Ping		int32
}

func (p PIUpdateLatency) Encode(writer *minecraft.Writer) {
	writer.WriteUUID(p.UUID)
	writer.WriteVarint(p.Ping)
}

type PIDisplayName struct {
	UUID			uuid.UUID
	DisplayName		*minecraft.Chat
}

func (p PIDisplayName) Encode(writer *minecraft.Writer) {
	writer.WriteUUID(p.UUID)
	if p.DisplayName != nil {
		writer.WriteBoolean(true)
		writer.WriteChat(*p.DisplayName)
	} else {
		writer.WriteBoolean(false)
	}
}

type PIRemovePlayer struct {
	UUID			uuid.UUID
}

func (p PIRemovePlayer) Encode(writer *minecraft.Writer) {
	writer.WriteUUID(p.UUID)
}

const (
	PlayerPositionXRelative = 0x01
	PlayerPositionYRelative = 0x02
	PlayerPositionZRelative = 0x04
	PlayerPositionYawRelative = 0x08
	PlayerPositionPitchRelative = 0x10
)

type PlayerPositionAndLook struct {
	X, Y, Z 	float64
	Yaw, Pitch 	float32
	Flags 		byte
	TeleportId	int32
}

func (r PlayerPositionAndLook) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x36)
	writer.WriteDouble(r.X)
	writer.WriteDouble(r.Y)
	writer.WriteDouble(r.Z)
	writer.WriteFloat(r.Yaw)
	writer.WriteFloat(r.Pitch)
	writer.WriteByte(r.Flags)
	writer.WriteVarint(r.TeleportId)
}
