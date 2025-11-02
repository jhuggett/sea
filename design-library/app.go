package design_library

import (
	"design-library/doodad"
	"design-library/position/box"
	"design-library/reaction"
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
			func(event *reaction.KeyDownEvent) {
				app.Children().PrettyPrint(0)
				fmt.Printf("App layout: %s\n", app.Default.Layout().String())
			},
		),
		reaction.NewKeyDownReaction(
			reaction.SpecificKeyDown(ebiten.KeyR),
			func(event *reaction.KeyDownEvent) {
				app.Default.Layout().Recalculate()
				slog.Info("Recalculated layout")
			},
		),
	)
	app.Reactions().Register(app.Gesturer(), app.Z())

	return app
}

type App struct {
	PageStack []doodad.Doodad

	startup func(*App)

	doodad.Default
}

func (g *App) Start() {
	g.startup(g)
}

func (g *App) Current() doodad.Doodad {
	if len(g.PageStack) == 0 {
		return nil
	}
	return g.PageStack[len(g.PageStack)-1]
}

func (g *App) SetupCurrentPage() {
	if g.Current() != nil {
		g.Current().Setup()
	}
}

func (g *App) Push(page doodad.Doodad) {
	if len(g.PageStack) > 0 {
		g.Current().Hide()
	}
	g.PageStack = append(g.PageStack, page)

	g.AddChild(page)
	g.SetupCurrentPage()
}

func (g *App) Replace(page doodad.Doodad) {
	if len(g.PageStack) > 0 {
		g.Children().Remove(g.Current())
		g.PageStack = g.PageStack[:len(g.PageStack)-1]
	}

	g.PageStack = append(g.PageStack, page)

	g.AddChild(page)

	g.SetupCurrentPage()

	g.Box.Recalculate()
}

func (g *App) Pop() {
	if len(g.PageStack) == 0 {
		return
	}

	g.Children().Remove(g.Current())
	g.PageStack = g.PageStack[:len(g.PageStack)-1]

	if len(g.PageStack) > 0 {
		g.Current().Show()
	}

	g.Box.Recalculate()
}

func (g *App) PopBy(count int) {
	if len(g.PageStack) == 0 || count <= 0 {
		return
	}

	if count > len(g.PageStack) {
		count = len(g.PageStack)
	}

	for i := 0; i < count; i++ {
		g.Children().Remove(g.Current())
		g.PageStack = g.PageStack[:len(g.PageStack)-1]
	}

	if len(g.PageStack) > 0 {
		g.Current().Show()
	}

	g.Box.Recalculate()
}

func (g *App) PopToRoot() {
	if len(g.PageStack) == 0 {
		return
	}

	g.PopBy(len(g.PageStack) - 1)
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

	x, y := g.Gesturer().CurrentMouseLocation()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse: %d, %d", x, y), 4, screen.Bounds().Dy()-14)

}

func (g *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// g.CurrentPage.SetWidthAndHeight(outsideWidth, outsideHeight)

	if g.Default.Layout().Width() != outsideWidth || g.Default.Layout().Height() != outsideHeight {
		g.Default.Layout().SetDimensions(outsideWidth, outsideHeight)
		g.Default.Layout().Recalculate()
	}

	return outsideWidth, outsideHeight
}
