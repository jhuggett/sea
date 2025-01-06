package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/timeline"
	"github.com/jhuggett/sea/utils/coordination"
)

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type PlotRouteReq struct {
	Coordinate Coordinate `json:"coordinate"`
}

type PlotRouteResp struct {
	Coordinates []Coordinate `json:"coordinates"`

	EstimatedSailingSpeed float64 `json:"speed"`
	EstimatedDuration     float64 `json:"duration"` // in days
}

func PlotRoute(conn Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var reqObj PlotRouteReq
		if err := json.Unmarshal(req, &reqObj); err != nil {
			return nil, err
		}

		ship, err := conn.Context().Ship()
		if err != nil {
			return nil, err
		}

		route, _, err := Plot(
			ship.Location(),
			coordination.Point{
				X: int(reqObj.Coordinate.X),
				Y: int(reqObj.Coordinate.Y),
			},
			conn,
		)

		if err != nil {
			slog.Error("Failed to plot route: ", err)
			return nil, err
		}

		var coordinates []Coordinate

		for _, point := range route {
			coordinates = append(coordinates, Coordinate{
				X: float64(point.X),
				Y: float64(point.Y),
			})
		}

		respObj := PlotRouteResp{
			Coordinates: coordinates,
		}

		sailingSpeed, err := ship.SailingSpeed()
		if err != nil {
			slog.Error("Failed to get sailing speed: ", err)
			return nil, err
		}

		if sailingSpeed == 0 {
			sailingSpeed = 0.001
		}

		slog.Debug("Sailing speed: ", sailingSpeed)
		slog.Debug("Route length: ", len(route))
		slog.Debug("Estimated duration: ", float64(len(route))/sailingSpeed/float64(timeline.Day))

		respObj.EstimatedSailingSpeed = sailingSpeed
		respObj.EstimatedDuration = float64(len(route)) / sailingSpeed / float64(timeline.Day)

		return respObj, nil
	}
}
