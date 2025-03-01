package main_menu

import (
	"bytes"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/jhuggett/sea/data/session"
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/test/economy/page"
	"github.com/jhuggett/sea/test/economy/pages/world_map"
	"github.com/jhuggett/sea/test/economy/widgets"
	"golang.org/x/image/font/gofont/goregular"
)

type MainMenuPage struct {
	ui *ebitenui.UI

	Width, Height int

	PageControls page.PageControls
}

func (m *MainMenuPage) Update() error {
	m.ui.Update()
	return nil
}

func (m *MainMenuPage) Draw(screen *ebiten.Image) {
	m.ui.Draw(screen)
}

func (m *MainMenuPage) SetWidthAndHeight(width, height int) {
	m.Width, m.Height = width, height
}

func New(pageControls page.PageControls) (*MainMenuPage, error) {
	// This creates the root container for this UI.
	// All other UI elements must be added to this container.
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			// widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	// This adds the root container to the UI, so that it will be rendered.
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	page := &MainMenuPage{
		ui: eui,

		PageControls: pageControls,
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}

	fontFace := &text.GoTextFace{
		Source: s,
		Size:   32,
	}

	helloWorldLabel := widget.NewText(
		widget.TextOpts.Text("Ships Colonies Commerce", fontFace, color.White),
	)
	rootContainer.AddChild(helloWorldLabel)

	(&widgets.Button{
		Text: "New Game",
		OnClick: func() {
			newPage, err := world_map.New(nil)
			if err != nil {
				panic(err)
			}
			page.PageControls.Push(newPage)
		},
	}).Setup(rootContainer)

	sessions, err := session.All()
	if err != nil {
		return nil, err
	}

	for _, gameSession := range sessions {
		(&widgets.Button{
			Text: "Load Game",
			OnClick: func() {
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
		}).Setup(rootContainer)
	}

	// (&widgets.TextArea{
	// 	Contents: "Hello",
	// }).Setup(rootContainer)

	return page, nil
}
