package ship

import (
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/models/port"
	"github.com/jhuggett/sea/timeline"
	"github.com/jhuggett/sea/utils/coordination"
)

var shipRoutesInProgress = map[uint]*RouteShip{}

func RegisterRoute(route *RouteShip) {
	if shipRoutesInProgress[route.Ship.Persistent.ID] != nil {
		shipRoutesInProgress[route.Ship.Persistent.ID].cancelled = true
	}

	shipRoutesInProgress[route.Ship.Persistent.ID] = route
}

func LookupRoute(shipId uint) *RouteShip {
	return shipRoutesInProgress[shipId]
}

type RouteShip struct {
	Ship  *Ship
	Route []coordination.Point

	ticks          uint64
	lastTilesMoved int

	started bool

	PortToDockTo *port.Port

	cancelled bool
	paused    bool

	Broadcast func(e *RouteShip)
}

func (e *RouteShip) CurrentSpeed() (float64, error) {
	return e.Ship.SailingSpeed()
}

func (e *RouteShip) TilesMoved() (int, error) {
	speed, err := e.CurrentSpeed()
	if err != nil {
		slog.Error("RouteShip failed to get speed", "id", e.Ship.Persistent.ID, "err", err)
		return 0, err
	}

	tilesMoved := int(float64(e.ticks) * (speed / float64(timeline.Day)))

	if tilesMoved > len(e.Route) {
		tilesMoved = len(e.Route)
	}

	return tilesMoved, nil
}

func (e *RouteShip) EstimatedTimeLeft() (uint64, error) {
	speed, err := e.CurrentSpeed()
	if err != nil {
		slog.Error("RouteShip failed to get speed", "id", e.Ship.Persistent.ID, "err", err)
		return 0, err
	}

	if speed == 0 {
		return 0, fmt.Errorf("speed is 0")
	}

	tilesMoved, err := e.TilesMoved()
	if err != nil {
		slog.Error("RouteShip failed to get tilesMoved", "id", e.Ship.Persistent.ID, "err", err)
		return 0, err
	}

	return uint64((float64(len(e.Route) - tilesMoved)) / speed), nil
}

func (e *RouteShip) Do(ticks uint64) (stop bool) {
	if e.cancelled {
		return true
	}

	if e.paused {
		return false
	}

	if !e.started {
		e.started = true
		e.Ship.AnchorRaised()
	}

	e.ticks += ticks

	// slog.Debug("RouteShip", "speed", speed, "id", e.Ship.Persistent.ID, "ticks", e.ticks, "route", len(e.Route), "travelled", float64(e.ticks)*speed)

	tilesMoved, err := e.TilesMoved()
	if err != nil {
		slog.Error("RouteShip failed to get tilesMoved", "id", e.Ship.Persistent.ID, "err", err)
		return true
	}

	if tilesMoved > e.lastTilesMoved {
		e.lastTilesMoved = tilesMoved

		e.Ship, err = e.Ship.Fetch()
		if err != nil {
			slog.Error("RouteShip failed to fetch ship", "id", e.Ship.Persistent.ID, "err", err)
			return true
		}

		slog.Info("RouteShip", "id", e.Ship.Persistent.ID, "tilesMoved", tilesMoved)
		e.Ship.Persistent.X = float64(e.Route[tilesMoved-1].X)
		e.Ship.Persistent.Y = float64(e.Route[tilesMoved-1].Y)
		err := e.Ship.Save()
		if err != nil {
			slog.Error("RouteShip failed to save", "id", e.Ship.Persistent.ID, "err", err)
			return true
		}
		// e.ticks = 0
		// e.Route = e.Route[tilesMoved:]

		if e.Ship.Persistent.IsDocked {
			e.Ship.Undocked()
		}

		e.Broadcast(e)
		e.Ship.Moved()
	}

	if len(e.Route) == tilesMoved {
		slog.Debug("Ship reached destination", "id", e.Ship.Persistent.ID)

		if e.PortToDockTo != nil {
			e.Ship.AnchorLowered(AnchorLoweredEventData{
				Location: AnchorLoweredLocationPort,
			})
			slog.Info("Ship docked", "id", e.Ship.Persistent.ID)
			e.Ship.Docked(e.PortToDockTo.Persistent.Point.Point())
		} else {
			slog.Info("Ship reached destination - not port", "id", e.Ship.Persistent.ID)
			e.Ship.AnchorLowered(AnchorLoweredEventData{
				Location: AnchorLoweredLocationOpenSea,
			})
		}

		return true
	}

	return false
}

func (e *RouteShip) Pause() {
	e.paused = true
	e.Ship.AnchorLowered(AnchorLoweredEventData{
		Location: AnchorLoweredLocationOpenSea,
	})
	e.Broadcast(e)
}

func (e *RouteShip) Resume() {
	e.paused = false
	e.Ship.AnchorRaised()
	e.Broadcast(e)
}

func (e *RouteShip) Cancel() {
	e.cancelled = true
	e.Ship.AnchorLowered(AnchorLoweredEventData{
		Location: AnchorLoweredLocationOpenSea,
	})
	e.Broadcast(e)
}

func (e *RouteShip) IsPaused() bool {
	return e.paused
}

func (e *RouteShip) IsCancelled() bool {
	return e.cancelled
}
