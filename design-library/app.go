package design_library

import (
	"design-library/doodad"
	"design-library/position/box"
	"design-library/reaction"
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewApp(startup func(*App)) *App {
	app := &App{
		startup: startup,
		Default: *doodad.NewDefault(nil),
	}

	app.Children().Parent = &app.Default
	app.SetGesturer(reaction.NewGesturer())
	app.SetLayout(box.Zeroed())

	slog.Info("App initialized", "app", app, "layout", fmt.Sprintf("%p", app.Default.Layout()))

	app.SetReactions(&reaction.Reactions{})

	app.Reactions().Add(
		reaction.NewKeyDownReaction(
			reaction.SpecificKeyDown(ebiten.KeyD),
			func(event reaction.KeyDownEvent) {
				app.Children().PrettyPrint(0)
				fmt.Printf("App layout: %s\n", app.Default.Layout().String())
			},
		),
		reaction.NewKeyDownReaction(
			reaction.SpecificKeyDown(ebiten.KeyR),
			func(event reaction.KeyDownEvent) {
				app.Default.Layout().Recalculate()
				slog.Info("Recalculated layout")
			},
		),
	)
	app.Reactions().Register(app.Gesturer())

	// app.Gesturer().OnKeyDown(func(key ebiten.Key) error {
	// 	if key == ebiten.KeyD {
	// 		app.Children().PrettyPrint(0)
	// 		fmt.Printf("App layout: %s\n", app.Default.Layout().String())
	// 	}

	// 	if key == ebiten.KeyR {
	// 		app.Default.Layout().Recalculate()
	// 		slog.Info("Recalculated layout")
	// 	}

	// 	return nil
	// })

	return app
}

type App struct {
	CurrentPage doodad.Doodad

	startup func(*App)

	doodad.Default
}

func (g *App) Start() {
	g.startup(g)
}

func (g *App) Push(page doodad.Doodad) {
	// g.CurrentPage = page
	g.AddChild(page)
	g.Children().Setup()
}

func (g *App) Replace(page doodad.Doodad) {
	// g.CurrentPage = page
	g.Children().Clear()
	g.AddChild(page)

	g.Children().Setup()

	g.Box.Recalculate()
}

func (g *App) Pop() {
	// TODO: Implement
}

func (g *App) Update() error {
	// if err := g.CurrentPage.Update(); err != nil {
	// 	slog.Error("Error updating page", "error", err)
	// }

	g.Gesturer().Update()

	return nil
}

func (g *App) Draw(screen *ebiten.Image) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from panic in Draw", "error", r)
		}
	}()

	// g.CurrentPage.Draw(screen)

	g.Children().Draw(screen)
}

func (g *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// g.CurrentPage.SetWidthAndHeight(outsideWidth, outsideHeight)

	if g.Default.Layout().Width() != outsideWidth || g.Default.Layout().Height() != outsideHeight {
		g.Default.Layout().SetDimensions(outsideWidth, outsideHeight)
		g.Default.Layout().Recalculate()
	}

	return outsideWidth, outsideHeight
}
