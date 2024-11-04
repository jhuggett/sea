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
}

type InboundFunc func(req json.RawMessage) (interface{}, error)

type Connection interface {
	Context() *game_context.GameContext
	Sender() *outbound.Sender
}
