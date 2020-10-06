package status

import (
	"github.com/itay2805/mcserver/game"
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/minecraft/proto/status"
	"math"
)

func HandleRequest(player *game.Player, reader *minecraft.Reader) {
	player.Send(status.Response{
		Response: status.ServerListResponse{
			Version:     status.ServerListVersion{
				Name:     "1.15.2",
				Protocol: 578,
			},
			Players:     status.ServerListPlayers{
				Max:    math.MaxInt32,
				Online: int(game.GetPlayerCount()),
			},
			Description: minecraft.Chat{
				Text:          	"Go Minecraft Server!",
				Italic:			true,
				Color: 			"red",
			},
		},
	})
}

func HandlePing(player *game.Player, reader *minecraft.Reader) {
	payload := reader.ReadLong()
	player.Send(status.Pong{
		Payload: payload,
	})
}
