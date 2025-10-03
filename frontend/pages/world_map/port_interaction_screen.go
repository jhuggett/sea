package world_map

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
)

type PortInteractionScreen struct {
	doodad.Default
	Manager *game.Manager
	Port    *game.Port
	OnClose func()
}

func NewPortInteractionScreen(manager *game.Manager, port *game.Port, onClose func()) *PortInteractionScreen {
	screen := &PortInteractionScreen{
		Manager: manager,
		Port:    port,
		OnClose: onClose,
	}

	return screen
}

func (p *PortInteractionScreen) Setup() {
	// Root panel to capture all input
	rootPanel := stack.New(stack.Config{
		Layout: box.Computed(func(b *box.Box) {
			b.Copy(p.Box)
		}),
	})
	p.AddChild(rootPanel)

	// Header with port name and close button
	portNameLabel := label.New(label.Config{
		Message:  p.Port.RawData.Name,
		FontSize: 30,
	})

	closeButton := button.New(button.Config{
		OnClick: func(*button.Button) {
			if p.OnClose != nil {
				p.OnClose()
			}
		},
		Config: label.Config{
			Message:         "✕",
			BackgroundColor: color.RGBA{120, 0, 0, 255},
			Padding: label.Padding{
				Top:    5,
				Bottom: 5,
				Left:   10,
				Right:  10,
			},
		},
	})

	header := stack.New(stack.Config{
		Children: doodad.NewChildren(
			p,
			[]doodad.Doodad{
				portNameLabel,
				closeButton,
			},
		),
		Type:         stack.Horizontal,
		SpaceBetween: 20,
		Padding: stack.Padding{
			Top:    10,
			Bottom: 10,
			Left:   20,
			Right:  20,
		},
		BackgroundColor: color.RGBA{40, 20, 0, 255},
	})

	// Create trade button
	tradeButton := button.New(button.Config{
		OnClick: func(*button.Button) {
			// TODO: Implement trade functionality
		},
		Config: label.Config{
			Message:         "Trade",
			BackgroundColor: color.RGBA{0, 80, 120, 255},
			Padding: label.Padding{
				Top:    10,
				Bottom: 10,
				Left:   20,
				Right:  20,
			},
		},
	})

	// Create repair button
	repairButton := button.New(button.Config{
		OnClick: func(*button.Button) {
			// TODO: Implement repair functionality
		},
		Config: label.Config{
			Message:         "Repair Ship",
			BackgroundColor: color.RGBA{0, 120, 80, 255},
			Padding: label.Padding{
				Top:    10,
				Bottom: 10,
				Left:   20,
				Right:  20,
			},
		},
	})

	// Create hire crew button
	hireCrewButton := button.New(button.Config{
		OnClick: func(*button.Button) {
			// TODO: Implement hire crew functionality
		},
		Config: label.Config{
			Message:         "Hire Crew",
			BackgroundColor: color.RGBA{120, 80, 0, 255},
			Padding: label.Padding{
				Top:    10,
				Bottom: 10,
				Left:   20,
				Right:  20,
			},
		},
	})

	// Port description
	descriptionLabel := label.New(label.Config{
		Message:  fmt.Sprintf("Port of %s is located at (%d, %d)", p.Port.RawData.Name, int(p.Port.Location().X), int(p.Port.Location().Y)),
		FontSize: 16,
	})

	// Create buttons container
	buttonContainer := stack.New(stack.Config{
		Children: doodad.NewChildren(
			p,
			[]doodad.Doodad{
				tradeButton,
				repairButton,
				hireCrewButton,
			},
		),
		Type:         stack.Horizontal,
		SpaceBetween: 20,
	})

	// Create main content container
	contentChildren := doodad.NewChildren(
		p,
		[]doodad.Doodad{
			descriptionLabel,
			buttonContainer,
		},
	)

	contentPanel := stack.New(stack.Config{
		Children:        contentChildren,
		Type:            stack.Vertical,
		BackgroundColor: colors.Panel,
		SpaceBetween:    20,
		Padding: stack.Padding{
			Top:    20,
			Bottom: 20,
			Left:   20,
			Right:  20,
		},
	})

	// Main layout container
	mainPanel := stack.New(stack.Config{
		Children: doodad.NewChildren(
			p,
			[]doodad.Doodad{
				header,
				contentPanel,
			},
		),
		Type:            stack.Vertical,
		BackgroundColor: color.RGBA{40, 20, 0, 255},
	})

	p.AddChild(mainPanel)
	p.Children().Setup()

	mainPanel.Box.Computed(func(b *box.Box) {
		b.Copy(p.Box)
	})
}

func (p *PortInteractionScreen) Draw(screen *ebiten.Image) {
	if !p.IsVisible() {
		return
	}

	p.Children().Draw(screen)
}

func (p *PortInteractionScreen) Update() error {
	return nil
}
