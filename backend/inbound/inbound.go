package inbound

import (
	"encoding/json"
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
}

type InboundFunc func(req json.RawMessage) (interface{}, error)
