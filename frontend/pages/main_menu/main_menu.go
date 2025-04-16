package main_menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/doodad"
	"github.com/jhuggett/frontend/pages/world_map"
	"github.com/jhuggett/sea/data/session"
	"github.com/jhuggett/sea/game_context"
)

type MainMenuPage struct {
	Width, Height int

	PageControls doodad.PageControls
	Gesturer     doodad.Gesturer

	Doodads []doodad.Doodad
}

func (m *MainMenuPage) Update() error {

	for _, widget := range m.Doodads {
		widget.Update()
	}

	m.Gesturer.Update()

	return nil
}

func (m *MainMenuPage) Draw(screen *ebiten.Image) {
	for _, widget := range m.Doodads {
		widget.Draw(screen)
	}
}

func (m *MainMenuPage) SetWidthAndHeight(width, height int) {
	m.Width, m.Height = width, height
}

func New(pageControls doodad.PageControls) (*MainMenuPage, error) {
	page := &MainMenuPage{
		PageControls: pageControls,
		Doodads:      []doodad.Doodad{},
		Gesturer:     doodad.NewGesturer(),
	}

	titleLabel := doodad.NewLabel()
	titleLabel.SetMessage("Ships Colonies Commerce: Main Menu")
	page.Doodads = append(page.Doodads, titleLabel)

	newGameButton := doodad.NewButton(
		"New Game",
		func() {
			newPage, err := world_map.New(nil)
			if err != nil {
				panic(err)
			}
			page.PageControls.Push(newPage)
		},
		page.Gesturer,
	)

	newGameButton.SetPosition(doodad.Below(titleLabel))

	page.Doodads = append(page.Doodads, newGameButton)

	sessions, err := session.All()
	if err != nil {
		return nil, err
	}

	for _, gameSession := range sessions {
		loadGameButton := doodad.NewButton(
			"Load Game",
			func() {
				newPage, err := world_map.New(&game_context.Snapshot{
					ShipID:    gameSession.ShipID,
					PlayerID:  gameSession.PlayerID,
					GameMapID: gameSession.GameMapID,
				})
				if err != nil {
					panic(err)
				}
				page.PageControls.Push(newPage)
			},
			page.Gesturer,
		)

		lastDoodad := page.Doodads[len(page.Doodads)-1]

		if d, ok := lastDoodad.(*doodad.Button); ok {
			loadGameButton.SetPosition(doodad.Below(d))
		}

		page.Doodads = append(page.Doodads, loadGameButton)
	}

	return page, nil
}
