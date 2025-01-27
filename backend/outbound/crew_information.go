package outbound

import (
	"fmt"
)

// bad name, it's colliding with the person type in inbound, outta pull these types up probably or something
type OutboundPerson struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	NickName string `json:"nick_name"`

	Age int `json:"age"`

	Morale float64 `json:"morale"`
}

type EmploymentContract struct {
	Title string `json:"title"`

	StartDate uint `json:"start_date"`
	EndDate   uint `json:"end_date"`
}

type CrewMember struct {
	Person   OutboundPerson     `json:"person"`
	Contract EmploymentContract `json:"contract"`
}

type CrewInformationReq struct {
	Size    int     `json:"size"`
	Wage    float64 `json:"wage"`
	Rations float64 `json:"rations"`

	Morale float64 `json:"morale"`

	MinimumSafeManning uint `json:"minimumSafeManning"`
	MaximumSafeManning uint `json:"maximumSafeManning"`

	CrewMembers []CrewMember `json:"crew_members"`
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

	summary, err := crew.Summary()
	if err != nil {
		return fmt.Errorf("failed to get crew summary: %w", err)
	}

	members, err := crew.Members()
	if err != nil {
		return fmt.Errorf("failed to get crew members: %w", err)
	}

	crewMembers := []CrewMember{}
	for _, member := range members {
		crewMembers = append(crewMembers, CrewMember{
			Person: OutboundPerson{
				FirstName: member.FirstName,
				LastName:  member.LastName,
				NickName:  member.NickName,
				Age:       int(member.Age),
			},
			Contract: EmploymentContract{
				Title:     member.Title,
				StartDate: member.StartDate,
				EndDate:   member.EndDate,
			},
		})
	}

	_, err = s.rpc.Send("CrewInformation", CrewInformationReq{
		// Size:               crew.Persistent.Size,
		// Wage:               crew.Persistent.Wage,
		// Rations:            crew.Persistent.Rations,
		// Morale:             crew.Persistent.Morale,
		Size:               summary.Size,
		Morale:             summary.AverageMorale,
		MinimumSafeManning: ship.Persistent.MinimumSafeManning,
		MaximumSafeManning: ship.Persistent.MaximumSafeManning,

		CrewMembers: crewMembers,
	})
	if err != nil {
		return err
	}

	return nil
}
