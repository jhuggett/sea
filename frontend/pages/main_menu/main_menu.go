package main_menu

import (
	design_library "design-library"
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"fmt"
	"os"

	"github.com/jhuggett/frontend/pages/world_map"
	"github.com/jhuggett/sea/data/session"
	"github.com/jhuggett/sea/game_context"
)

type MainMenuPage struct {
	App *design_library.App
	doodad.Default
}

func (m *MainMenuPage) Setup() {
	titleLabel := label.New(label.Config{
		Message:  "Ships Colonies & Commerce",
		FontSize: 36,
		Padding: label.Padding{
			Bottom: 20,
		},
	})

	// This should really go into some sort of controller
	sessions, err := session.All()
	if err != nil {
		panic(fmt.Errorf("failed to retrieve game sessions: %w", err))
	}

	loadSaves := []doodad.Doodad{}

	for _, gameSession := range sessions {
		loadGameButton := button.New(button.Config{
			OnClick: func(b *button.Button) {
				newPage := world_map.New(&game_context.Snapshot{
					ShipID:    gameSession.ShipID,
					PlayerID:  gameSession.PlayerID,
					GameMapID: gameSession.GameMapID,
				}, m.App)
				// m.PageControls.Push(newPage)
				m.App.Replace(newPage)
			},
			Config: label.Config{
				Message: fmt.Sprintf("Load Game: %s", gameSession.UpdatedAt.Format("2006-01-02 15:04:05")),
				Layout:  box.Zeroed(),
			},
		})
		loadSaves = append(loadSaves, loadGameButton)
	}

	newGameButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			newPage := world_map.New(nil, m.App)
			// m.PageControls.Push(newPage)
			m.App.Replace(newPage)
		},
		Config: label.Config{
			Message: "New Game",
			Layout:  box.Zeroed(),
		},
	})

	optionsButton := button.New(button.Config{
		OnClick: func(b *button.Button) {},
		Config: label.Config{
			Message: "Options",
			Layout:  box.Zeroed(),
		},
	})

	exitButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			os.Exit(0) // Exit the application
		},
		Config: label.Config{
			Message: "Exit",
			Layout:  box.Zeroed(),
		},
	})

	panelChildren := doodad.NewChildren(
		m,
		[]doodad.Doodad{
			titleLabel,
		},
		loadSaves,
		[]doodad.Doodad{
			newGameButton,
			optionsButton,
			exitButton,
		},
	)

	panel := stack.New(stack.Config{
		Children: panelChildren,
		Type:     stack.Vertical,
		Layout: box.Computed(func(b *box.Box) {
			boundingBox := box.Bounding(panelChildren.Boxes())
			b.SetDimensions(boundingBox.Width(), boundingBox.Height()).CenterWithin(m.Box)
		}),
		SpaceBetween: 10,
	})

	m.AddChild(panel)

	m.Children().Setup()
}

func New(
	app *design_library.App,
) *MainMenuPage {
	page := &MainMenuPage{
		App: app,
	}

	// page.Setup()

	// titleLabel, err := label.New(label.Config{
	// 	Message: "Ships Colonies & Commerce",
	// 	Layout: position.NewBox(position.BoxConfig{}).Computed(func(b *position.Box) *position.Box {
	// 		return b
	// 	}),
	// 	FontSize: 48,
	// 	Padding: label.Padding{
	// 		Top:    20,
	// 		Bottom: 20,
	// 		Left:   20,
	// 		Right:  20,
	// 	},
	// 	ForegroundColor: color.RGBA{
	// 		A: 255,
	// 		R: 255,
	// 		G: 155,
	// 		B: 105,
	// 	},
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create title label: %w", err)
	// }
	// page.Children.Add(titleLabel)

	// newGameButton, err := button.New(button.Config{
	// 	OnClick: func() {
	// 		// newPage, err := world_map.New(nil)
	// 		// if err != nil {
	// 		// 	panic(err)
	// 		// }
	// 		// page.PageControls.Push(newPage)
	// 	},
	// 	Gesturer: page.Gesturer,
	// 	Config: label.Config{
	// 		Message: "New Game",
	// 		Layout: position.NewBox(position.BoxConfig{}).Computed(func(b *position.Box) *position.Box {
	// 			return b.MoveAbove(titleLabel.Box)
	// 		}),
	// 	},
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create new game button: %w", err)
	// }
	// page.Children.Add(newGameButton)

	// sessions, err := session.All()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to retrieve game sessions: %w", err)
	// }

	// var previousButton *button.Button

	// for _, gameSession := range sessions {

	// 	previousButtonReferenceCopy := previousButton

	// 	loadGameButton, err := button.New(button.Config{
	// 		Config: label.Config{
	// 			Message: fmt.Sprintf("Load Game: %s", gameSession.UpdatedAt.Format("2006-01-02 15:04:05")),
	// 			Layout: position.NewBox(position.BoxConfig{}).Computed(func(b *position.Box) *position.Box {
	// 				if previousButtonReferenceCopy == nil {
	// 					return b.MoveBelow(titleLabel.Box)
	// 				}
	// 				return b.MoveBelow(previousButtonReferenceCopy.Box)
	// 			}),
	// 			Padding: label.Padding{
	// 				Left: 20,
	// 			},
	// 		},
	// 		OnClick: func() {
	// 			// newPage, err := world_map.New(&game_context.Snapshot{
	// 			// 	ShipID:    gameSession.ShipID,
	// 			// 	PlayerID:  gameSession.PlayerID,
	// 			// 	GameMapID: gameSession.GameMapID,
	// 			// })
	// 			// if err != nil {
	// 			// 	panic(err)
	// 			// }
	// 			// page.PageControls.Push(newPage)
	// 		},
	// 		Gesturer: page.Gesturer,
	// 	})
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to create load game button: %w", err)
	// 	}

	// 	page.Children.Add(loadGameButton)
	// 	previousButton = loadGameButton
	// }

	return page
}
