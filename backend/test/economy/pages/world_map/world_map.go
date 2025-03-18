package world_map

import (
	"image/color"
	"log/slog"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/outbound"
	"github.com/jhuggett/sea/start"
	"github.com/jhuggett/sea/test/economy/widgets"
	"github.com/jhuggett/sea/timeline"
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
	Conn *start.Connection

	PlottedRoute PlottedRoute

	Camera Camera
	Canvas *ebiten.Image

	TileSize float64

	Continents []Continent
	Ports      []Port
	Ship       Ship

	ui *ebitenui.UI

	Press *Press

	MouseLocationX, MouseLocationY int

	Width, Height int

	SmallestPointX, SmallestPointY int
	LargestPointX, LargestPointY   int

	MapImage *ebiten.Image

	OnTileClicked func(x, y int)

	CursorImage *ebiten.Image
}

func New(snapshot *game_context.Snapshot) (*WorldMapPage, error) {
	// This creates the root container for this UI.
	// All other UI elements must be added to this container.
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			// widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	// This adds the root container to the UI, so that it will be rendered.
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	// s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	// if err != nil {
	// 	return nil, err
	// }

	// fontFace := &text.GoTextFace{
	// 	Source: s,
	// 	Size:   32,
	// }
	// This creates a text widget that says "Hello World!"
	// helloWorldLabel := widget.NewText(
	// 	widget.TextOpts.Text("Hello World!", fontFace, color.White),
	// 	widget.TextOpts.MaxWidth(20),
	// )

	// To display the text widget, we have to add it to the root container.
	// rootContainer.AddChild(helloWorldLabel)

	page := &WorldMapPage{
		Camera: Camera{
			ViewPort: f64.Vec2{800, 500},
			TileSize: 32,
		},
		TileSize: 32,

		ui: eui,
	}

	conn := &start.Connection{
		Receiver: &outbound.Receiver{
			OnCrewInformation: func(cir outbound.CrewInformationReq) (outbound.CrewInformationResp, error) {
				slog.Info("CrewInformation called", "req", cir)
				return outbound.CrewInformationResp{}, nil
			},
			OnShipMoved: func(smr outbound.ShipMovedReq) (outbound.ShipMovedResp, error) {
				slog.Info("ShipMoved called", "req", smr)

				page.Ship.X = float64(smr.Location.X)
				page.Ship.Y = float64(smr.Location.Y)

				return outbound.ShipMovedResp{}, nil
			},
			OnShipDocked: func(sdr outbound.ShipDockedReq) (outbound.ShipDockedResp, error) {
				slog.Info("ShipDocked called", "req", sdr)
				return outbound.ShipDockedResp{}, nil
			},
			OnTimeChanged: func(tcr outbound.TimeChangedReq) (outbound.TimeChangedResp, error) {
				slog.Info("TimeChanged called", "req", tcr)
				return outbound.TimeChangedResp{}, nil
			},
			OnShipInventoryChanged: func(sicr outbound.ShipInventoryChangedReq) (outbound.ShipInventoryChangedResp, error) {
				slog.Info("ShipInventoryChanged called", "req", sicr)
				return outbound.ShipInventoryChangedResp{}, nil
			},
			OnShipChanged: func(scr outbound.ShipChangedReq) (outbound.ShipChangedResp, error) {
				slog.Info("ShipChanged called", "req", scr)
				return outbound.ShipChangedResp{}, nil
			},
		},
		GameCtx: &game_context.GameContext{},
	}

	exampleImage, _, err := ebitenutil.NewImageFromFile("./assets/images/tile_0006.png")
	if err != nil {
		return nil, err
	}
	//

	if snapshot == nil {
		resp, err := inbound.Register(inbound.RegisterReq{})
		if err != nil {
			return nil, err
		}

		snapshot = &resp.GameCtx

		slog.Info("Register response", "resp", resp)
	}

	conn.GameCtx = game_context.New(*snapshot)

	loginResp, err := inbound.Login(inbound.LoginReq{
		Snapshot: *snapshot,
	}, conn)
	if err != nil {
		return nil, err
	}

	slog.Info("Login response", "loginResp", loginResp)

	conn.GameCtx.Timeline = timeline.New()

	start.Game(conn)

	page.Ship = Ship{
		X:     loginResp.Ship.X,
		Y:     loginResp.Ship.Y,
		Image: exampleImage,
	}

	shipImg := ebiten.NewImage(
		int(page.TileSize),
		int(page.TileSize),
	)
	vector.DrawFilledRect(
		shipImg,
		0,
		0,
		float32(page.TileSize),
		float32(page.TileSize),
		color.RGBA{
			B: 255,
			G: 255,
			A: 255,
		},
		false,
	)

	page.Ship.Image = shipImg

	gameMap, err := inbound.GetWorldMap(inbound.GetWorldMapReq{}, snapshot.GameMapID)
	if err != nil {
		return nil, err
	}

	slog.Info("GetWorldMap response", "gameMap", gameMap)

	// lowestColor := color.RGBA{
	// 	R: 100,
	// 	G: 100,
	// 	B: 100,
	// 	A: 255,
	// }
	// highestColor := color.RGBA{
	// 	R: 220,
	// 	G: 202,
	// 	B: 127,
	// 	A: 255,
	// }

	a := uint8(50)
	b := uint8(100)
	lowestColor := color.RGBA{
		R: 23 + a,
		G: 18 + a,
		B: 13 + a,
		A: 255,
	}
	highestColor := color.RGBA{
		R: 58 + b,
		G: 45 + b,
		B: 25 + b,
		A: 255,
	}

	// rgb(179, 146, 83)
	// rgb(86, 48, 16)
	page.Continents = []Continent{}
	for _, continent := range gameMap.Continents {
		smallestX := 0
		smallestY := 0
		largestX := 0
		largestY := 0

		for _, point := range continent.Points {
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

			if page.SmallestPointX > point.X {
				page.SmallestPointX = point.X
			}
			if page.SmallestPointY > point.Y {
				page.SmallestPointY = point.Y
			}
			if page.LargestPointX < point.X {
				page.LargestPointX = point.X
			}
			if page.LargestPointY < point.Y {
				page.LargestPointY = point.Y
			}
		}

		img := ebiten.NewImage(
			(largestX-smallestX+1)*int(page.TileSize),
			(largestY-smallestY+1)*int(page.TileSize),
		)
		for _, point := range continent.Points {
			e := (point.Elevation - .5) * 2
			vector.DrawFilledRect(
				img,
				float32(float64(point.X-smallestX)*page.TileSize),
				float32(float64(point.Y-smallestY)*page.TileSize),
				float32(page.TileSize),
				float32(page.TileSize),
				color.RGBA{
					R: uint8(float64(highestColor.R-lowestColor.R)*e) + lowestColor.R,
					G: uint8(float64(highestColor.G-lowestColor.G)*e) + lowestColor.G,
					B: uint8(float64(highestColor.B-lowestColor.B)*e) + lowestColor.B,
					A: 255,
				},
				false,
			)
		}
		page.Continents = append(page.Continents, Continent{
			OriginX: smallestX,
			OriginY: smallestY,
			Image:   img,
		})
	}

	portsData, err := inbound.GetPorts(inbound.GetPortsReq{}, snapshot.GameMapID)
	if err != nil {
		return nil, err
	}

	slog.Info("GetPorts response", "portsData", portsData)

	page.Ports = []Port{}

	for _, port := range portsData.Ports {
		img := ebiten.NewImage(int(page.TileSize), int(page.TileSize))
		vector.DrawFilledRect(
			img,
			float32(page.TileSize/4),
			float32(page.TileSize/4),
			float32(page.TileSize/2),
			float32(page.TileSize/2),
			color.RGBA{
				R: 255,
				A: 255,
			},
			false,
		)
		page.Ports = append(page.Ports, Port{
			X:     port.Point.X,
			Y:     port.Point.Y,
			Image: img,
		})
	}

	// should add border
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

	page.OnTileClicked = func(x, y int) {
		resp, err := inbound.PlotRoute(inbound.PlotRouteReq{
			Coordinate: inbound.Coordinate{
				X: float64(x),
				Y: float64(y),
			},
		}, conn)
		if err != nil {
			slog.Error("Failed to plot route", "error", err)
			return
		}
		page.PlottedRoute = PlottedRoute{
			Points:                resp.Coordinates,
			EstimatedSailingSpeed: resp.EstimatedSailingSpeed,
			EstimatedDuration:     resp.EstimatedDuration,
			EndTileX:              x,
			EndTileY:              y,
			IsActive:              false,
		}
		page.PlottedRoute.RenderImage(page.TileSize)
	}

	page.CursorImage = ebiten.NewImage(int(page.TileSize), int(page.TileSize))
	vector.DrawFilledRect(
		page.CursorImage,
		0,
		0,
		float32(page.TileSize),
		float32(page.TileSize),
		color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		},
		false,
	)

	//

	page.Canvas = ebiten.NewImage(
		int(float64(page.LargestPointX-page.SmallestPointX+1)*page.TileSize),
		int(float64(page.LargestPointY-page.SmallestPointY+1)*page.TileSize),
	)

	button := (&widgets.Button{
		Text: "Set Sail",
		OnClick: func() {
			slog.Info("Button clicked")

			slog.Info("MoveShip called", "conn", conn.Context())

			resp, err := inbound.MoveShip(inbound.MoveShipReq{
				X: float64(page.PlottedRoute.EndTileX),
				Y: float64(page.PlottedRoute.EndTileY),
			}, conn)

			if err != nil {
				slog.Error("Failed to move ship", "error", err)
				return
			}

			slog.Info("MoveShip response", "resp", resp)
		},
	})

	button.Setup(rootContainer)

	return page, nil
}

func (w *WorldMapPage) SetWidthAndHeight(width, height int) {
	w.Width = width
	w.Height = height
}
