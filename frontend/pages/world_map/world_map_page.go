package world_map

import (
	"image/color"
	"log/slog"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jhuggett/frontend/doodad"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/inbound"
	"golang.org/x/image/math/f64"
)

type PlottedRoute struct {
	Points []inbound.Coordinate

	EndTileX, EndTileY int

	EstimatedSailingSpeed float64
	EstimatedDuration     float64 // in days

	OriginX, OriginY int
	Image            *ebiten.Image

	IsActive bool
}

func (p *PlottedRoute) RenderImage(tileSize float64) {
	smallestX, smallestY := 0.0, 0.0
	largestX, largestY := 0.0, 0.0

	for _, point := range p.Points {
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

	p.OriginX = int(smallestX)
	p.OriginY = int(smallestY)

	img := ebiten.NewImage(
		int((largestX-smallestX+1)*tileSize),
		int((largestY-smallestY+1)*tileSize),
	)

	for _, point := range p.Points {
		vector.DrawFilledRect(
			img,
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

	p.Image = img
}

type Press struct {
	StartX, StartY int
	X, Y           int
	TimeStart      time.Time
}

// type Camera struct {
// 	X, Y float64
// 	Zoom float64
// }

type Continent struct {
	OriginX, OriginY int

	Image *ebiten.Image
}

type Port struct {
	X, Y  int
	Image *ebiten.Image
}

type Ship struct {
	X, Y float64

	Image *ebiten.Image
}

// type Connection struct {
// 	Receiver outbound.Receiver
// 	gameCtx  *game_context.GameContext
// }

// func (c *Connection) Context() *game_context.GameContext {
// 	return c.gameCtx
// }

// func (c *Connection) Sender() *outbound.Sender {
// 	return outbound.NewSender(c.gameCtx, c.Receiver)
// }

type WorldMapPage struct {
	// Conn *start.Connection

	Doodads  []doodad.Doodad
	Gesturer doodad.Gesturer

	// PlottedRoute PlottedRoute

	Camera   *Camera
	TileSize float64

	// Canvas *ebiten.Image

	// Continents []Continent
	// Ports      []Port
	// Ship       Ship

	Press *Press

	MouseLocationX, MouseLocationY int

	Width, Height int

	SpaceTranslator SpaceTranslator

	// SmallestPointX, SmallestPointY int
	// LargestPointX, LargestPointY   int

	// MapImage *ebiten.Image

	// OnTileClicked func(x, y int)

	// CursorImage *ebiten.Image
}

func New(snapshot *game_context.Snapshot) (*WorldMapPage, error) {

	p := &WorldMapPage{
		Camera: &Camera{
			ViewPort:   f64.Vec2{800, 500},
			TileSize:   32,
			Position:   f64.Vec2{0, 0},
			ZoomFactor: 1,
		},
		TileSize: 32,
	}

	p.SpaceTranslator = &spaceTranslator{
		camera:   p.Camera,
		tileSize: p.TileSize,
	}

	p.Gesturer = doodad.NewGesturer()

	// page.OnTileClicked = func(x, y int) {
	// 	resp, err := inbound.PlotRoute(inbound.PlotRouteReq{
	// 		Coordinate: inbound.Coordinate{
	// 			X: float64(x),
	// 			Y: float64(y),
	// 		},
	// 	}, conn)
	// 	if err != nil {
	// 		slog.Error("Failed to plot route", "error", err)
	// 		return
	// 	}
	// 	page.PlottedRoute = PlottedRoute{
	// 		Points:                resp.Coordinates,
	// 		EstimatedSailingSpeed: resp.EstimatedSailingSpeed,
	// 		EstimatedDuration:     resp.EstimatedDuration,
	// 		EndTileX:              x,
	// 		EndTileY:              y,
	// 		IsActive:              false,
	// 	}
	// 	page.PlottedRoute.RenderImage(page.TileSize)
	// }

	// page.CursorImage = ebiten.NewImage(int(page.TileSize), int(page.TileSize))
	// vector.DrawFilledRect(
	// 	page.CursorImage,
	// 	0,
	// 	0,
	// 	float32(page.TileSize),
	// 	float32(page.TileSize),
	// 	color.RGBA{
	// 		R: 0,
	// 		G: 0,
	// 		B: 0,
	// 		A: 255,
	// 	},
	// 	false,
	// )

	//

	// page.Canvas = ebiten.NewImage(
	// 	int(float64(page.LargestPointX-page.SmallestPointX+1)*page.TileSize),
	// 	int(float64(page.LargestPointY-page.SmallestPointY+1)*page.TileSize),
	// )

	// button := (&widgets.Button{
	// 	Text: "Set Sail",
	// 	OnClick: func() {
	// 		slog.Info("Button clicked")

	// 		slog.Info("MoveShip called", "conn", conn.Context())

	// 		resp, err := inbound.MoveShip(inbound.MoveShipReq{
	// 			X: float64(page.PlottedRoute.EndTileX),
	// 			Y: float64(page.PlottedRoute.EndTileY),
	// 		}, conn)

	// 		if err != nil {
	// 			slog.Error("Failed to move ship", "error", err)
	// 			return
	// 		}

	// 		slog.Info("MoveShip response", "resp", resp)
	// 	},
	// })

	// buttonContainer := widget.NewContainer(
	// 	widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255})),
	// 	widget.ContainerOpts.WidgetOpts(
	// 		//The widget in this cell has a MaxHeight and MaxWidth less than the
	// 		//Size of the grid cell so it will use the Position fields below to
	// 		//Determine where the widget should be displayed within that grid cell.
	// 		widget.WidgetOpts.LayoutData(widget.GridLayoutData{}),
	// 		widget.WidgetOpts.MinSize(50, 50),
	// 	),
	// )

	// button.Setup(buttonContainer)

	// rootContainer.AddChild(buttonContainer)

	gameManager := game.NewManager(snapshot)

	gameManager.SetTileSize(32)

	worldMapDoodad := &WorldMapDoodad{
		SpaceTranslator: p.SpaceTranslator,
		Gesturer:        p.Gesturer,
	}

	loadingDoodad := &LoadingDoodad{}

	p.Doodads = []doodad.Doodad{
		loadingDoodad,
	}

	loadingDoodad.Setup()
	loadingDoodad.SetMessage("Loading world map...")
	loadingDoodad.Show()

	p.Gesturer.OnMouseDrag(func(lastX, lastY, currentX, currentY int) error {
		p.Camera.Position[0] += float64(lastX-currentX) / p.Camera.ZoomFactor
		p.Camera.Position[1] += float64(lastY-currentY) / p.Camera.ZoomFactor
		return nil
	})

	p.Gesturer.OnMouseWheel(func(offset float64) error {
		if offset > 0 {
			p.Camera.ZoomFactor += 0.1
		} else {
			p.Camera.ZoomFactor -= 0.1
			if p.Camera.ZoomFactor < 0.1 {
				p.Camera.ZoomFactor = 0.1
			}
		}
		return nil
	})

	go func() {
		err := gameManager.Start()
		if err != nil {
			slog.Error("Failed to start game manager", "error", err)
			loadingDoodad.SetMessage("Failed to load world map")
			return
		}

		worldMapDoodad.WorldMap = gameManager.WorldMap

		err = worldMapDoodad.Setup()
		if err != nil {
			slog.Error("Failed to setup world map doodad", "error", err)
			loadingDoodad.SetMessage("Failed to load world map")
			return
		}
		loadingDoodad.Hide()

		p.Doodads = []doodad.Doodad{
			worldMapDoodad,
		}

		// --- Ship

		shipDoodad := &ShipDoodad{
			Origin: func() (int, int) {
				return int(p.Camera.Position[0]), int(p.Camera.Position[1])
			},
			Scale: func() (float64, float64) {
				return p.Camera.ZoomFactor, p.Camera.ZoomFactor
			},
			Ship: gameManager.PlayerShip,
		}
		err = shipDoodad.Setup()
		if err != nil {
			slog.Error("Failed to setup ship doodad", "error", err)
			loadingDoodad.SetMessage("Failed to load world map")
			return
		}
		p.Doodads = append(p.Doodads, shipDoodad)

		// --- Route Doodad

		routeDoodad := &RouteDoodad{
			SpaceTranslator: p.SpaceTranslator,
			Gesturer:        p.Gesturer,
			Ship:            gameManager.PlayerShip,
		}
		err = routeDoodad.Setup()
		if err != nil {
			slog.Error("Failed to setup route doodad", "error", err)
			loadingDoodad.SetMessage("Failed to load world map")
			return
		}
		p.Doodads = append(p.Doodads, routeDoodad)

		// --- Cursor Doodad

		cursorDoodad := &CursorDoodad{
			Gesturer:        p.Gesturer,
			SpaceTranslator: p.SpaceTranslator,
		}
		err = cursorDoodad.Setup()
		if err != nil {
			slog.Error("Failed to setup cursor doodad", "error", err)
			loadingDoodad.SetMessage("Failed to load world map")
			return
		}
		p.Doodads = append(p.Doodads, cursorDoodad)

		timeControlDoodad := &TimeControlDoodad{
			Gesturer: p.Gesturer,
			Manager:  gameManager,
			Position: func() doodad.Position {
				return doodad.Position{
					X: p.Width - 200,
					Y: p.Height - 200,
				}
			},
		}
		err = timeControlDoodad.Setup()
		if err != nil {
			slog.Error("Failed to setup time control doodad", "error", err)
			loadingDoodad.SetMessage("Failed to load world map")
			return
		}
		p.Doodads = append(p.Doodads, timeControlDoodad)

		tileInfoDoodad := &TileInformationDoodad{
			SpaceTranslator: p.SpaceTranslator,
			Gesturer:        p.Gesturer,

			HideTileCursor: func() {
				cursorDoodad.Hidden = true
			},
			ShowTileCursor: func() {
				cursorDoodad.Hidden = false
			},
		}
		err = tileInfoDoodad.Setup()
		if err != nil {
			slog.Error("Failed to setup tile info doodad", "error", err)
			loadingDoodad.SetMessage("Failed to load world map")
			return
		}
		p.Doodads = append(p.Doodads, tileInfoDoodad)
	}()

	return p, nil
}

func (w *WorldMapPage) SetWidthAndHeight(width, height int) {
	w.Width = width
	w.Height = height
}
