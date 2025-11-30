package world_map

import (
	"design-library/doodad"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/frontend/pages/world_map/camera"
)

func NewPlayerShipDoodad(gameManager *game.Manager, Camera *camera.Camera) *ShipDoodad {
	shipDoodad := &ShipDoodad{
		GameManager: gameManager,
		Camera:      Camera,
	}

	return shipDoodad
}

type ShipDoodad struct {
	GameManager *game.Manager
	Ship        *game.Ship

	img *ebiten.Image

	Origin func() (int, int)
	Scale  func() (float64, float64)

	Camera *camera.Camera

	doodad.Default
}

func (w *ShipDoodad) Draw(screen *ebiten.Image) {

	originX, originY := w.Origin()
	scaleX, scaleY := w.Scale()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64((w.Ship.Location().X)*w.Ship.Manager.TileSize()-originX),
		float64((w.Ship.Location().Y)*w.Ship.Manager.TileSize()-originY),
	)

	op.GeoM.Scale(
		scaleX,
		scaleY,
	)
	screen.DrawImage(w.img, op)
}

func (w *ShipDoodad) Setup() {
	w.Origin = func() (int, int) {
		return int(w.Camera.Position[0]), int(w.Camera.Position[1])
	}
	w.Scale = func() (float64, float64) {
		return w.Camera.ZoomFactor, w.Camera.ZoomFactor
	}

	w.Ship = w.GameManager.PlayerShip

	tileSize := w.Ship.Manager.TileSize()
	w.img = ebiten.NewImage(tileSize, tileSize)

	// Draw a blue square to represent the ship
	shipColor := color.RGBA{
		R: 0,
		G: 0,
		B: 255,
		A: 255, // Blue color
	}

	w.img.Fill(shipColor)
}
