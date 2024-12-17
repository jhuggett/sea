package inbound

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/jhuggett/sea/log"
	continent_model "github.com/jhuggett/sea/models/continent"
	"github.com/jhuggett/sea/models/port"
	"github.com/jhuggett/sea/models/ship"
	"github.com/jhuggett/sea/timeline"
	"github.com/jhuggett/sea/utils/coordination"
	"gorm.io/gorm"
)

type MoveShipReq struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type MoveShipResp struct {
	Success bool `json:"success"`
}

type RouteShip struct {
	ship  *ship.Ship
	route []coordination.Point

	ticks uint64

	conn Connection

	started bool

	portToDockTo *port.Port
}

func (e *RouteShip) Do(ticks uint64) (stop bool) {
	slog.Info("RouteShip.Do", "ticks", ticks, "e.ticks", e.ticks)

	if !e.started {
		e.started = true
		e.ship.AnchorRaised()
	}

	e.ticks += ticks

	speed := 3.0 // points per tick

	slog.Debug("RouteShip", "id", e.ship.Persistent.ID, "ticks", e.ticks, "route", len(e.route), "travelled", float64(e.ticks)*speed)

	if float64(e.ticks)*speed > 1.0 {
		tilesMoved := int(float64(e.ticks) * speed)
		if tilesMoved > len(e.route) {
			tilesMoved = len(e.route)
		}
		slog.Info("RouteShip", "id", e.ship.Persistent.ID, "tilesMoved", tilesMoved)
		e.ship.Persistent.X = float64(e.route[tilesMoved-1].X)
		e.ship.Persistent.Y = float64(e.route[tilesMoved-1].Y)
		err := e.ship.Save()
		if err != nil {
			slog.Error("RouteShip failed to save", "id", e.ship.Persistent.ID, "err", err)
			return true
		}
		e.ticks = 0
		e.route = e.route[tilesMoved:]

		if e.ship.Persistent.IsDocked {
			e.ship.Undocked()
		}
		e.ship.Moved()
	}

	if len(e.route) == 0 {
		slog.Debug("Ship reached destination", "id", e.ship.Persistent.ID)

		if e.portToDockTo != nil {
			e.ship.AnchorLowered(ship.AnchorLoweredEventData{
				Location: ship.AnchorLoweredLocationPort,
			})
			e.ship.Docked()
		} else {
			e.ship.AnchorLowered(ship.AnchorLoweredEventData{
				Location: ship.AnchorLoweredLocationOpenSea,
			})
		}

		return true
	}

	return false
}

func MoveShip(conn Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		slog := slog.With("rid", log.RandID())

		var r MoveShipReq
		if err := json.Unmarshal(req, &r); err != nil {
			return nil, err
		}

		slog = slog.With("x", r.X, "y", r.Y)

		slog.Debug("MoveShip called")

		worldMap, err := conn.Context().GameMap()
		if err != nil {
			return nil, err
		}

		ship, err := conn.Context().Ship()
		if err != nil {
			return nil, err
		}

		starting := coordination.Point{X: int(ship.Persistent.X), Y: int(ship.Persistent.Y)}
		ending := coordination.Point{X: int(r.X), Y: int(r.Y)}

		obstacles := coordination.ObstacleMap{}

		if starting.SameAs(ending) {
			slog.Info("Starting and ending points are the same")
			return []coordination.Point{}, nil
		}

		slog.Debug("Plotting Route")
		slog.Debug("route", "starting", starting, "ending", ending)

		shipRouter := RouteShip{}

		for _, continent := range worldMap.Continents() {
			for _, p := range continent.Persistent.Points {
				obstacles.AddObstacle(p.Point())
			}

			land, err := continent.Contains(ending)
			if err != nil && !errors.Is(err, continent_model.ErrNotInContinent) {
				return nil, err
			}

			// slog.Debug("Checking if ending point is in a continent", "contains", contains, "pointInfo", pointInfo, "err", err)

			if land != nil && land.Coastal == true {
				slog.Debug("Ending point is a coastal point", "point", land)
				port, err := port.Find(*land)
				if errors.Is(err, gorm.ErrRecordNotFound) {
				} else if err != nil {
					return nil, err
				} else {
					slog.Debug("Ending point is a port", "port", port)
					shipRouter.portToDockTo = port
					obstacles.RemoveObstacle(coordination.Point{
						X: land.X,
						Y: land.Y,
					})
				}
			} else if land != nil {
				slog.Debug("Ending point is in a continent")
				return nil, errors.New("ending point is in a continent")
			}
		}

		route, err := worldMap.PlotRoute(
			starting,
			ending,
			obstacles,
		)
		if err != nil {
			return nil, err
		}

		if shipRouter.portToDockTo != nil {
			slog.Debug("Dropping last point in route", "route", route)
			route = route[:len(route)-1]
		}

		shipRouter.ship = ship
		shipRouter.route = route
		shipRouter.conn = conn

		conn.Context().Timeline.Do(func(routeShipEvent *RouteShip) timeline.EventDo {
			lastTickTimestamp := conn.Context().Timeline.CurrentTick()

			return func() uint64 {
				elapsedTicks := conn.Context().Timeline.CurrentTick() - lastTickTimestamp
				stop := routeShipEvent.Do(elapsedTicks)

				if stop {
					return 0
				}

				lastTickTimestamp = conn.Context().Timeline.CurrentTick()

				return conn.Context().Timeline.TicksPerCycle()
			}
		}(&shipRouter), 0)

		return MoveShipResp{
			Success: true,
		}, nil
	}
}
