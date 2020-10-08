package server

import (
	"errors"
	"github.com/itay2805/mcserver/game"
	"github.com/itay2805/mcserver/minecraft"
	"github.com/itay2805/mcserver/server/handshaking"
	"github.com/itay2805/mcserver/server/login"
	"github.com/itay2805/mcserver/server/play"
	"github.com/itay2805/mcserver/server/socket"
	"github.com/itay2805/mcserver/server/status"
	"github.com/panjf2000/ants"
	"io"
	"log"
	"net"
	"runtime/debug"
)

type PacketHandler func(player *game.Player, reader *minecraft.Reader)

var packetHandlers = [][]PacketHandler {
	socket.Handshaking: {
		0x00:	handshaking.HandleHandshaking,
	},
	socket.Status: {
		0x00: status.HandleRequest,
		0x01: status.HandlePing,
	},
	socket.Login: {
		0x00: login.HandleLoginStart,
	},
	socket.Play: {
		// 0x00
		// 0x01
		// 0x02
		// 0x03
		// 0x04
		0x05: play.HandleClientSettings,
		// 0x06
		// 0x07
		// 0x08
		// 0x09
		// 0x0A
		// 0x0B
		// 0x0C
		// 0x0D
		// 0x0E
		0x0F: play.HandleKeepAlive,
		//0x10:
		0x11: play.HandlePlayerPosition,
		0x12: play.HandlePlayerPositionAndRotation,
		0x13: play.HandlePlayerRotation,
		0x14: play.HandlePlayerMovement,
		// 0x15
		// 0x16
		// 0x17
		// 0x18
		0x19: play.HandlePlayerAbilities,
		0x1A: play.HandlePlayerDigging,
		0x1B: play.HandleEntityAction,
		// 0x1C
		// 0x1D
		// 0x1E
		// 0x1F
		// 0x20
		// 0x21
		// 0x22
		0x23: play.HandleHeldItemChange,
		// 0x24
		// 0x25
		0x26: play.HandleCreativeInventoryAction,
		// 0x27
		// 0x28
		// 0x29
		0x2A: play.HandleAnimation,
		// 0x2B
		0x2C: play.HandlePlayerBlockPlacement,
		// 0x2D
		// 0x2E
		// 0x2F
	},
}

func handleConnection(conn net.Conn) {
	// create the net player, which will be used
	// for the entity and for the socket
	player := game.NewPlayer(socket.NewSocket(conn))

	// do cleanup in here, this will make sure so even if we
	// get a panic in an handler everything will still work
	defer func() {
		// remove the player
		if player.State == socket.Play {
			player.State = socket.Disconnected
			game.LeftPlayer(player)
		}

		// close the connection
		player.Close()

		// recover from error
		if r := recover(); r != nil {
			log.Println("Got error", r, "\n" + string(debug.Stack()))
		}
	}()

	// start the send goroutine for this client
	go player.StartSend()

	// start the recv loop which will recv
	// packets and handle them correctly
	for {
		data, err := player.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// socket closed, ignore
			} else {
				log.Println("Got an error recving packet:", err)
			}
			break
		}

		reader := minecraft.Reader{
			Data:   data,
		}
		packetId := reader.ReadVarint()

		// check the id is even valid
		handlers := packetHandlers[player.State]
		if len(handlers) <= int(packetId) {
			log.Println("Got invalid packet id", packetId, "from", player, "state", player.State)
			continue
		}

		// check the handler is all good and call it if so
		if handlers[packetId] != nil {
			handlers[packetId](player, &reader)
		} else {
			log.Println("Got invalid packet id", packetId, "from", player, "state", player.State)
			continue
		}

		// check if we need to disconnect
		if player.State == socket.Disconnected {
			break
		}
	}
}

func StartServer() {
	log.Println("Starting server on :25565")

	// we are going to use ants for sending shit
	defer ants.Release()

	//
	// start the accept loop
	//
	server, err := net.Listen("tcp4", ":25565")
	if err != nil {
		log.Panicln("Got error listening:", err)
	}

	for {
		// accept a connection
		conn, err := server.Accept()
		if err != nil {
			log.Panicln("Got error accepting connection:", err)
		}

		// start the connection handling in
		// a new goroutine
		go handleConnection(conn)
	}
}
