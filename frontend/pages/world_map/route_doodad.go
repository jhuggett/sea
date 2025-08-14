package world_map

import (
	"design-library/doodad"
	"design-library/reaction"
	"fmt"
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/sea/inbound"
)

func NewRouteDoodad(
	spaceTranslator SpaceTranslator,
	ship *game.Ship,
) *RouteDoodad {
	routeDoodad := &RouteDoodad{
		SpaceTranslator: spaceTranslator,
		Ship:            ship,
	}

	return routeDoodad
}

type RouteDoodad struct {
	SpaceTranslator SpaceTranslator
	Ship            *game.Ship

	img              *ebiten.Image
	originX, originY int

	doodad.Default
}

func (w *RouteDoodad) Update() error {
	return nil
}

func (w *RouteDoodad) Draw(screen *ebiten.Image) {
	if w.img == nil {
		return
	}

	x, y := w.SpaceTranslator.FromDataToWorld(
		float64(w.originX),
		float64(w.originY),
	)

	x, y = w.SpaceTranslator.FromWorldToScreen(
		float64(x),
		float64(y),
	)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		x,
		y,
	)
	op.GeoM.Scale(w.SpaceTranslator.ScreenScale())
	screen.DrawImage(w.img, op)

	w.Children().Draw(screen)
}

func (w *RouteDoodad) Setup() {
	w.Reactions().Add(
		reaction.NewMouseUpReaction(
			doodad.MouseMovedWithin[reaction.MouseUpEvent](w),
			func(mm reaction.MouseUpEvent) {
				if mm.Button != ebiten.MouseButtonLeft {
					return
				}

				x, y := mm.X, mm.Y
				fmt.Println("RouteDoodad.OnClick", x, y)

				worldX, worldY := w.SpaceTranslator.FromScreenToWorld(
					float64(x),
					float64(y),
				)

				dataX, dataY := w.SpaceTranslator.FromWorldToData(
					worldX,
					worldY,
				)

				route, err := w.Ship.PlotRoute(int(dataX), int(dataY))
				if err != nil {
					slog.Error("RouteDoodad.OnClick", "error", err)
					return
				}

				scale, _ := w.SpaceTranslator.TileSize()
				w.img = ebiten.NewImage(Box(route.Points, float64(scale)))

				// For Debugging
				// w.img.Fill(color.RGBA{
				// 	R: 0,
				// 	G: 100,
				// 	B: 0,
				// 	A: 125,
				// })

				tileSize, _ := w.SpaceTranslator.TileSize()
				smallestX, smallestY := 0.0, 0.0

				for _, point := range route.Points {
					if point.X < smallestX {
						smallestX = point.X
					}
					if point.Y < smallestY {
						smallestY = point.Y
					}
				}

				w.originX = int(smallestX)
				w.originY = int(smallestY)

				for _, point := range route.Points {
					vector.DrawFilledRect(
						w.img,
						float32((point.X-smallestX)*tileSize+tileSize/4),
						float32((point.Y-smallestY)*tileSize+tileSize/4),
						float32(tileSize/2),
						float32(tileSize/2),
						color.RGBA{
							B: 255,
							A: 255,
						},
						false,
					)
				}
			},
		),
	)
}

func Box(points []inbound.Coordinate, scale float64) (width, height int) {
	smallestX := 0.0
	smallestY := 0.0
	largestX := 0.0
	largestY := 0.0

	for _, point := range points {
		if point.X < smallestX {
			smallestX = point.X
		}
		if point.Y < smallestY {
			smallestY = point.Y
		}
		if point.X > largestX {
			largestX = point.X
		}
		if point.Y > largestY {
			largestY = point.Y
		}
	}

	width = int(largestX - smallestX + 1)
	height = int(largestY - smallestY + 1)

	width = int(float64(width) * scale)
	height = int(float64(height) * scale)

	return width, height
}
