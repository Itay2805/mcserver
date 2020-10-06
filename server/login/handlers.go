package login

import (
	"crypto/md5"
	"github.com/google/uuid"
	"github.com/itay2805/mcserver/game"
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/entity"
	"github.com/itay2805/mcserver/minecraft/proto/login"
	"github.com/itay2805/mcserver/server/socket"
	"log"
)

func offlineUUID(name string) uuid.UUID {
	var version = 3
	h := md5.New()
	h.Reset()
	h.Write([]byte("OfflinePlayer:" + name))
	s := h.Sum(nil)
	var uuid uuid.UUID
	copy(uuid[:], s)
	uuid[6] = (uuid[6] & 0x0f) | uint8((version&0xf)<<4)
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // RFC 4122 variant
	return uuid
}

func HandleLoginStart(player *game.Player, reader *minecraft.Reader) {
	name := reader.ReadString(16)
	log.Println("Got login request from", player.RemoteAddr(), "connecting from", player.RemoteAddr())

	uuid := offlineUUID(name)

	// Set Compression (block until done)
	player.SendSync(login.SetCompression{Threshold: 128})

	// sleep a bit before enabling compression so we will
	// send the packet without compression
	player.EnableCompression()

	// create the player object
	player.Player = entity.NewPlayer(name, uuid)

	// now we can actually send the login success
	player.Send(login.LoginSuccess{
		Uuid:     uuid,
		Username: name,
	})

	// we are not in play state
	player.State = socket.Play
	game.JoinPlayer(player)
}
