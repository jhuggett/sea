package design_library

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewApp(startup func(*App)) *App {
	app := &App{
		startup: startup,
	}

	if startup != nil {
		startup(app)
	}

	return app
}

type App struct {
	CurrentPage Page

	startup func(*App)
}

func (g *App) Start() {
	g.startup(g)
}

func (g *App) Push(page Page) {
	g.CurrentPage = page
}

func (g *App) Pop() {
	// TODO: Implement
}

func (g *App) Update() error {
	if err := g.CurrentPage.Update(); err != nil {
		slog.Error("Error updating page", "error", err)
	}

	return nil
}

func (g *App) Draw(screen *ebiten.Image) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from panic in Draw", "error", r)
		}
	}()

	g.CurrentPage.Draw(screen)
}

func (g *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.CurrentPage.SetWidthAndHeight(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
