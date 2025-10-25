package port_map

import (
	"design-library/doodad"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/frontend/pages/world_map/camera"
	"github.com/jhuggett/frontend/utils/space_translator"
)

func NewPortMapDoodad(
	camera *camera.Camera,
	spaceTranslator space_translator.SpaceTranslator,
	port *game.Port,
	buildings []*game.Building,
	tileSize int,
) *PortMapDoodad {
	portMapDoodad := &PortMapDoodad{
		Camera:          camera,
		SpaceTranslator: spaceTranslator,
		Port:            port,
		Buildings:       buildings,
		TileSize:        tileSize,
	}

	return portMapDoodad
}

type PortMapDoodad struct {
	doodad.Default

	Camera          *camera.Camera
	SpaceTranslator space_translator.SpaceTranslator

	Port      *game.Port
	Buildings []*game.Building

	Background *ebiten.Image

	SmallestPointX int
	SmallestPointY int

	TileSize int
}

func (p *PortMapDoodad) Setup() {
	largestX := 0
	largestY := 0

	for _, building := range p.Buildings {
		if building.X > largestX {
			largestX = building.X
		}
		if building.Y > largestY {
			largestY = building.Y
		}
		if building.X < p.SmallestPointX {
			p.SmallestPointX = building.X
		}
		if building.Y < p.SmallestPointY {
			p.SmallestPointY = building.Y
		}
	}

	tileSize := float64(p.TileSize)

	p.Background = ebiten.NewImage(
		int(float64(largestX+1)*tileSize)-int(float64(p.SmallestPointX)*tileSize),
		int(float64(largestY+1)*tileSize)-int(float64(p.SmallestPointY)*tileSize),
	)

	// w.Background.Fill(color.RGBA{
	// 	R: 220,
	// 	G: 202,
	// 	B: 127,
	// 	A: 20,
	// })

	p.Background.Fill(color.RGBA{
		R: 25,
		G: 0,
		B: 0,
		A: 255,
	})

	for _, building := range p.Buildings {
		b := building
		buildingDoodad := NewBuildingDoodad(
			p.Camera,
			p.SpaceTranslator,
			b,
			p.TileSize,
		)
		p.AddChild(buildingDoodad)
	}

	p.Children().Setup()
}

func (p *PortMapDoodad) Draw(screen *ebiten.Image) {
	if p.Background != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(
			p.SpaceTranslator.FromWorldToScreen(
				p.SpaceTranslator.FromDataToWorld(
					float64(p.SmallestPointX),
					float64(p.SmallestPointY),
				),
			),
		)
		op.GeoM.Scale(
			p.SpaceTranslator.ScreenScale(),
		)
		screen.DrawImage(p.Background, op)
	}

	// Draw the doodads
	p.Children().Draw(screen)
}
