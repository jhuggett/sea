package inbound

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/jhuggett/sea/log"
	"github.com/jhuggett/sea/models/port"
	"github.com/jhuggett/sea/models/ship"
	"github.com/jhuggett/sea/models/world_map"
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
	route []world_map.Point

	ticks uint64

	conn Connection

	portToDockTo *port.Port
}

func (e *RouteShip) Do(ticks uint64) (stop bool) {
	e.ticks += ticks

	speed := 3.0 // points per tick

	slog.Debug("RouteShip", "id", e.ship.ID, "ticks", e.ticks, "route", len(e.route), "travelled", float64(e.ticks)*speed)

	if float64(e.ticks)*speed > 1.0 {
		tilesMoved := int(float64(e.ticks) * speed)
		if tilesMoved > len(e.route) {
			tilesMoved = len(e.route)
		}
		slog.Info("RouteShip", "id", e.ship.ID, "tilesMoved", tilesMoved)
		e.ship.X = float64(e.route[tilesMoved-1].X)
		e.ship.Y = float64(e.route[tilesMoved-1].Y)
		err := e.ship.Save()
		if err != nil {
			slog.Error("RouteShip failed to save", "id", e.ship.ID, "err", err)
			return true
		}
		e.ticks = 0
		e.route = e.route[tilesMoved:]

		if e.ship.IsDocked {
			e.ship.Undocked()
		}
		e.ship.Moved()
	}

	if len(e.route) == 0 {
		slog.Debug("Ship reached destination", "id", e.ship.ID)

		if e.portToDockTo != nil {
			slog.Debug("Docking at port", "id", e.portToDockTo.ID, "coastalPoint", e.portToDockTo.CoastalPoint)

			e.ship.Docked()
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

		starting := &world_map.Point{X: int(ship.X), Y: int(ship.Y)}
		ending := &world_map.Point{X: int(r.X), Y: int(r.Y)}

		obstacles := world_map.ObstacleMap{}

		if starting.SameAs(ending) {
			slog.Info("Starting and ending points are the same")
			return []world_map.Point{}, nil
		}

		slog.Debug("Plotting Route")
		slog.Debug("route", "starting", *starting, "ending", *ending)

		shipRouter := RouteShip{}

		for _, continent := range worldMap.Continents {
			for _, coastalPoint := range continent.CoastalPoints {
				obstacles.AddObstacle(&world_map.Point{
					X: coastalPoint.X,
					Y: coastalPoint.Y,
				})
			}

			contains, pointInfo, err := continent.Contains(*ending)
			if err != nil {
				return nil, err
			}

			slog.Debug("Checking if ending point is in a continent", "contains", contains, "pointInfo", pointInfo, "err", err)

			if pointInfo.CoastalPoint != nil {
				slog.Debug("Ending point is a coastal point", "point", pointInfo.CoastalPoint)
				port, err := port.Find(*pointInfo.CoastalPoint)
				if errors.Is(err, gorm.ErrRecordNotFound) {
				} else if err != nil {
					return nil, err
				} else {
					slog.Debug("Ending point is a port", "port", port)
					shipRouter.portToDockTo = port
					obstacles.RemoveObstacle(&world_map.Point{
						X: pointInfo.CoastalPoint.X,
						Y: pointInfo.CoastalPoint.Y,
					})
				}
			} else if contains {
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

		conn.Context().Timeline.RegisterContinualEvent(&shipRouter)

		return MoveShipResp{
			Success: true,
		}, nil
	}
}
