package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/models/ship"
	"github.com/jhuggett/sea/models/world_map"
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

	conn connection
}

func (e *RouteShip) Do(ticks uint64) (stop bool) {
	e.ticks += ticks

	speed := 3.0 // points per tick

	slog.Info("RouteShip", "id", e.ship.ID, "ticks", e.ticks, "route", len(e.route), "travelled", float64(e.ticks)*speed)

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

		err = e.conn.Sender().ShipMoved(e.ship.ID, world_map.Point{X: int(e.ship.X), Y: int(e.ship.Y)})
		if err != nil {
			slog.Error("RouteShip failed to notify", "id", e.ship.ID, "err", err)
		}
	}

	if len(e.route) == 0 {
		slog.Info("Ship reached destination", "id", e.ship.ID)
		return true
	}

	return false
}

func MoveShip(conn connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r MoveShipReq
		if err := json.Unmarshal(req, &r); err != nil {
			return nil, err
		}

		slog.Info("MoveShip", "id", conn.Context().ShipID, "x", r.X, "y", r.Y)

		worldMap, err := conn.Context().GameMap()
		if err != nil {
			return nil, err
		}

		ship, err := conn.Context().Ship()
		if err != nil {
			return nil, err
		}

		route, err := worldMap.PlotRoute(
			&world_map.Point{X: int(ship.X), Y: int(ship.Y)},
			&world_map.Point{X: int(r.X), Y: int(r.Y)},
		)
		if err != nil {
			return nil, err
		}

		conn.Context().Timeline.RegisterContinualEvent(
			&RouteShip{
				ship:  ship,
				route: route,
				conn:  conn,
			},
		)

		return MoveShipResp{
			Success: true,
		}, nil
	}
}
