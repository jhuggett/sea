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
	ZoomFactor int
	Rotation   int
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
	screen.DrawImage(
		world, &ebiten.DrawImageOptions{
			GeoM: c.worldMatrix(),
		},
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

}

// func (w *WorldMapPage) TranslatedTileOrigin(x, y int) (float64, float64) {
// 	tx := (float64(x)*(w.TileSize) - w.Camera.X)
// 	ty := (float64(y)*(w.TileSize) - w.Camera.Y)

// 	return tx, ty
// }

// func (w *WorldMapPage) ScaledTranslatedTileOrigin(x, y int) (float64, float64) {
// 	tx := (float64(x)*(w.TileSize/w.Camera.Zoom) - w.Camera.X)
// 	ty := (float64(y)*(w.TileSize/w.Camera.Zoom) - w.Camera.Y)

// 	return tx, ty
// }

// func (w *WorldMapPage) MouseLocationAsTileLocation() (int, int) {
// 	x := int((float64(w.MouseLocationX) + w.Camera.X) / (w.TileSize * w.Camera.Zoom))
// 	y := int((float64(w.MouseLocationY) + w.Camera.Y) / (w.TileSize * w.Camera.Zoom))
// 	return x, y
// }

func (w *WorldMapPage) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(float64(w.SmallestPointX)*w.TileSize, float64(w.SmallestPointY)*w.TileSize)
	// op.GeoM.Translate(-w.Camera.X, -w.Camera.Y)
	op.GeoM.Translate(float64(w.SmallestPointX)*w.TileSize, float64(w.SmallestPointY)*w.TileSize)
	w.Canvas.DrawImage(w.MapImage, op)

	for _, continent := range w.Continents {
		op := &ebiten.DrawImageOptions{}

		// op.GeoM.Translate(float64(continent.OriginX)*w.TileSize, float64(continent.OriginY)*w.TileSize)
		// op.GeoM.Translate(-w.Camera.X, -w.Camera.Y)
		// op.GeoM.Translate(w.TranslatedTileOrigin(continent.OriginX, continent.OriginY))
		op.GeoM.Translate(float64(continent.OriginX)*w.TileSize, float64(continent.OriginY)*w.TileSize)
		// op.GeoM.Scale(w.Camera.Zoom, w.Camera.Zoom)

		w.Canvas.DrawImage(continent.Image, op)
	}

	for _, port := range w.Ports {
		op := &ebiten.DrawImageOptions{}

		// op.GeoM.Translate(float64(port.X)*w.TileSize, float64(port.Y)*w.TileSize)
		// op.GeoM.Translate(-w.Camera.X, -w.Camera.Y)
		op.GeoM.Translate(float64(port.X)*w.TileSize, float64(port.Y)*w.TileSize)
		// op.GeoM.Scale(w.Camera.Zoom, w.Camera.Zoom)

		w.Canvas.DrawImage(port.Image, op)
	}

	if w.PlottedRoute.Image != nil {
		op := &ebiten.DrawImageOptions{}
		// op.GeoM.Translate(-w.Camera.X, -w.Camera.Y)
		// op.GeoM.Translate(float64(w.PlottedRoute.OriginX)*w.TileSize, float64(w.PlottedRoute.OriginY)*w.TileSize)
		op.GeoM.Translate(float64(w.PlottedRoute.OriginX)*w.TileSize, float64(w.PlottedRoute.OriginY)*w.TileSize)
		// op.GeoM.Scale(w.Camera.Zoom, w.Camera.Zoom)
		w.Canvas.DrawImage(w.PlottedRoute.Image, op)
	}

	worldX, worldY := w.Camera.ScreenToWorld(ebiten.CursorPosition())

	if w.CursorImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(math.Floor(worldX/w.TileSize)*w.TileSize, math.Floor(worldY/w.TileSize)*w.TileSize)
		w.Canvas.DrawImage(w.CursorImage, op)
	}

	op = &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(w.Ship.X*w.TileSize, w.Ship.Y*w.TileSize)
	// op.GeoM.Translate(-w.Camera.X, -w.Camera.Y)
	// op.GeoM.Translate(w.TranslatedTileOrigin(int(w.Ship.X), int(w.Ship.Y)))
	// op.GeoM.Scale(w.Camera.Zoom, w.Camera.Zoom)
	op.GeoM.Translate(w.Ship.X*w.TileSize, w.Ship.Y*w.TileSize)
	w.Canvas.DrawImage(w.Ship.Image, op)

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse: %d, %d", w.MouseLocationX, w.MouseLocationY), 0, 30)
	// x, y := w.MouseLocationAsTileLocation()
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse tile: %d, %d",
	// 	x, y,
	// ), 0, 40)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Camera: %0.2f, %0.2f", w.Camera.X, w.Camera.Y), 0, 50)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Zoom: %0.2f", w.Camera.Zoom), 0, 60)

	w.Camera.Render(w.Canvas, screen)

	w.ui.Draw(screen)

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()),
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("%s\nCursor World Pos: %.2f,%.2f",
			w.Camera.String(),
			worldX, worldY),
		0, w.Height-32,
	)
}
