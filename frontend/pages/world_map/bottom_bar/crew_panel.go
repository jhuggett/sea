package bottom_bar

import (
	"design-library/config"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"fmt"

	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/sea/outbound"
)

func NewCrewPanel(manager *game.Manager) *CrewPanel {
	crewPanel := &CrewPanel{
		Manager: manager,
	}
	return crewPanel
}

type CrewPanel struct {
	doodad.Default

	Manager *game.Manager

	CrewData outbound.CrewInformationReq

	Subbed string
}

func (c *CrewPanel) Setup() {
	// Subscribe to crew updates if not already subscribed
	if c.Subbed == "" {
		// Add callback for crew information changes
		c.Subbed = c.Manager.OnCrewInformationCallback.Add(func(cir outbound.CrewInformationReq) error {
			c.CrewData = cir
			doodad.ReSetup(c)
			return nil
		})

	}

	if c.CrewData.Size == 0 {
		// Request initial crew information
		go c.Manager.PlayerShip.TriggerCrewInfoRequest()
		return // No crew data available yet
	}

	// Main container stack
	mainStack := stack.New(stack.Config{})
	c.AddChild(mainStack)

	// Title
	mainStack.AddChild(
		label.New(label.Config{
			Message: "Ship Crew",
		}),
	)

	// If no crew data yet, show a message
	if len(c.CrewData.CrewMembers) == 0 {
		mainStack.AddChild(
			label.New(label.Config{
				Message: "No crew data available yet",
			}),
		)
		c.Children().Setup()
		return
	}

	// Display general crew information
	infoStack := stack.New(stack.Config{})
	mainStack.AddChild(infoStack)

	infoStack.AddChild(
		label.New(label.Config{
			Message: fmt.Sprintf("Crew Size: %d", c.CrewData.Size),
		}),
	)

	infoStack.AddChild(
		label.New(label.Config{
			Message: fmt.Sprintf("Morale: %.1f%%", c.CrewData.Morale*100),
		}),
	)

	// Safe manning info
	infoStack.AddChild(
		label.New(label.Config{
			Message: fmt.Sprintf("Safe Manning: %d-%d", c.CrewData.MinimumSafeManning, c.CrewData.MaximumSafeManning),
		}),
	)

	// Create header for crew members table
	headerStack := stack.New(stack.Config{
		Flow:            config.LeftToRight,
		BackgroundColor: colors.SemiTransparent(colors.Panel),

		SpaceBetween: 10,
		Padding: config.Padding{
			Top: 15,
		},
	})
	mainStack.AddChild(headerStack)

	// Add headers
	headerStack.AddChild(
		label.New(label.Config{
			Message: "Name",
			Layout:  box.Zeroed().SetWidth(150),
		}),
	)

	headerStack.AddChild(
		label.New(label.Config{
			Message: "Title",
			Layout:  box.Zeroed().SetWidth(100),
		}),
	)

	headerStack.AddChild(
		label.New(label.Config{
			Message: "Age",
			Layout:  box.Zeroed().SetWidth(50),
		}),
	)

	headerStack.AddChild(
		label.New(label.Config{
			Message: "Morale",
			Layout:  box.Zeroed().SetWidth(70),
		}),
	)

	// Create a scrollable list for crew members
	crewContainer := stack.New(stack.Config{})
	mainStack.AddChild(crewContainer)

	// Add each crew member to the list
	for _, member := range c.CrewData.CrewMembers {
		crewRow := stack.New(stack.Config{
			Flow:         config.LeftToRight,
			SpaceBetween: 10,
		})
		crewContainer.AddChild(crewRow)

		// Crew member name
		name := member.Person.FirstName + " " + member.Person.LastName
		if member.Person.NickName != "" {
			name += " '" + member.Person.NickName + "'"
		}

		crewRow.AddChild(
			label.New(label.Config{
				Message: name,
				Layout:  box.Zeroed().SetWidth(150),
			}),
		)

		// Crew member title
		crewRow.AddChild(
			label.New(label.Config{
				Message: member.Contract.Title,
				Layout:  box.Zeroed().SetWidth(100),
			}),
		)

		// Crew member age
		crewRow.AddChild(
			label.New(label.Config{
				Message: fmt.Sprintf("%d", member.Person.Age),
				Layout:  box.Zeroed().SetWidth(50),
			}),
		)

		// Crew member morale
		crewRow.AddChild(
			label.New(label.Config{
				Message: fmt.Sprintf("%.1f%%", member.Person.Morale*100),
				Layout:  box.Zeroed().SetWidth(70),
			}),
		)
	}

	c.Children().Setup()

	c.Layout().Recalculate()
}

func (c *CrewPanel) Teardown() error {
	if c.Subbed != "" {
		c.Manager.OnCrewInformationCallback.Remove(c.Subbed)
		c.Subbed = ""
	}

	c.Children().Teardown()

	return nil
}
