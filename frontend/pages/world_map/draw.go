package world_map

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"
)

type Camera struct {
	ViewPort   f64.Vec2
	Position   f64.Vec2
	ZoomFactor float64
	Rotation   float64

	TileSize float64

	CameraIsMoving bool
	LastPosition   f64.Vec2
}

func (c *Camera) String() string {
	return fmt.Sprintf("Camera{ViewPort: %v, Position: %v, ZoomFactor: %d, Rotation: %d}", c.ViewPort, c.Position, c.ZoomFactor, c.Rotation)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] / 2,
		c.ViewPort[1] / 2,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])

	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])

	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	}
	screen.DrawImage(
		world,
		op,
	)
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) Reset() {
	c.Position[0] = 0
	c.Position[1] = 0
}

func (w *WorldMapPage) CanvasTranslateTo(x, y float64) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}

	// op.GeoM.Translate(x+float64(-w.SmallestPointX)*w.TileSize, y+float64(-w.SmallestPointY)*w.TileSize)

	return op
}

func (w *WorldMapPage) CanvasTranslateFrom(x, y float64) (float64, float64) {
	// return x - float64(-w.SmallestPointX)*w.TileSize, y - float64(-w.SmallestPointY)*w.TileSize

	return 0, 0
}

func (w *WorldMapPage) Draw(screen *ebiten.Image) {
	// w.Canvas.DrawImage(
	// 	w.MapImage,
	// 	w.CanvasTranslateTo(float64(w.SmallestPointX)*w.TileSize, float64(w.SmallestPointY)*w.TileSize),
	// )

	// for _, continent := range w.Continents {
	// 	w.Canvas.DrawImage(
	// 		continent.Image,
	// 		w.CanvasTranslateTo(float64(continent.OriginX)*w.TileSize, float64(continent.OriginY)*w.TileSize),
	// 	)
	// }

	// for _, port := range w.Ports {
	// 	w.Canvas.DrawImage(
	// 		port.Image,
	// 		w.CanvasTranslateTo(float64(port.X)*w.TileSize, float64(port.Y)*w.TileSize),
	// 	)
	// }

	// if w.PlottedRoute.Image != nil {
	// 	w.Canvas.DrawImage(
	// 		w.PlottedRoute.Image,
	// 		w.CanvasTranslateTo(float64(w.PlottedRoute.OriginX)*w.TileSize, float64(w.PlottedRoute.OriginY)*w.TileSize),
	// 	)

	// }

	// worldX, worldY := w.Camera.ScreenToWorld(ebiten.CursorPosition())
	// translatedX, translatedY := w.CanvasTranslateFrom(float64(worldX), float64(worldY))

	// if w.CursorImage != nil {
	// 	w.Canvas.DrawImage(
	// 		w.CursorImage,
	// 		w.CanvasTranslateTo(math.Floor(translatedX/w.TileSize)*w.TileSize, math.Floor(translatedY/w.TileSize)*w.TileSize),
	// 	)
	// }

	// w.Canvas.DrawImage(
	// 	w.Ship.Image,
	// 	w.CanvasTranslateTo(w.Ship.X*w.TileSize, w.Ship.Y*w.TileSize),
	// )

	// w.Camera.Render(w.Canvas, screen)

	// w.ui.Draw(screen)

	for _, doodad := range w.Doodads {
		doodad.Draw(screen)
	}

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()),
	)

	worldX, worldY := w.SpaceTranslator.FromScreenToWorld(float64(w.MouseLocationX), float64(w.MouseLocationY))

	dataX, dataY := w.SpaceTranslator.FromWorldToData(
		w.SpaceTranslator.FromScreenToWorld(float64(w.MouseLocationX), float64(w.MouseLocationY)),
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Screen: %d,%d | World: %d,%d | Data: %d,%d | Zoom: %.2f",
			int(w.MouseLocationX),
			int(w.MouseLocationY),

			int(worldX),
			int(worldY),

			int(dataX),
			int(dataY),

			w.Camera.ZoomFactor,
		),
		0, w.Height-32,
	)
}
