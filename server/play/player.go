package play

import (
	"github.com/itay2805/mcserver/game"
	"github.com/itay2805/mcserver/minecraft"
)

func HandleClientSettings(player *game.Player, reader *minecraft.Reader) {
	_ = reader.ReadString(16) // locale
	player.Change(&player.ViewDistance, int(reader.ReadByte()), nil)
	_ = reader.ReadVarint() // Chat Mode
	_ = reader.ReadBoolean() // Chat colors
	player.Change(&player.SkinMask, reader.ReadByte(), nil)
	player.Change(&player.MainHand, byte(reader.ReadVarint()), nil)
	player.Change(&player.MetadataChanged, true, nil)
}

func HandlePlayerPosition(player *game.Player, reader *minecraft.Reader) {
	player.Change(&player.Position[0], reader.ReadDouble(), &player.Moved)
	player.Change(&player.Position[1], reader.ReadDouble(), &player.Moved)
	player.Change(&player.Position[2], reader.ReadDouble(), &player.Moved)
	player.Change(&player.OnGround, reader.ReadBoolean(), &player.OnGroundChanged)
}

func HandlePlayerPositionAndRotation(player *game.Player, reader *minecraft.Reader) {
	player.Change(&player.Position[0], reader.ReadDouble(), &player.Moved)
	player.Change(&player.Position[1], reader.ReadDouble(), &player.Moved)
	player.Change(&player.Position[2], reader.ReadDouble(), &player.Moved)
	yaw := minecraft.ToAngle(reader.ReadFloat())
	player.Change(&player.Yaw, yaw, &player.Rotated)
	player.Change(&player.HeadYaw, yaw, nil)
	player.Change(&player.Pitch, minecraft.ToAngle(reader.ReadFloat()), &player.Rotated)
	player.Change(&player.OnGround, reader.ReadBoolean(), &player.OnGroundChanged)
}

func HandlePlayerRotation(player *game.Player, reader *minecraft.Reader) {
	yaw := minecraft.ToAngle(reader.ReadFloat())
	player.Change(&player.Yaw, yaw, &player.Rotated)
	player.Change(&player.HeadYaw, yaw, nil)
	player.Change(&player.Pitch, minecraft.ToAngle(reader.ReadFloat()), &player.Rotated)
	player.Change(&player.OnGround, reader.ReadBoolean(), &player.OnGroundChanged)
}

func HandlePlayerMovement(player *game.Player, reader *minecraft.Reader) {
	player.Change(&player.OnGround, reader.ReadBoolean(), &player.OnGroundChanged)
}
