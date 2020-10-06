package game

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/itay2805/mcserver/minecraft/entity"
	"github.com/itay2805/mcserver/minecraft/proto/play"
	"github.com/itay2805/mcserver/minecraft/world"
	"github.com/itay2805/mcserver/server/socket"
	"github.com/panjf2000/ants"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/eapache/queue"
)

type PendingChange struct {
	Field		interface{}
	Value		interface{}
	ChangeFlag	interface{}
}

type Player struct {
	// the socket
	socket.Socket

	// The entity player
	*entity.Player

	// Client settings we need to know about
	ViewDistance		int

	// pending changes
	changeQueue			*queue.Queue
	changeMutex			sync.Mutex

	// ping related
	Ping				time.Duration
	PingChanged			bool

	// The world we are in
	World 				*World

	// the chunks loaded by the client of this player
	loadedChunks 		map[world.ChunkPos]bool

	// player info
	waitingForPlayers	map[uuid.UUID]chan bool
	knownPlayers		map[uuid.UUID]bool
	joined				bool
}

func (p *Player) String() string  {
	username := "<none>"
	uuid := "<none>"
	if p.Player != nil {
		username = p.Username
		uuid = p.UUID.String()
	}
	return fmt.Sprintf("Player{ Username: %s, UUID: %s, Socket: %s }", username, uuid, p.RemoteAddr())
}

func NewPlayer(socket socket.Socket) *Player {
	return &Player{
		Socket: socket,
		Player: nil,

		ViewDistance: 2,

		changeQueue: queue.New(),
		changeMutex: sync.Mutex{},

		World: nil,

		loadedChunks: make(map[world.ChunkPos]bool),

		waitingForPlayers: make(map[uuid.UUID]chan bool),
		knownPlayers: make(map[uuid.UUID]bool),
		joined: true,
	}
}

func (p *Player) Change(field, value, flag interface{}) {
	p.changeMutex.Lock()
	defer p.changeMutex.Unlock()

	// TODO: in debug only
	if reflect.ValueOf(field).Elem().Kind() != reflect.ValueOf(value).Kind() {
		log.Panicln("Can't assign", reflect.TypeOf(value), "to", reflect.TypeOf(field))
	}

	p.changeQueue.Add(PendingChange{
		Field:      field,
		Value:      value,
		ChangeFlag: flag,
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Sync server
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *Player) syncChanges() {
	p.changeMutex.Lock()
	defer p.changeMutex.Unlock()

	// apply all changes
	for p.changeQueue.Length() > 0 {
		change := p.changeQueue.Remove().(PendingChange)
		value := reflect.ValueOf(change.Field)

		// check if the value has changed
		if change.ChangeFlag != nil && value.Elem().Interface() != change.Value {
			reflect.ValueOf(change.ChangeFlag).Elem().SetBool(true)
		}

		// set the value
		value.Elem().Set(reflect.ValueOf(change.Value))
	}

	// if moved update the entity
	// position and bounding box
	if p.Moved {
		p.World.UpdateEntityPosition(p)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tick player
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *Player) tick() {
	// TODO: do shit
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Sync client
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
// Sync the player infos, needed before we can send
//
func (p *Player) syncPlayerInfo() {
	var pladd []play.PIAddPlayer
	var plrem []play.PIRemovePlayer
	var pllat []play.PIUpdateLatency

	// check for any players that we already got
	for uid, c := range p.waitingForPlayers {
		select {
		case _, ok := <-c:
			if ok {
				// we are the first to see this change
				close(c)
			}

			// we know about this player
			delete(p.waitingForPlayers, uid)
			p.knownPlayers[uid] = true
		default:
			// not sent yet, ignore
		}
	}

	// check for any new players
	if p.joined {
		pladd = make([]play.PIAddPlayer, 0, len(players))
		for _, p := range players {
			pladd = append(pladd, play.PIAddPlayer{
				UUID:        p.UUID,
				Name:        p.Username,
				Gamemode:    0,
				Ping:        int32(p.Ping.Milliseconds()),
			})
		}
	} else {
		// the player has a list of everyone, only send updates
		pladd = make([]play.PIAddPlayer, 0, len(newPlayers))
		plrem = make([]play.PIRemovePlayer, 0, len(leftPlayers))

		for _, p := range newPlayers {
			pladd = append(pladd, play.PIAddPlayer{
				UUID:        p.UUID,
				Name:        p.Username,
				Gamemode:    0,
				Ping:        int32(p.Ping),
			})
		}

		for _, p := range leftPlayers {
			plrem = append(plrem, play.PIRemovePlayer{UUID: p.UUID})
		}
	}

	// update latencies
	for _, p := range players {
		if p.PingChanged {
			pllat = append(pllat, play.PIUpdateLatency{
				UUID: p.UUID,
				Ping: int32(p.Ping),
			})
		}
	}

	if len(pladd) > 0 {
		// insert all the players we are waiting for to
		// the map
		done := make(chan bool)
		p.SendChan(play.PlayerInfo{ AddPlayer: pladd }, done)
		for _, p := range newPlayers {
			p.waitingForPlayers[p.UUID] = done
		}
	}

	if len(plrem) > 0 {
		p.Send(play.PlayerInfo{ RemovePlayer: plrem })
	}

	if len(pllat) > 0 {
		p.Send(play.PlayerInfo{ UpdateLatency: pllat })
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *Player) syncChunks() {
	// first set all the chunks as uneeded
	for pos := range p.loadedChunks {
		p.loadedChunks[pos] = false
	}

	// go over all the chunks in the view distance
	forEachChunkInRange(int(p.Position.X()) >> 4, int(p.Position.Z()) >> 4, p.ViewDistance,
	func(x, z int) {
		pos := world.ChunkPos{X: x, Z: z}

		if _, ok := p.loadedChunks[pos]; !ok {
			// this chunk is not loaded at the player
			// load the chunk in an async manner
			_ = ants.Submit(func() {
				p.World.SendChunkToPlayer(x, z, p)
			})

		} else {
			// TODO: handle block updates
		}

		// mark this chunk as needed
		p.loadedChunks[pos] = true
	})

	// now any chunk that is still loaded client side
	// and is not needed send an unload packet
	for pos, val := range p.loadedChunks {
		if !val {
			delete(p.loadedChunks, pos)

			p.Send(play.UnloadChunk{
				ChunkX: int32(pos.X),
				ChunkZ: int32(pos.Z),
			})
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *Player) syncEntities() {

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
// sync the client state
//
func (p *Player) syncClient() {
	p.syncPlayerInfo()
	p.syncChunks()
	p.syncEntities()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Cleanup
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *Player) cleanupTick() {
	p.Moved = false
	p.Rotated = false
	p.MetadataChanged = false
	p.joined = false
}