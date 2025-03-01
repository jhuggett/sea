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

	ShipInventoryChanged struct {
		Req  ShipInventoryChangedReq  `json:"req"`
		Resp ShipInventoryChangedResp `json:"resp"`
	}

	CrewInformation struct {
		Req  CrewInformationReq  `json:"req"`
		Resp CrewInformationResp `json:"resp"`
	}

	ShipChanged struct {
		Req  ShipChangedReq  `json:"req"`
		Resp ShipChangedResp `json:"resp"`
	}
}

type Receiver struct {
	OnShipMoved            func(ShipMovedReq) (ShipMovedResp, error)
	OnShipDocked           func(ShipDockedReq) (ShipDockedResp, error)
	OnTimeChanged          func(TimeChangedReq) (TimeChangedResp, error)
	OnShipInventoryChanged func(ShipInventoryChangedReq) (ShipInventoryChangedResp, error)
	OnCrewInformation      func(CrewInformationReq) (CrewInformationResp, error)
	OnShipChanged          func(ShipChangedReq) (ShipChangedResp, error)
}

type Sender struct {
	gameContext *game_context.GameContext
	Receiver    Receiver
}

func NewRPCReceiver(rpc jsonrpc.JSONRPC) *Receiver {
	return &Receiver{
		OnShipMoved: func(req ShipMovedReq) (ShipMovedResp, error) {
			_, err := rpc.Send("ShipMoved", req)
			return ShipMovedResp{}, err
		},
		OnShipDocked: func(req ShipDockedReq) (ShipDockedResp, error) {
			_, err := rpc.Send("ShipDocked", req)
			return ShipDockedResp{}, err
		},
		OnTimeChanged: func(req TimeChangedReq) (TimeChangedResp, error) {
			_, err := rpc.Send("TimeChanged", req)
			return TimeChangedResp{}, err
		},
		OnShipInventoryChanged: func(req ShipInventoryChangedReq) (ShipInventoryChangedResp, error) {
			_, err := rpc.Send("ShipInventoryChanged", req)
			return ShipInventoryChangedResp{}, err
		},
		OnCrewInformation: func(req CrewInformationReq) (CrewInformationResp, error) {
			_, err := rpc.Send("CrewInformation", req)
			return CrewInformationResp{}, err
		},
		OnShipChanged: func(req ShipChangedReq) (ShipChangedResp, error) {
			_, err := rpc.Send("ShipChanged", req)
			return ShipChangedResp{}, err
		},
	}
}

func NewSender(gameContext *game_context.GameContext, receiver Receiver) *Sender {
	return &Sender{
		gameContext: gameContext,
		Receiver:    receiver,
	}
}
