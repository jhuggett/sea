package main

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/sea/db"

	"github.com/jhuggett/sea/log"

	"github.com/jhuggett/sea/test/economy/page"
	"github.com/jhuggett/sea/test/economy/pages/main_menu"
)

type Game struct {
	CurrentPage page.Page
}

func (g *Game) Push(page page.Page) {
	g.CurrentPage = page
}

func (g *Game) Pop() {
	// TODO: Implement
}

func (g *Game) Update() error {

	// 	x, y := ebiten.CursorPosition()
	// 	if g.dragStartX == 0 && g.dragStartY == 0 {
	// 		g.dragStartX = x
	// 		g.dragStartY = y
	// 	} else {
	// 		g.Camera.X += float64(g.dragStartX-x) / g.Camera.Zoom
	// 		g.Camera.Y += float64(g.dragStartY-y) / g.Camera.Zoom
	// 		g.dragStartX = x
	// 		g.dragStartY = y
	// 	}
	// } else {
	// 	g.dragStartX = 0
	// 	g.dragStartY = 0
	//}

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

func run() {
	ebiten.SetWindowSize(2000, 800)
	ebiten.SetWindowTitle("Ships Colonies Commerce")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	runGame()
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

func runGame() {
	slog.SetDefault(
		slog.New(log.NewHandler(&log.HandlerOptions{
			HandlerOptions: slog.HandlerOptions{
				AddSource: true,
				Level:     log.OptInDebug,
			},
			UseColor: true,

			BlockList: []string{"backend/timeline", "backend/utils/callback"},
			Allowlist: []string{},
		})),
	)

	db.Conn()
	db.Migrate()
}
