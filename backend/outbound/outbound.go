package outbound

import (
	"github.com/jhuggett/sea/game_context"
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

func NewSender(gameContext *game_context.GameContext, receiver Receiver) *Sender {
	return &Sender{
		gameContext: gameContext,
		Receiver:    receiver,
	}
}
