package world_map

import (
	"design-library/doodad"
	"design-library/reaction"
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/outbound"
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
	w.DoOnTeardown(w.Ship.Manager.OnShipMovedCallback.Register(func(smr outbound.ShipMovedReq) error {
		w.Ship.Route.ShipMovedReq = &smr

		if w.Ship.Route.ShipMovedReq.RouteInfo.ReachedDestination {
			w.Ship.Route.Active = false
		}

		w.RenderRoute()
		return nil
	}))

	w.Reactions().Add(
		reaction.NewMouseUpReaction(
			doodad.MouseMovedWithin[*reaction.MouseUpEvent](w),
			func(mm *reaction.MouseUpEvent) {
				if mm.Button != ebiten.MouseButtonLeft {
					return
				}

				if w.Ship.IsRouteActive() {
					slog.Info("RouteDoodad.OnClick: Ship already has an active route")
					return
				}

				worldX, worldY := w.SpaceTranslator.FromScreenToWorld(
					float64(mm.X),
					float64(mm.Y),
				)

				dataX, dataY := w.SpaceTranslator.FromWorldToData(
					worldX,
					worldY,
				)

				_, err := w.Ship.PlotRoute(int(dataX), int(dataY))
				if err != nil {
					slog.Error("RouteDoodad.OnClick", "error", err)
					return
				}

				w.RenderRoute()

				mm.StopPropagation()
			},
		),
	)
}

func (w *RouteDoodad) RenderRoute() {

	scale, _ := w.SpaceTranslator.TileSize()
	w.img = ebiten.NewImage(Box(w.Ship.Route.Points, float64(scale)))

	tileSize, _ := w.SpaceTranslator.TileSize()
	smallestX, smallestY := 0.0, 0.0

	for _, point := range w.Ship.Route.Points {
		if point.X < smallestX {
			smallestX = point.X
		}
		if point.Y < smallestY {
			smallestY = point.Y
		}
	}

	w.originX = int(smallestX)
	w.originY = int(smallestY)

	routeColor := color.RGBA{
		R: 255,
		A: 255,
	}
	if w.Ship.IsRouteActive() {
		routeColor = color.RGBA{
			G: 255,
			A: 255,
		}
	}

	for _, point := range w.Ship.Route.Points {
		vector.DrawFilledRect(
			w.img,
			float32((point.X-smallestX)*tileSize+tileSize/4),
			float32((point.Y-smallestY)*tileSize+tileSize/4),
			float32(tileSize/2),
			float32(tileSize/2),
			routeColor,
			false,
		)
	}
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
