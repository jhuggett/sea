package main

import (
	"bytes"
	"log"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

type Camera struct {
	X, Y float64
	Zoom float64
}

var tiles = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 1, 1, 1, 1, 0, 0, 0, 0},
	{0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 1, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 1, 0, 1, 1, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

type Game struct {
	ui *ebitenui.UI

	Camera Camera

	BackgroundImage *ebiten.Image

	ExampleImage *ebiten.Image
}

func (g *Game) Update() error {
	g.ui.Update()

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Camera.Y -= 4 / (1 + g.Camera.Zoom)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Camera.Y += 4 / (1 + g.Camera.Zoom)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Camera.X -= 4 / (1 + g.Camera.Zoom)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Camera.X += 4 / (1 + g.Camera.Zoom)
	}

	_, yoff := ebiten.Wheel()
	if yoff > 0 {
		slog.Info("Zooming in", "zoom", g.Camera.Zoom)
		g.Camera.Zoom -= .25
	}
	if yoff < 0 {
		slog.Info("Zooming out", "zoom", g.Camera.Zoom)
		g.Camera.Zoom += .25
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y, row := range tiles {
		for x, tile := range row {
			if tile == 1 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*8), float64(y*8))
				op.GeoM.Translate(-g.Camera.X, -g.Camera.Y)
				op.GeoM.Scale(1+g.Camera.Zoom, 1+g.Camera.Zoom)
				screen.DrawImage(g.ExampleImage, op)
			}
		}
	}

	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func run() {
	ebiten.SetWindowSize(2000, 800)
	ebiten.SetWindowTitle("Ships Colonies Commerce")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

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

	// popSize, graph := do()

	_, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}

	// fontFace := &text.GoTextFace{
	// 	Source: s,
	// 	Size:   32,
	// }
	// This creates a text widget that says "Hello World!"
	// helloWorldLabel := widget.NewText(
	// 	widget.TextOpts.Text("Hello World!", fontFace, color.White),
	// 	widget.TextOpts.MaxWidth(20),
	// )

	// (&widgets.Button{
	// 	Text: fmt.Sprintf("Population: %d", popSize),
	// }).Setup(rootContainer)

	// (&widgets.TextArea{
	// 	Contents: graph,
	// }).Setup(rootContainer)

	// To display the text widget, we have to add it to the root container.
	//rootContainer.AddChild(helloWorldLabel)

	game := &Game{
		ui: eui,
	}

	exampleImage, _, err := ebitenutil.NewImageFromFile("./assets/images/tile_0006.png")
	if err != nil {
		log.Fatal(err)
	}
	game.ExampleImage = exampleImage

	backendImage, _, err := ebitenutil.NewImageFromFile("./assets/images/parchmentBasic.png")
	if err != nil {
		log.Fatal(err)
	}
	game.BackgroundImage = backendImage

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
