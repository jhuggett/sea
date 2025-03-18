package world_map

import (
	"log/slog"
	"math"
	"time"

	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
)

func (w *WorldMapPage) Update() error {
	w.ui.Update()

	if input.UIHovered {
		return nil
	}

	x, y := ebiten.CursorPosition()

	w.MouseLocationX = x
	w.MouseLocationY = y

	// moveAmount := w.TileSize

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		// w.Camera.Y -= moveAmount
		w.Camera.Position[1] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		// w.Camera.Y += moveAmount
		w.Camera.Position[1] += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		// w.Camera.X -= moveAmount
		w.Camera.Position[0] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		// w.Camera.X += moveAmount
		w.Camera.Position[0] += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// center on ship
		w.CanvasTranslateTo(w.Ship.X, w.Ship.Y)
	}

	_, yoff := ebiten.Wheel()
	if yoff > 0 {
		w.Camera.ZoomFactor += 5
	}
	if yoff < 0 {
		w.Camera.ZoomFactor -= 5
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

		if w.Press == nil {
			w.Press = &Press{
				StartX:    x,
				StartY:    y,
				X:         x,
				Y:         y,
				TimeStart: time.Now(),
			}
		}

		if time.Since(w.Press.TimeStart) > 100*time.Millisecond || (math.Abs(float64(w.Press.StartX-x)) > 25 || math.Abs(float64(w.Press.StartY-y)) > 25) {
			w.Camera.Position[0] += float64(w.Press.X-x) / math.Pow(1.01, float64(w.Camera.ZoomFactor))
			w.Camera.Position[1] += float64(w.Press.Y-y) / math.Pow(1.01, float64(w.Camera.ZoomFactor))
		}

		w.Press.X = x
		w.Press.Y = y

	} else {
		if w.Press != nil {
			if time.Since(w.Press.TimeStart) < 100*time.Millisecond || (math.Abs(float64(w.Press.StartX-w.Press.X)) < 8 && math.Abs(float64(w.Press.StartY-w.Press.Y)) < 8) {
				slog.Info("Click", "x", w.Press.X, "y", w.Press.Y)
				if w.OnTileClicked != nil {
					worldX, worldY := w.Camera.ScreenToWorld(ebiten.CursorPosition())
					translatedX, translatedY := w.CanvasTranslateFrom(float64(worldX), float64(worldY))
					w.OnTileClicked(
						int(translatedX/float64(w.TileSize)),
						int(translatedY/float64(w.TileSize)),
					)
				}
			}

			w.Press = nil
		}
	}
	return nil
}
