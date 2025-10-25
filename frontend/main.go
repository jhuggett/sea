package main

import (
	"context"
	"log/slog"

	design_library "design-library"
	"design-library/sound"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/pages/main_menu"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/log"
)

func main() {
	// game := &Game{}
	// game.Start()

	sound.Setup()

	backgroundMusic := sound.NewBackgroundMusic(
		[]string{
			"../assets/audio/IMG_0948.mp3",
			"../assets/audio/IMG_0960.mp3",
			"../assets/audio/IMG_0961.mp3",
			"../assets/audio/IMG_0962.mp3",
		},
	)

	go backgroundMusic.Run(context.Background())

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
