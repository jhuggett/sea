package inbound

import (
	"encoding/json"

	ship_model "github.com/jhuggett/sea/models/ship"
)

type MangeRouteAction string

const (
	ManageRouteActionStart MangeRouteAction = "start"
	ManageRouteActionPause MangeRouteAction = "pause"
	ManageRouteActionStop  MangeRouteAction = "cancel"
)

type ManageRouteReq struct {
	ShipID uint             `json:"ship_id"`
	Action MangeRouteAction `json:"action"` // start, pause, cancel
}

type ManageRouteResp struct {
}

func ManageRoute(conn Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var reqObj ManageRouteReq
		if err := json.Unmarshal(req, &reqObj); err != nil {
			return nil, err
		}

		route := ship_model.LookupRoute(reqObj.ShipID)

		switch reqObj.Action {
		case ManageRouteActionStart:
			route.Resume()
			//ship.StartRoute()
		case ManageRouteActionPause:
			route.Pause()
			//ship.PauseRoute()
		case ManageRouteActionStop:
			route.Cancel()
			//ship.CancelRoute()
		}

		return ManageRouteResp{}, nil
	}
}
