package outbound

import (
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/data/fleet"
	"github.com/jhuggett/sea/data/ship"
)

type ShipSummary struct {
	ID        uint `json:"id"`
	IsCapital bool `json:"isCapital"`

	Name string `json:"name"`
}

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

	Fleet []ShipSummary `json:"fleet"`
}

type ShipChangedResp struct{}

func (s *Sender) ShipChanged(shipID uint) error {

	ship, err := ship.Get(shipID)
	if err != nil {
		return err
	}

	fleet, err := fleet.GetFleetByShipID(shipID)
	if err != nil {
		return fmt.Errorf("failed to get fleet for ship %d: %w", shipID, err)
	}

	var fleetSummary []ShipSummary
	for _, s := range fleet.Persistent.Ships {
		fleetSummary = append(fleetSummary, ShipSummary{
			ID:        s.ID,
			IsCapital: fleet.Persistent.CapitalShipID == s.ID,
			Name:      s.Name,
		})
	}

	slog.Warn("fleetSummary:", "fleet", fleet, "fleetSummary", fleetSummary)

	inventory := ship.Inventory()

	totalCargoWeight := inventory.TotalWeight()

	totalUsedSpace := inventory.OccupiedSpace()

	sailingSpeed, err := ship.SailingSpeed()
	if err != nil {
		return err
	}

	_, err = s.Receiver.OnShipChanged(ShipChangedReq{
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

		Fleet: fleetSummary,
	})
	if err != nil {
		return err
	}

	return nil
}
