package inbound

import (
	"encoding/json"

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

	// GetWorldMapChunk struct {
	// 	Req  GetWorldMapChunkReq  `json:"req"`
	// 	Resp GetWorldMapChunkResp `json:"resp"`
	// }

	GetWorldMap struct {
		Req  GetWorldMapReq  `json:"req"`
		Resp GetWorldMapResp `json:"resp"`
	}
}

type InboundFunc func(req json.RawMessage) (interface{}, error)

type connection interface {
	Context() GameContext
	Sender() *outbound.Sender
}
