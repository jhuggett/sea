package world_map

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
)

type DockConfirmation struct {
	doodad.Default
	Manager       *game.Manager
	Port          *game.Port
	OnDockAccept  func()
	OnDockDecline func()
}

func NewDockConfirmation(manager *game.Manager, port *game.Port, onDockAccept func(), onDockDecline func()) *DockConfirmation {
	confirmation := &DockConfirmation{
		Manager:       manager,
		Port:          port,
		OnDockAccept:  onDockAccept,
		OnDockDecline: onDockDecline,
	}

	return confirmation
}

func (d *DockConfirmation) Setup() {
	// Root panel catches all mouse input to prevent clicks from passing through
	rootPanel := stack.New(stack.Config{
		Layout: box.Computed(func(b *box.Box) {
			b.Copy(d.Box)
		}),
	})
	d.AddChild(rootPanel)

	if d.Port == nil {
		return
	}

	panel := stack.New(stack.Config{
		Type:            stack.Vertical,
		BackgroundColor: colors.Panel,
		SpaceBetween:    20,
		Padding: stack.Padding{
			Top:    20,
			Bottom: 20,
			Left:   40,
			Right:  40,
		},
		FitContents: true,
	})

	d.AddChild(panel)

	// Create the title and message
	titleLabel := label.New(label.Config{
		Message:  "Dock at " + d.Port.RawData.Name + "?",
		FontSize: 24,
	})
	panel.AddChild(titleLabel)

	messageLabel := label.New(label.Config{
		Message:  "Would you like to dock your ship at this port?",
		FontSize: 18,
	})
	panel.AddChild(messageLabel)

	// Create button container
	buttonContainer := stack.New(stack.Config{
		Type:         stack.Horizontal,
		SpaceBetween: 20,
		FitContents:  true,
	})
	panel.AddChild(buttonContainer)

	// Create buttons
	acceptButton := button.New(button.Config{
		OnClick: func(*button.Button) {
			if d.OnDockAccept != nil {
				d.OnDockAccept()
			}
		},
		Config: label.Config{
			Message:         "Dock",
			BackgroundColor: color.RGBA{0, 120, 0, 255},
			Padding: label.Padding{
				Top:    10,
				Bottom: 10,
				Left:   20,
				Right:  20,
			},
		},
	})
	buttonContainer.AddChild(acceptButton)

	declineButton := button.New(button.Config{
		OnClick: func(*button.Button) {
			if d.OnDockDecline != nil {
				d.OnDockDecline()
			}
		},
		Config: label.Config{
			Message:         "Cancel",
			BackgroundColor: color.RGBA{120, 0, 0, 255},
			Padding: label.Padding{
				Top:    10,
				Bottom: 10,
				Left:   20,
				Right:  20,
			},
		},
	})
	buttonContainer.AddChild(declineButton)

	d.Children().Setup()

	panel.Box.Computed(func(b *box.Box) {
		b.CenterWithin(rootPanel.Layout())
	})

	d.Layout().Recalculate()
}

func (d *DockConfirmation) Draw(screen *ebiten.Image) {
	if !d.IsVisible() {
		return
	}

	// Apply semi-transparent overlay to the background
	background := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	background.DrawImage(screen, nil)
	options := &ebiten.DrawImageOptions{}
	options.ColorM.Scale(0.3, 0.3, 0.3, 1) // Darken the background
	screen.DrawImage(background, options)

	// Draw the dialog
	d.Children().Draw(screen)
}

func (d *DockConfirmation) Update() error {
	return nil
}
