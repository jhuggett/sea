package outbound

import (
	"fmt"
)

type CrewInformationReq struct {
	Size    int     `json:"size"`
	Wage    float64 `json:"wage"`
	Rations float64 `json:"rations"`

	Morale float64 `json:"morale"`

	MinimumSafeManning uint `json:"minimumSafeManning"`
	MaximumSafeManning uint `json:"maximumSafeManning"`
}

type CrewInformationResp struct{}

func (s *Sender) CrewInformation() error {
	ship, err := s.gameContext.Ship()
	if err != nil {
		return fmt.Errorf("failed to get ship: %w", err)
	}

	crew, err := ship.Crew()
	if err != nil {
		return fmt.Errorf("failed to get crew: %w", err)
	}

	_, err = s.rpc.Send("CrewInformation", CrewInformationReq{
		Size:               crew.Persistent.Size,
		Wage:               crew.Persistent.Wage,
		Rations:            crew.Persistent.Rations,
		Morale:             crew.Persistent.Morale,
		MinimumSafeManning: ship.Persistent.MinimumSafeManning,
		MaximumSafeManning: ship.Persistent.MaximumSafeManning,
	})
	if err != nil {
		return err
	}

	return nil
}
