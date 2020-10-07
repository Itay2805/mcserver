package game

import (
	"github.com/google/uuid"
	"github.com/itay2805/mcserver/math"
	"github.com/itay2805/mcserver/minecraft"
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

type PlayerActionType int

const (
	PlayerActionNone = PlayerActionType(iota)
	PlayerActionDig
)

type PlayerAction struct {
	Type	PlayerActionType
	Data	interface{}
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
	LastKeepAlive		time.Time
	PingChanged			bool

	// The world we are in
	World 				*World

	// the chunks loaded by the client of this player
	loadedChunks 		map[world.ChunkPos]bool
	loadedEntities		map[int32]bool

	// player info
	waitingForPlayers	map[uuid.UUID]chan bool
	knownPlayers		map[uuid.UUID]bool
	joined				bool

	// sound related
	lastSoundPos		math.Point
	lastSoundTime		time.Time

	ActionQueue			*queue.Queue
}

func (p *Player) String() string  {
	if p.Player != nil {
		return p.Username
	} else {
		return p.RemoteAddr().String()
	}
}

func NewPlayer(socket socket.Socket) *Player {
	return &Player{
		Socket: socket,
		Player: nil,

		ViewDistance: 2,

		changeQueue: queue.New(),
		changeMutex: sync.Mutex{},

		Ping: -1,
		LastKeepAlive: time.Now(),
		PingChanged: false,

		World: nil,

		loadedChunks: make(map[world.ChunkPos]bool),
		loadedEntities: make(map[int32]bool),

		waitingForPlayers: make(map[uuid.UUID]chan bool),
		knownPlayers: make(map[uuid.UUID]bool),
		joined: true,

		ActionQueue: queue.New(),
	}
}

func (p *Player) Change(cb func()) {
	p.changeMutex.Lock()
	defer p.changeMutex.Unlock()
	p.changeQueue.Add(cb)
}

//
// Get a rect with the area that a player can see
//
func (p *Player) ViewRect() *math.Rect {
	return math.NewRectFromPoints(
		math.NewPoint(
			p.Position.X() - float64(p.ViewDistance * 16),
			0,
			p.Position.Z() - float64(p.ViewDistance * 16),
		),
		math.NewPoint(
			p.Position.X() + float64(p.ViewDistance * 16),
			256,
			p.Position.Z() + float64(p.ViewDistance * 16),
		),
	)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Sync server
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *Player) syncChanges() {
	// apply all changes
	p.changeMutex.Lock()
	for p.changeQueue.Length() > 0 {
		p.changeQueue.Remove().(func())()
	}
	p.changeMutex.Unlock()

	// if moved update the entity
	// position and bounding box
	if p.Moved {
		p.World.UpdateEntityPosition(p)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tick player
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *Player) tickDig(position minecraft.Position) {
	// TODO: don't assume creative

	// get the prev block
	blkstate := p.World.GetBlockState(position.X, position.Y, position.Z)

	// register a block change
	chunkX := position.X >> 4
	chunkZ := position.Z >> 4
	pos := world.ChunkPos{ X: chunkX, Z: chunkZ }
	p.World.BlockChanges[pos] = append(p.World.BlockChanges[pos], play.BlockRecord{
		BlockX:     byte(position.X & 0xf),
		BlockZ:     byte(position.Z & 0xf),
		BlockY:     byte(position.Y),
		BlockState: 0,
	})

	// send the sound and animation to other players which are in a 16 block range from
	// the given effect, farther than that they will just not be able to see it
	soundDistance := math.NewRect(position.ToPoint().SubScalar(16), [3]float64{ 32, 32, 32 })
	p.World.ForEntitiesInRange(soundDistance, func(ient entity.IEntity) {
		// ignore local player, their client does this anyways
		if ient == p {
			return
		}

		// ignore non-players
		if other, ok := ient.(*Player); ok {
			other.Send(play.Effect{
				EffectID:              play.EffectBlockBreak,
				Location:              position,
				Data:                  int32(blkstate),
				DisableRelativeVolume: false,
			})
		} else {
			log.Println("Not sending to", ient, reflect.TypeOf(ient))
		}
	})
}

func (p *Player) tick() {
	for p.ActionQueue.Length() > 0 {
		action := p.ActionQueue.Remove().(PlayerAction)
		switch action.Type {
			case PlayerActionDig:
				p.tickDig(action.Data.(minecraft.Position))
		}
	}
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
		// check about this channel
		select {
			case _, ok := <-c:
				if ok {
					close(c)
				}

				delete(p.waitingForPlayers, uid)
				p.knownPlayers[uid] = true
			default:
				// not sent yet, ignore
		}
	}

	// check for any new players
	if p.joined {
		pladd = make([]play.PIAddPlayer, 0, len(players))
		for _, pNew := range players {
			pladd = append(pladd, play.PIAddPlayer{
				UUID:        pNew.UUID,
				Name:        pNew.Username,
				Gamemode:    0,
				Ping:        int32(pNew.Ping.Milliseconds()),
			})
		}
	} else {
		// the player has a list of everyone, only send updates
		pladd = make([]play.PIAddPlayer, 0, len(newPlayers))
		plrem = make([]play.PIRemovePlayer, 0, len(leftPlayers))

		for _, pNew := range newPlayers {
			pladd = append(pladd, play.PIAddPlayer{
				UUID:        pNew.UUID,
				Name:        pNew.Username,
				Gamemode:    0,
				Ping:        int32(pNew.Ping.Milliseconds()),
			})
		}

		for _, pNew := range leftPlayers {
			plrem = append(plrem, play.PIRemovePlayer{UUID: pNew.UUID})
		}
	}

	// update latencies
	for _, pNew := range players {
		if p.PingChanged {
			pllat = append(pllat, play.PIUpdateLatency{
				UUID: pNew.UUID,
				Ping: int32(pNew.Ping.Milliseconds()),
			})
		}
	}

	if len(pladd) > 0 {
		// insert all the players we are waiting for to
		// the map
		done := make(chan bool)
		p.SendChan(play.PlayerInfo{ AddPlayer: pladd }, done)
		for _, pNew := range pladd {
			p.waitingForPlayers[pNew.UUID] = done
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
			changes := p.World.BlockChanges[world.ChunkPos{ X: x, Z: z }]
			if len(changes) > 0 {
				if len(changes) > 1 {
					// multi-block change
					p.Send(play.MultiBlockChange{
						ChunkX:  int32(x),
						ChunkZ:  int32(z),
						Records: changes,
					})
				} else {
					// single block change
					p.Send(play.BlockChange{
						Location: minecraft.Position{
							X: (x << 4) | int(changes[0].BlockX),
							Y: (z << 4) | int(changes[0].BlockZ),
							Z: int(changes[0].BlockY),
						},
						BlockID:  changes[0].BlockState,
					})
				}
			}
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
	// set all entities as not needed
	for eid := range p.loadedEntities {
		p.loadedEntities[eid] = false
	}

	// go over the entities and tick them
	p.World.ForEntitiesInRange(p.ViewRect(), func(ient entity.IEntity) {
		// skip current player
		if ient == p {
			return
		}

		// get the entity data
		e := ient.GetEntity()
		newEntity := false

		// first check if the entity is known or not
		if _, ok := p.loadedEntities[e.EID]; ok {

			// this entity is known, tick its position if needed
			if e.Moved && e.Rotated {
				p.Send(play.EntityTeleport{
					EntityID: e.EID,
					X:        e.Position.X(),
					Y:        e.Position.Y(),
					Z:        e.Position.Z(),
					Yaw:      e.Yaw,
					Pitch:    e.Pitch,
					OnGround: e.OnGround,
				})

				p.Send(play.EntityHeadLook{
					EntityID: e.EID,
					HeadYaw:  e.HeadYaw,
				})

			} else if e.Moved {

				p.Send(play.EntityTeleport{
					EntityID: e.EID,
					X:        e.Position.X(),
					Y:        e.Position.Y(),
					Z:        e.Position.Z(),
					Yaw:      e.Yaw,
					Pitch:    e.Pitch,
					OnGround: e.OnGround,
				})

			} else if e.Rotated {
				p.Send(play.EntityRotation{
					EntityID: e.EID,
					Yaw:      e.Yaw,
					Pitch:    e.Pitch,
					OnGround: e.OnGround,
				})
				p.Send(play.EntityHeadLook{
					EntityID: e.EID,
					HeadYaw:  e.HeadYaw,
				})
			} else {
				// nothing happend, just keep reminding the player about
				// this entity
				// TODO: do I really need to do this?
				p.Send(play.EntityMovement{
					EntityID: e.EID,
				})
			}

			// TODO: handle other entity specific stuff

		} else {
			newEntity = true

			// this entity is unknown, spawn it
			switch other := ient.(type) {
				case *Player:
					// check the client has the player info
					// about this player before sending it
					if _, ok := p.knownPlayers[other.UUID]; !ok {
						return
					}

					// This is a player, spawn it
					p.Send(play.SpawnPlayer{
						EntityID: other.EID,
						UUID:     other.UUID,
						X:        other.Position.X(),
						Y:        other.Position.Y(),
						Z:        other.Position.Z(),
						Yaw:      other.Yaw,
						Pitch:    other.Pitch,
					})

				default:
					// This is a mob
			}
		}

		// TODO: entity animation

		// TODO: entity equipment


		// check if there is metadata to update
		if newEntity || e.MetadataChanged {
			p.Send(play.EntityMetadata{
				EntityID: e.EID,
				Metadata: e,
			})
		}

		// set entity as known
		p.loadedEntities[e.EID] = true
	})

	// unload all uneeded entities
	ids := make([]int32, 0)
	for eid, val := range p.loadedEntities {
		if !val {
			delete(p.loadedEntities, eid)
			ids = append(ids, eid)
		}
	}

	// if there are entities to unload, unload them
	if len(ids) > 0 {
		p.Send(play.DestroyEntities{EntityIDs: ids})
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
// sync the client state
//
func (p *Player) syncClient() {
	p.syncPlayerInfo()
	p.syncChunks()
	p.syncEntities()

	// send keepalive if needed
	if time.Since(p.LastKeepAlive) > time.Second * 25 {
		now := time.Now()
		p.LastKeepAlive = now
		p.Send(play.KeepAlive{ KeepAliveId: now.UnixNano() })
	}
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