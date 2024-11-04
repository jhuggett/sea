package outbound

import (
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/jsonrpc"
)

type Outbound struct {
	ShipMoved struct {
		Req  ShipMovedReq  `json:"req"`
		Resp ShipMovedResp `json:"resp"`
	}

	ShipDocked struct {
		Req  ShipDockedReq  `json:"req"`
		Resp ShipDockedResp `json:"resp"`
	}

	TimeChanged struct {
		Req  TimeChangedReq  `json:"req"`
		Resp TimeChangedResp `json:"resp"`
	}
}

type Sender struct {
	rpc         jsonrpc.JSONRPC
	gameContext *game_context.GameContext
}

func NewSender(rpc jsonrpc.JSONRPC, gameContext *game_context.GameContext) *Sender {
	return &Sender{
		rpc:         rpc,
		gameContext: gameContext,
	}
}
