package main

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/sea/db"

	"github.com/jhuggett/sea/log"

	"github.com/jhuggett/frontend/doodad"
	"github.com/jhuggett/frontend/pages/main_menu"
)

type Game struct {
	CurrentPage doodad.Page

	running bool
}

func (g *Game) Start() {
	if g.running {
		return
	}
	g.running = true

	ebiten.SetWindowSize(1200, 800)
	ebiten.SetWindowTitle("Ships Colonies Commerce")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	slog.SetDefault(
		slog.New(log.NewHandler(&log.HandlerOptions{
			HandlerOptions: slog.HandlerOptions{
				AddSource: true,
				Level:     log.OptInDebug,
			},
			UseColor:  true,
			BlockList: []string{},
			Allowlist: []string{},
		})),
	)

	db.Conn()
	db.Migrate()

	defer db.Close()

	game := &Game{}

	firstPage, err := main_menu.New(game)
	if err != nil {
		panic(err)
	}

	game.Push(firstPage)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func (g *Game) Push(page doodad.Page) {
	g.CurrentPage = page
}

func (g *Game) Pop() {
	// TODO: Implement
}

func (g *Game) Update() error {
	if err := g.CurrentPage.Update(); err != nil {
		slog.Error("Error updating page", "error", err)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.CurrentPage.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.CurrentPage.SetWidthAndHeight(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
