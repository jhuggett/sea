package inbound

import (
	"encoding/json"

	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/outbound"
)

type Inbound struct {
	MoveShip struct {
		Req  MoveShipReq  `json:"req"`
		Resp MoveShipResp `json:"resp"`
	}

	Login struct {
		Req  LoginReq  `json:"req"`
		Resp LoginResp `json:"resp"`
	}

	Register struct {
		Req  RegisterReq  `json:"req"`
		Resp RegisterResp `json:"resp"`
	}

	GetWorldMap struct {
		Req  GetWorldMapReq  `json:"req"`
		Resp GetWorldMapResp `json:"resp"`
	}

	GetPorts struct {
		Req  GetPortsReq  `json:"req"`
		Resp GetPortsResp `json:"resp"`
	}

	ControlTime struct {
		Req  ControlTimeReq  `json:"req"`
		Resp ControlTimeResp `json:"resp"`
	}

	Trade struct {
		Req  TradeReq  `json:"req"`
		Resp TradeResp `json:"resp"`
	}

	PlotRoute struct {
		Req  PlotRouteReq  `json:"req"`
		Resp PlotRouteResp `json:"resp"`
	}

	HireCrew struct {
		Req  HireCrewReq  `json:"req"`
		Resp HireCrewResp `json:"resp"`
	}

	RepairShip struct {
		Req  RepairShipReq  `json:"req"`
		Resp RepairShipResp `json:"resp"`
	}

	ManageRoute struct {
		Req  ManageRouteReq  `json:"req"`
		Resp ManageRouteResp `json:"resp"`
	}
}

type InboundFunc func(req json.RawMessage) (interface{}, error)

type Connection interface {
	Context() *game_context.GameContext
	Sender() *outbound.Sender
}
