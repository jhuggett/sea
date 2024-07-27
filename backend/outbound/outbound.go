package outbound

import (
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/jsonrpc"
)

type ExampleReq struct {
	Name string `json:"name"`
}

type ExampleResp struct {
	Age int `json:"age"`
}

type Outbound struct {
	Example struct {
		Req  ExampleReq  `json:"req"`
		Resp ExampleResp `json:"resp"`
	}

	ShipMoved struct {
		Req  ShipMovedReq  `json:"req"`
		Resp ShipMovedResp `json:"resp"`
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
