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
		//0x02: play.HandleChatMessage,
		//0x03: play.HandleClientStatus,
		0x05: play.HandleClientSettings,
		//0x09: play.HandlePluginMessage,
		//
		//0x0B: play.HandleKeepAlive,
		0x11: play.HandlePlayerPosition,
		0x12: play.HandlePlayerPositionAndRotation,
		0x13: play.HandlePlayerRotation,
		0x14: play.HandlePlayerMovement,
		//
		//0x13: play.HandlePlayerAbilities,
		//0x14: play.HandlePlayerDigging,
		//0x15: play.HandleEntityAction,
		//
		//0x1A: play.HandleHeldItemChange,
		//0x1B: play.HandleCreativeInventoryAction,
		//
		//0x1D: play.HandleAnimation,
		//
		//0x1F: play.HandlePlayerBlockPlacement,
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
