package game_context

import (
	"github.com/jhuggett/sea/models/ship"
	"github.com/jhuggett/sea/models/world_map"
	"github.com/jhuggett/sea/timeline"
)

type GameContext struct {
	snapshot Snapshot

	Timeline *timeline.Timeline
}

type Snapshot struct {
	ShipID    uint
	GameMapID uint

	// Sign this for auth in future
}

func (g *GameContext) Ship() (*ship.Ship, error) {
	return ship.Get(g.snapshot.ShipID)
}

func (g *GameContext) GameMap() (*world_map.WorldMap, error) {
	return world_map.Get(g.snapshot.GameMapID)
}

func (g *GameContext) GameMapID() uint {
	return g.snapshot.GameMapID
}

func (g *GameContext) ShipID() uint {
	return g.snapshot.ShipID
}

func New(snapshot Snapshot) *GameContext {
	return &GameContext{
		snapshot: snapshot,
	}
}
