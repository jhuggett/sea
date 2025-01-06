package outbound

import "github.com/jhuggett/sea/models/ship"

type ShipChangedReq struct {
	ID uint `json:"id"`

	X float64 `json:"x"`
	Y float64 `json:"y"`

	IsDocked bool `json:"isDocked"`

	StateOfRepair float64 `json:"stateOfRepair"`

	EstimatedSailingSpeed float64 `json:"estimatedSailingSpeed"`

	MinimumSafeManning uint `json:"minimumSafeManning"`
	MaximumSafeManning uint `json:"maximumSafeManning"`

	RecommendedMaxCargoWeightCapacity uint `json:"recommendedMaxCargoWeightCapacity"`
	MaxCargoSpaceCapacity             uint `json:"maxCargoSpaceCapacity"`

	CurrentCargoWeight float32 `json:"currentCargoWeight"`
	CurrentCargoSpace  float32 `json:"currentCargoSpace"`
}

type ShipChangedResp struct{}

func (s *Sender) ShipChanged(shipID uint) error {

	ship, err := ship.Get(shipID)
	if err != nil {
		return err
	}

	inventory := ship.Inventory()

	totalCargoWeight := inventory.TotalWeight()

	totalUsedSpace := inventory.OccupiedSpace()

	sailingSpeed, err := ship.SailingSpeed()
	if err != nil {
		return err
	}

	_, err = s.rpc.Send("ShipChanged", ShipChangedReq{
		ID:                                shipID,
		X:                                 ship.Persistent.X,
		Y:                                 ship.Persistent.Y,
		IsDocked:                          ship.Persistent.IsDocked,
		StateOfRepair:                     ship.Persistent.StateOfRepair,
		EstimatedSailingSpeed:             sailingSpeed,
		MinimumSafeManning:                ship.Persistent.MinimumSafeManning,
		MaximumSafeManning:                ship.Persistent.MaximumSafeManning,
		RecommendedMaxCargoWeightCapacity: ship.Persistent.RecommendedMaxCargoWeightCapacity,
		MaxCargoSpaceCapacity:             ship.Persistent.MaxCargoSpaceCapacity,
		CurrentCargoWeight:                totalCargoWeight,
		CurrentCargoSpace:                 totalUsedSpace,
	})
	if err != nil {
		return err
	}

	return nil
}
