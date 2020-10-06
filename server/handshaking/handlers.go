package handshaking

import (
	"github.com/itay2805/mcserver/game"
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/proto/login"
	"github.com/itay2805/mcserver/server/socket"
	"log"
)

func HandleHandshaking(player *game.Player, reader *minecraft.Reader) {
	protocolVersion := reader.ReadVarint()
	_ = reader.ReadString(255)
	_ = reader.ReadUShort()
	nextState := reader.ReadVarint()

	if nextState == 1 {
		// go to status, no need for protocol version check
		player.State = socket.Status

	} else if nextState == 2 {
		// go to login, only if the protocol version is good
		if protocolVersion != 578 {
			log.Println("Client at", player.RemoteAddr(), "tried to login from invalid client", protocolVersion, "(expected 340)")

			// kick the player
			player.Send(login.Disconnect{
				Reason: minecraft.Chat{
					Text:          "Invalid protocol version",
					Bold:          true,
					Italic:        true,
					Underlined:    true,
					Strikethrough: true,
					Color:         "red",
				},
			})

			// set the player as disconnected
			player.State = socket.Disconnected
		} else {
			// move the player to login state
			player.State = socket.Login
		}

	} else {
		log.Panicln("Requested invalid state", nextState)
	}
}
