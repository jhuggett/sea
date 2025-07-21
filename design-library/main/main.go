package main

import (
	design_library "design-library"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	game := design_library.NewApp(func(app *design_library.App) {
		ebiten.SetWindowSize(1200, 800)
		ebiten.SetWindowTitle("Design Library Example")
		ebiten.SetWindowDecorated(true)
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

		ebiten.SetCursorMode(ebiten.CursorModeVisible)

		firstPage := NewFirstPage(
			app,
		)
		app.Push(firstPage)

		if err := ebiten.RunGame(app); err != nil {
			panic(err)
		}
	})
	game.Start()
}
