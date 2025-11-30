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

	"github.com/jhuggett/frontend/pages/loading"
	"github.com/jhuggett/frontend/pages/settings"
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

	navigateToGame := func(snapshot *game_context.Snapshot) {
		loadingPage := loading.New(func(lpc *loading.LoadingPageController) {
			newPage := world_map.New(snapshot, m.App)

			newPage.SetLayout(box.Computed(func(b *box.Box) {
				b.Copy(m.App.Box)
			}))

			newPage.LoadRequiredData()

			m.App.Replace(newPage)
		})

		m.App.Push(loadingPage)
	}

	loadSaves := []doodad.Doodad{}

	for _, gameSession := range sessions {
		loadGameButton := button.New(button.Config{
			OnClick: func(b *button.Button) {
				navigateToGame(&game_context.Snapshot{
					ShipID:    gameSession.ShipID,
					PlayerID:  gameSession.PlayerID,
					GameMapID: gameSession.GameMapID,
				})
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
			navigateToGame(nil)

		},
		Config: label.Config{
			Message: "New Game",
			Layout:  box.Zeroed(),
		},
	})

	optionsButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			settingsPage := settings.New(m.App)
			m.App.Push(settingsPage)
		},
		Config: label.Config{
			Message: "Settings",
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

	panel := stack.New(stack.Config{

		SpaceBetween: 10,
	})

	m.AddChild(panel)

	panel.AddChild(
		titleLabel,
	)
	panel.AddChild(
		loadSaves...,
	)
	panel.AddChild(
		newGameButton,
		optionsButton,
		exitButton,
	)

	m.Children().Setup()

	panel.Box.Computed(func(b *box.Box) {
		b.CenterWithin(m.App.Box)
	})
}

func New(
	app *design_library.App,
) *MainMenuPage {
	page := &MainMenuPage{
		App: app,
	}

	return page
}
