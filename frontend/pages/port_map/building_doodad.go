package port_map

import (
	"design-library/doodad"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

	Image *ebiten.Image
}

func (b *BuildingDoodad) Setup() {

	b.Image = ebiten.NewImage(b.TileSize+10, b.TileSize+10)
	b.Image.Fill(color.RGBA{100, 100, 200, 255})

	textImg := ebiten.NewImage(b.TileSize+10, 20)
	ebitenutil.DebugPrintAt(textImg, b.Building.Name, 0, 0)
	b.Image.DrawImage(textImg, &ebiten.DrawImageOptions{})

	b.Children().Setup()
}

func (b *BuildingDoodad) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		b.SpaceTranslator.FromWorldToScreen(
			b.SpaceTranslator.FromDataToWorld(
				float64(b.Building.X),
				float64(b.Building.Y),
			),
		),
	)
	op.GeoM.Scale(
		b.SpaceTranslator.ScreenScale(),
	)
	screen.DrawImage(b.Image, op)

	b.Children().Draw(screen)
}
