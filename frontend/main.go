package main

import (
	"design-library/sound"
	"log/slog"

	design_library "design-library"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/pages/main_menu"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/log"
)

func main() {
	// game := &Game{}
	// game.Start()

	go func() {
		sound.Setup()
		_, err := sound.Play("../assets/audio/IMG_0960.mp3")
		if err != nil {
			panic(err)
		}
	}()

	app := design_library.NewApp(func(app *design_library.App) {
		ebiten.SetWindowSize(1200, 800)
		ebiten.SetWindowTitle("Ships Colonies Commerce")
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
		ebiten.SetScreenClearedEveryFrame(true)

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

		firstPage := main_menu.New(app)

		app.Push(firstPage)

		if err := ebiten.RunGame(app); err != nil {
			panic(err)
		}
	})

	app.Start()

	slog.Info("Game over.")
}
