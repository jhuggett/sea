package outbound

import (
	"fmt"

	ship_model "github.com/jhuggett/sea/models/ship"
	"github.com/jhuggett/sea/timeline"
	"github.com/jhuggett/sea/utils/coordination"
)

type RouteInfo struct {
	TotalTilesMoved    int  `json:"total_tiles_moved"`
	TilesInRoute       int  `json:"tiles_in_route"`
	ReachedDestination bool `json:"reached_destination"`

	EstimatedTimeLeft uint64 `json:"estimated_time_left"` // in days

	Trajectory []coordination.Point `json:"trajectory"`

	IsPaused    bool `json:"is_paused"`
	IsCancelled bool `json:"is_cancelled"`

	HeadedToPort bool `json:"headed_to_port"`
}

type ShipMovedReq struct {
	ShipID   uint               `json:"ship_id"`
	Location coordination.Point `json:"location"`

	RouteInfo RouteInfo `json:"route_info"`
}

type ShipMovedResp struct{}

func (s *Sender) ShipMoved(shipId uint) error {
	// slog.Info("ShipMoved", "id", shipId, "location", location)

	ship, err := ship_model.Get(shipId)
	if err != nil {
		return err
	}

	route := ship_model.LookupRoute(shipId)
	if route == nil {
		return fmt.Errorf("no route found for ship %d", shipId)
	}

	totalTilesMoved, err := route.TilesMoved()
	if err != nil {
		return err
	}

	reachedDestination := totalTilesMoved >= len(route.Route)

	estimatedTimeLeft, err := route.EstimatedTimeLeft()
	if err != nil {
		return err
	}

	_, err = s.rpc.Send("ShipMoved", ShipMovedReq{
		ShipID:   shipId,
		Location: ship.Location(),

		RouteInfo: RouteInfo{
			TotalTilesMoved:    totalTilesMoved,
			TilesInRoute:       len(route.Route),
			ReachedDestination: reachedDestination,

			EstimatedTimeLeft: estimatedTimeLeft / timeline.Day,

			Trajectory:  route.Route[totalTilesMoved:],
			IsPaused:    route.IsPaused(),
			IsCancelled: route.IsCancelled(),

			HeadedToPort: route.PortToDockTo != nil,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
