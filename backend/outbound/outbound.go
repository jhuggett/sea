package outbound

import "github.com/jhuggett/sea/jsonrpc"

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

	ShipChangedTarget struct {
		Req  ShipChangedTargetReq  `json:"req"`
		Resp ShipChangedTargetResp `json:"resp"`
	}
}

type Sender struct {
	rpc jsonrpc.JSONRPC
}

func NewSender(rpc jsonrpc.JSONRPC) *Sender {
	return &Sender{
		rpc: rpc,
	}
}
