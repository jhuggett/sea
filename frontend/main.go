package main

import (
	"log/slog"
	"os"
	"runtime/pprof"

	design_library "design-library"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/pages/main_menu"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/log"

	_ "net/http/pprof"
)

func main() {
	// game := &Game{}
	// game.Start()

	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	// --- OPTIONAL: HEAP PROFILER ---
	// Will write the heap profile when the game exits.
	defer func() {
		heapFile, err := os.Create("heap.prof")
		if err != nil {
			panic(err)
		}
		pprof.WriteHeapProfile(heapFile)
	}()

	// Switch to WAV format for sound files or OGG

	// sound.Setup()

	// backgroundMusic := sound.NewBackgroundMusic(
	// 	[]string{
	// 		"../assets/audio/IMG_0948.mp3",
	// 		"../assets/audio/IMG_0960.mp3",
	// 		"../assets/audio/IMG_0961.mp3",
	// 		"../assets/audio/IMG_0962.mp3",
	// 	},
	// )

	// go backgroundMusic.Run(context.Background())

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
