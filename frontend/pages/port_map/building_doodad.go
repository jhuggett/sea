package port_map

import (
	"design-library/doodad"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/frontend/pages/world_map/camera"
	"github.com/jhuggett/frontend/utils/space_translator"
)

func NewBuildingDoodad(
	camera *camera.Camera,
	spaceTranslator space_translator.SpaceTranslator,
	building *game.Building,
	tileSize int,
) *BuildingDoodad {
	buildingDoodad := &BuildingDoodad{
		Camera:          camera,
		SpaceTranslator: spaceTranslator,
		Building:        building,
		TileSize:        tileSize,
	}

	return buildingDoodad
}

type BuildingDoodad struct {
	doodad.Default

	Camera          *camera.Camera
	SpaceTranslator space_translator.SpaceTranslator

	Building *game.Building

	TileSize int

	Image ebiten.Image
}

func (b *BuildingDoodad) Setup() {

	b.Image = *ebiten.NewImage(b.TileSize*b.Building.X, b.TileSize*b.Building.Y)
	b.Image.Fill(color.RGBA{100, 100, 200, 255})

	b.Children().Setup()
}

func (b *BuildingDoodad) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		b.SpaceTranslator.FromWorldToScreen(
			b.SpaceTranslator.FromDataToWorld(
				float64(b.Building.X*b.TileSize),
				float64(b.Building.Y*b.TileSize),
			),
		),
	)
	op.GeoM.Scale(
		b.SpaceTranslator.ScreenScale(),
	)
	screen.DrawImage(&b.Image, op)

	b.Children().Draw(screen)
}
