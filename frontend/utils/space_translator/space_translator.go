package space_translator

import (
	"math"

	"github.com/jhuggett/frontend/pages/world_map/camera"
)

type SpaceTranslator interface {
	// translates a point on the screen to a point in the world, based on camera, tile size, etc.
	FromScreenToWorld(x, y float64) (float64, float64)

	// translates a point in the world to a point on the screen
	FromWorldToScreen(x, y float64) (float64, float64)

	// translates a point from game data to world coordinates, e.g. scales by tile size
	FromDataToWorld(x, y float64) (float64, float64)

	// i.e. how zoomed in/out the camera is
	ScreenScale() (float64, float64)

	// returns the tile size
	TileSize() (float64, float64)

	// translates a point from world coordinates to game data coordinates, e.g. removes tile size scaling
	FromWorldToData(x, y float64) (float64, float64)
}

type spaceTranslator struct {
	camera   *camera.Camera
	tileSize float64
}

func (s *spaceTranslator) FromScreenToWorld(x, y float64) (float64, float64) {
	screenX := (x/s.camera.ZoomFactor + s.camera.Position[0])
	screenY := (y/s.camera.ZoomFactor + s.camera.Position[1])

	return screenX, screenY
}

func (s *spaceTranslator) FromWorldToScreen(x, y float64) (float64, float64) {
	screenX := (x - s.camera.Position[0])
	screenY := (y - s.camera.Position[1])

	return screenX, screenY
}

func (s *spaceTranslator) FromDataToWorld(x, y float64) (float64, float64) {
	worldX := x * s.tileSize
	worldY := y * s.tileSize

	return worldX, worldY
}

func (s *spaceTranslator) ScreenScale() (float64, float64) {
	return s.camera.ZoomFactor, s.camera.ZoomFactor
}

func (s *spaceTranslator) TileSize() (float64, float64) {
	return s.tileSize, s.tileSize
}

func (s *spaceTranslator) FromWorldToData(x, y float64) (float64, float64) {
	worldX := x / s.tileSize
	worldY := y / s.tileSize

	return worldX, worldY
}

func Floor(x, y float64) (float64, float64) {
	return math.Floor(x), math.Floor(y)
}

func ToInt(x, y float64) (int, int) {
	return int(x), int(y)
}

func New(camera *camera.Camera, tileSize float64) SpaceTranslator {
	return &spaceTranslator{
		camera:   camera,
		tileSize: tileSize,
	}
}
