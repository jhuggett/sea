package world_map

import (
	"design-library/doodad"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/frontend/utils/space_translator"
)

func NewWorldMapDoodad(
	worldMap *game.WorldMap,
	spaceTranslator space_translator.SpaceTranslator,
) *WorldMapDoodad {
	worldMapDoodad := &WorldMapDoodad{
		WorldMap:        worldMap,
		SpaceTranslator: spaceTranslator,
	}

	return worldMapDoodad
}

type WorldMapDoodad struct {
	WorldMap *game.WorldMap

	Background *ebiten.Image

	SpaceTranslator space_translator.SpaceTranslator

	SmallestContinentPointX int
	SmallestContinentPointY int

	doodad.Default

	ContinentDoodads []*ContinentDoodad
}

func (w *WorldMapDoodad) Draw(screen *ebiten.Image) {

	// originX, originY := w.Origin()
	// scaleX, scaleY := w.Scale()

	// // Draw the background
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(
	// 	float64(w.SmallestContinentPointX)*float64(w.WorldMap.Manager.TileSize())-float64(originX),
	// 	float64(w.SmallestContinentPointY)*float64(w.WorldMap.Manager.TileSize())-float64(originY),
	// )

	// op.GeoM.Scale(
	// 	scaleX,
	// 	scaleY,
	// )
	// screen.DrawImage(w.Background, op)

	// Draw the background

	if w.Background != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(
			w.SpaceTranslator.FromWorldToScreen(
				w.SpaceTranslator.FromDataToWorld(
					float64(w.SmallestContinentPointX),
					float64(w.SmallestContinentPointY),
				),
			),
		)
		op.GeoM.Scale(
			w.SpaceTranslator.ScreenScale(),
		)
		screen.DrawImage(w.Background, op)
	}

	// Draw the doodads
	w.Children().Draw(screen)
}

func (w *WorldMapDoodad) Load() { // So setting up all the renderables doesn't block the UI thread
	if w.ContinentDoodads != nil {
		slog.Info("WorldMapDoodad Load: already loaded, skipping")
		return
	}

	w.ContinentDoodads = []*ContinentDoodad{}

	for _, continent := range w.WorldMap.Continents {
		continentDoodad := &ContinentDoodad{
			Continent:       continent,
			SpaceTranslator: w.SpaceTranslator,
		}

		continentDoodad.Load()

		w.ContinentDoodads = append(w.ContinentDoodads, continentDoodad)

		// w.Background = ebiten.NewImage(
		// 	int(float64(largestX+1)*float64(w.WorldMap.Manager.TileSize()))-int(float64(w.SmallestContinentPointX)*float64(w.WorldMap.Manager.TileSize())),
		// 	int(float64(largestY+1)*float64(w.WorldMap.Manager.TileSize()))-int(float64(w.SmallestContinentPointY)*float64(w.WorldMap.Manager.TileSize())),
		// )

		// w.Background.Fill(color.RGBA{
		// 	R: 0,
		// 	G: 0,
		// 	B: 0,
		// 	A: 255,
		// })

	}
	largestX := 0
	largestY := 0

	for _, continent := range w.WorldMap.Continents {
		if continent.LargestX > largestX {
			largestX = continent.LargestX
		}
		if continent.LargestY > largestY {
			largestY = continent.LargestY
		}
		if continent.SmallestX < w.SmallestContinentPointX {
			w.SmallestContinentPointX = continent.SmallestX
		}
		if continent.SmallestY < w.SmallestContinentPointY {
			w.SmallestContinentPointY = continent.SmallestY
		}
	}

	mapWidth := int(float64(largestX+1)*float64(w.WorldMap.Manager.TileSize())) - int(float64(w.SmallestContinentPointX)*float64(w.WorldMap.Manager.TileSize()))
	mapHeight := int(float64(largestY+1)*float64(w.WorldMap.Manager.TileSize())) - int(float64(w.SmallestContinentPointY)*float64(w.WorldMap.Manager.TileSize()))

	slog.Info("WorldMapDoodad Load", "mapWidth", mapWidth, "mapHeight", mapHeight)
}

func (w *WorldMapDoodad) Setup() {

	if w.WorldMap == nil {
		panic("WorldMap is nil, cannot setup WorldMapDoodad")
	}

	// for _, continent := range w.WorldMap.Continents {
	// 	continentDoodad := &ContinentDoodad{
	// 		Continent:       continent,
	// 		SpaceTranslator: w.SpaceTranslator,
	// 	}

	// 	w.AddChild(continentDoodad)
	// 	continentDoodad.Setup()
	// }

	for _, continentDoodad := range w.ContinentDoodads {
		w.AddChild(continentDoodad)
		continentDoodad.Setup()
	}

}

/*

mapImg := ebiten.NewImage(
		int(float64(page.LargestPointX-page.SmallestPointX+1)*page.TileSize),
		int(float64(page.LargestPointY-page.SmallestPointY+1)*page.TileSize),
	)

	mapImg.Fill(color.RGBA{
		R: 220,
		G: 202,
		B: 127,
		A: 255,
	})

	page.MapImage = mapImg

*/
