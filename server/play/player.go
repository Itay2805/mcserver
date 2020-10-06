package play

import (
	"github.com/itay2805/mcserver/game"
	"github.com/itay2805/mcserver/minecraft"
	"log"
)

func HandleClientSettings(player *game.Player, reader *minecraft.Reader) {
	_ = reader.ReadString(16) // locale
	viewDistance := reader.ReadByte()
	_ = reader.ReadVarint() // Chat Mode
	_ = reader.ReadBoolean() // Chat colors
	skinMask := reader.ReadByte()
	mainHand := reader.ReadVarint()

	player.Change(func() {
		player.ViewDistance = int(viewDistance)
		player.SkinMask = skinMask
		player.MainHand = byte(mainHand)
		player.MetadataChanged = true
	})
}

func HandlePlayerPosition(player *game.Player, reader *minecraft.Reader) {
	x := reader.ReadDouble()
	y := reader.ReadDouble()
	z := reader.ReadDouble()
	onGround := reader.ReadBoolean()
	player.Change(func() {
		prevX := player.Position[0]
		prevY := player.Position[1]
		prevZ := player.Position[2]

		player.Position[0] = x
		player.Position[1] = y
		player.Position[2] = z
		player.OnGround = onGround

		player.Moved = prevX != x || prevY != y || prevZ != z
	})
}

func HandlePlayerPositionAndRotation(player *game.Player, reader *minecraft.Reader) {
	x := reader.ReadDouble()
	y := reader.ReadDouble()
	z := reader.ReadDouble()
	yaw := minecraft.ToAngle(reader.ReadFloat())
	pitch := minecraft.ToAngle(reader.ReadFloat())
	onGround := reader.ReadBoolean()
	player.Change(func() {
		prevX := player.Position[0]
		prevY := player.Position[1]
		prevZ := player.Position[2]
		prevYaw := player.Yaw
		prevPitch := player.Pitch

		player.Position[0] = x
		player.Position[1] = y
		player.Position[2] = z
		player.Yaw = yaw
		player.HeadYaw = yaw
		player.Pitch = pitch
		player.OnGround = onGround

		player.Moved = prevX != x || prevY != y || prevZ != z
		player.Rotated = pitch != prevPitch || yaw != prevYaw
	})

}

func HandlePlayerRotation(player *game.Player, reader *minecraft.Reader) {
	yaw := minecraft.ToAngle(reader.ReadFloat())
	pitch := minecraft.ToAngle(reader.ReadFloat())
	onGround := reader.ReadBoolean()
	player.Change(func() {
		prevYaw := player.Yaw
		prevPitch := player.Pitch

		player.Yaw = yaw
		player.Pitch = pitch
		player.HeadYaw = yaw
		player.OnGround = onGround

		player.Rotated = pitch != prevPitch || yaw != prevYaw
	})
}

func HandlePlayerMovement(player *game.Player, reader *minecraft.Reader) {
	onGround := reader.ReadBoolean()
	player.Change(func() {
		player.OnGround = onGround
	})
}

func HandlePlayerAbilities(player *game.Player, reader *minecraft.Reader) {
	flags := reader.ReadByte()

	player.Change(func() {
		player.Flying = (flags & 0x2) != 0
	})
}

func HandleEntityAction(player *game.Player, reader *minecraft.Reader) {
	eid := reader.ReadVarint()
	aid := reader.ReadVarint()

	if eid != player.EID {
		log.Println("HandleEntityAction:", player, "sent invalid entity id", eid, "(expected", player.EID, ")")
		return
	}

	switch aid {
	case 0:
		player.Change(func() {
			player.Pose = minecraft.PoseSneaking
			player.MetadataChanged = true
		})
	case 1:
		player.Change(func() {
			if player.Pose == minecraft.PoseSneaking {
				player.Pose = minecraft.PoseStanding
				player.MetadataChanged = true
			}
		})
	case 3:
		player.Change(func() {
			player.Sprinting = true
			player.MetadataChanged = true
		})
	case 4:
		player.Change(func() {
			player.Sprinting = false
			player.MetadataChanged = true
		})
	default:
		log.Println("HandleEntityAction:", player, "sent invalid action id", aid)
	}
}
