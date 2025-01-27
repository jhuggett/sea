package inbound

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/jhuggett/sea/data/port"
	ship_model "github.com/jhuggett/sea/data/ship"
	"github.com/jhuggett/sea/log"
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

func Plot(start coordination.Point, end coordination.Point, conn Connection) ([]coordination.Point, *port.Port, error) {
	starting := start
	ending := end

	obstacles := coordination.ObstacleMap{}

	if starting.SameAs(ending) {
		return []coordination.Point{}, nil, nil
	}

	worldMap, err := conn.Context().GameMap()
	if err != nil {
		return nil, nil, err
	}

	var portToDockTo *port.Port

	/*

		- Need to check if ending is a Point
			- if it is, check if there is a port there
				- if there is, dock there


		- Need to build an obstacle map of all the coastal points
		  - need to able to query all the coastal points of a world map
	*/

	// ---

	land, err := worldMap.HasLand(ending.X, ending.Y)
	if err != nil {
		return nil, nil, err
	}

	if land != nil {
		if land.Coastal == true {
			port, err := port.Find(*land)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil, errors.New("no port found")
			} else if err != nil {
				return nil, nil, err
			} else {
				portToDockTo = port
			}
		} else {
			return nil, nil, errors.New("ending point is in a continent")
		}
	}

	coastalPoints, err := worldMap.CoastalPoints()
	if err != nil {
		return nil, nil, err
	}

	for _, p := range coastalPoints {
		obstacles.AddObstacle(p.Point())
	}

	if portToDockTo != nil {
		obstacles.RemoveObstacle(coordination.Point{
			X: land.X,
			Y: land.Y,
		})
	}

	route, err := worldMap.PlotRoute(
		starting,
		ending,
		obstacles,
	)
	if err != nil {
		return nil, nil, err
	}

	if portToDockTo != nil {
		route = route[:len(route)-1]
	}

	return route, portToDockTo, nil
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

		ship, err := conn.Context().Ship()
		if err != nil {
			return nil, err
		}

		shipRouter := ship_model.RouteShip{}

		route, portToDockTo, err := Plot(
			ship.Location(),
			coordination.Point{
				X: int(r.X),
				Y: int(r.Y),
			},
			conn,
		)
		if err != nil {
			return nil, err
		}

		shipRouter.PortToDockTo = portToDockTo

		shipRouter.Ship = ship
		shipRouter.Route = route
		// shipRouter.conn = conn

		shipRouter.Broadcast = func(routeShipEvent *ship_model.RouteShip) {
			conn.Sender().ShipMoved(
				ship.Persistent.ID,
			)
		}

		conn.Context().Timeline.Do(func(routeShipEvent *ship_model.RouteShip) timeline.EventDo {
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

		ship_model.RegisterRoute(&shipRouter)

		return MoveShipResp{
			Success: true,
		}, nil
	}
}
