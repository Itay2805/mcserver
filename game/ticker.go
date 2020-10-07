package game

import (
	"container/list"
)

type tickObject struct {
	targetTick	int
	callback  	func()
}

//
// The current tick we are processing
//
var currentTick = 0

//
// The objects that need to be ticked
//
var tickedObjects = list.New()

//
// Inserts a new ticked object to the list, finds
// a place where the tick need to be placed
//
func insertTickObject(el tickObject) {
	for e := tickedObjects.Front(); e != nil; e = e.Next() {
		val := e.Value.(tickObject)
		if val.targetTick > el.targetTick {
			tickedObjects.InsertBefore(el, e)
			return
		}
	}
	tickedObjects.PushBack(el)
}

//
// Tick all the objects that need to be ticked
// on this tick
//
func tickScheduledObjects() {
	var next *list.Element
	for e := tickedObjects.Front(); e != nil; e = next {
		next = e.Next()
		val := e.Value.(tickObject)

		if val.targetTick == currentTick {
			val.callback()
			tickedObjects.Remove(e)
		} else {
			break
		}
	}
}

func RegisterForTick(cb func(), ticks int) {
	insertTickObject(tickObject{
		targetTick: currentTick + ticks,
		callback:   cb,
	})
}
