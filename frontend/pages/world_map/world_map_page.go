package world_map

import (
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/frontend/pages/world_map/camera"
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/inbound"
	"golang.org/x/image/math/f64"
)

func (w *WorldMapPage) CanvasTranslateTo(x, y float64) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}

	return op
}

func (w *WorldMapPage) CanvasTranslateFrom(x, y float64) (float64, float64) {
	return 0, 0
}

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

type WorldMapPage struct {
	// Children doodad.Children
	// Gesturer doodad.Gesturer

	Camera          *camera.Camera
	SpaceTranslator SpaceTranslator

	GameSnapshot *game_context.Snapshot
	GameManager  *game.Manager

	// TileSize float64

	// Press *Press

	// MouseLocationX, MouseLocationY int

	// Width, Height int

	doodad.Default

	loaded bool
}

func (w *WorldMapPage) SetupLoadingScreen() {
	loadingMessage := label.New(label.Config{
		Message:  "Loading...",
		Layout:   box.Zeroed(),
		FontSize: 24,
	})
	loadingMessage.Setup()
	loadingMessage.Box.Computed(func(b *box.Box) *box.Box {
		return b.CenterWithin(w.Box)
	})

	w.AddChild(loadingMessage)
	w.Children.Setup()
}

func (w *WorldMapPage) LoadRequiredData() error {
	// 	err := gameManager.Start()
	// 	if err != nil {
	// 		slog.Error("Failed to start game manager", "error", err)
	// 		loadingDoodad.SetMessage("Failed to load world map")
	// 		return
	// 	}

	// 	worldMapDoodad.WorldMap = gameManager.WorldMap

	w.Camera = &camera.Camera{
		ViewPort:   f64.Vec2{float64(w.Box.Width()), float64(w.Box.Height())},
		TileSize:   32,
		Position:   f64.Vec2{0, 0},
		ZoomFactor: 1,
	}

	w.SpaceTranslator = &spaceTranslator{
		camera:   w.Camera,
		tileSize: w.Camera.TileSize,
	}

	gameManager := game.NewManager(w.GameSnapshot)
	gameManager.SetTileSize(int(w.Camera.TileSize))

	w.Gesturer.OnMouseDrag(func(lastX, lastY, currentX, currentY int) error {
		w.Camera.Position[0] += float64(lastX-currentX) / w.Camera.ZoomFactor
		w.Camera.Position[1] += float64(lastY-currentY) / w.Camera.ZoomFactor
		return nil
	})

	w.Gesturer.OnMouseWheel(func(offset float64) error {
		if offset > 0 {
			w.Camera.ZoomFactor += 0.1
		} else {
			w.Camera.ZoomFactor -= 0.1
			if w.Camera.ZoomFactor < 0.1 {
				w.Camera.ZoomFactor = 0.1
			}
		}
		return nil
	})

	err := gameManager.Start()
	if err != nil {
		return fmt.Errorf("failed to start game manager: %w", err)
	}

	w.GameManager = gameManager

	w.loaded = true
	return nil
}

func (w *WorldMapPage) SetupMainScreen() {

	// worldMapDoodad := &WorldMapDoodad{
	// 	SpaceTranslator: p.SpaceTranslator,
	// 	Gesturer:        p.Gesturer,
	// }

	worldMapDoodad := NewWorldMapDoodad(
		w.GameManager.WorldMap,
		w.SpaceTranslator,
		w.Gesturer,
	)
	w.AddChild(worldMapDoodad)

	// Need to preload the things we need before we get to this method
	// then we can just render

	timeControlDoodad := NewTimeControlDoodad(
		w.GameManager,
		w.Gesturer,
	)
	timeControlDoodad.Layout().Computed(func(b *box.Box) *box.Box {
		return b.Copy(w.Box)
	})
	w.AddChild(timeControlDoodad)

	w.AddChild(NewPlayerShipDoodad(
		w.GameManager,
		w.Camera,
	))

	bottomBar := NewBottomBar(
		w.Gesturer,
		box.Computed(func(b *box.Box) *box.Box {
			return b.Copy(w.Box)
		}),
	)
	w.AddChild(bottomBar)

	w.AddChild(NewCursorDoodad(
		w.Gesturer,
		w.SpaceTranslator,
	))

	w.AddChild(NewRouteDoodad(
		w.Gesturer,
		w.SpaceTranslator,
		w.GameManager.PlayerShip,
	))

	w.Children.Setup()
}

func (w *WorldMapPage) Update() error {
	w.Gesturer.Update()
	return nil
}

func (w *WorldMapPage) Setup() {
	w.Children = doodad.NewChildren()

	if w.loaded {
		// show map
		w.SetupMainScreen()
	} else {
		// show loading screen
		w.SetupLoadingScreen()
		go func() {
			err := w.LoadRequiredData()
			if err != nil {
				panic(fmt.Errorf("failed to load required data: %w", err))
			}
			w.Setup() // recursively call Setup to switch to main screen
		}()
	}
}

func New(snapshot *game_context.Snapshot) *WorldMapPage {
	p := &WorldMapPage{
		Default: doodad.Default{
			Children: doodad.NewChildren(),
			Gesturer: doodad.NewGesturer(),
			Box:      box.Zeroed(),
		},
		GameSnapshot: snapshot,
	}

	p.Setup()

	// p := &WorldMapPage{
	// 	Camera: &Camera{
	// 		ViewPort:   f64.Vec2{800, 500},
	// 		TileSize:   32,
	// 		Position:   f64.Vec2{0, 0},
	// 		ZoomFactor: 1,
	// 	},
	// 	TileSize: 32,
	// }

	// p.SpaceTranslator = &spaceTranslator{
	// 	camera:   p.Camera,
	// 	tileSize: p.TileSize,
	// }

	// p.Gesturer = doodad.NewGesturer()

	// gameManager := game.NewManager(snapshot)

	// gameManager.SetTileSize(32)

	// worldMapDoodad := &WorldMapDoodad{
	// 	SpaceTranslator: p.SpaceTranslator,
	// 	Gesturer:        p.Gesturer,
	// }

	// loadingDoodad := &LoadingDoodad{}

	// p.Children = doodad.Children{}
	// p.Children.Add(loadingDoodad)

	// loadingDoodad.Setup()
	// loadingDoodad.SetMessage("Loading world map...")
	// loadingDoodad.Show()

	// p.Gesturer.OnMouseDrag(func(lastX, lastY, currentX, currentY int) error {
	// 	p.Camera.Position[0] += float64(lastX-currentX) / p.Camera.ZoomFactor
	// 	p.Camera.Position[1] += float64(lastY-currentY) / p.Camera.ZoomFactor
	// 	return nil
	// })

	// p.Gesturer.OnMouseWheel(func(offset float64) error {
	// 	if offset > 0 {
	// 		p.Camera.ZoomFactor += 0.1
	// 	} else {
	// 		p.Camera.ZoomFactor -= 0.1
	// 		if p.Camera.ZoomFactor < 0.1 {
	// 			p.Camera.ZoomFactor = 0.1
	// 		}
	// 	}
	// 	return nil
	// })

	// go func() {
	// 	err := gameManager.Start()
	// 	if err != nil {
	// 		slog.Error("Failed to start game manager", "error", err)
	// 		loadingDoodad.SetMessage("Failed to load world map")
	// 		return
	// 	}

	// 	worldMapDoodad.WorldMap = gameManager.WorldMap

	// 	err = worldMapDoodad.Setup()
	// 	if err != nil {
	// 		slog.Error("Failed to setup world map doodad", "error", err)
	// 		loadingDoodad.SetMessage("Failed to load world map")
	// 		return
	// 	}

	// 	err = p.Children.Clear()
	// 	if err != nil {
	// 		slog.Error("Failed to clear children", "error", err)
	// 		loadingDoodad.SetMessage("Failed to load world map")
	// 		return
	// 	}

	// 	p.Children.Add(worldMapDoodad)

	// 	// --- Ship

	// 	shipDoodad := &ShipDoodad{
	// 		Origin: func() (int, int) {
	// 			return int(p.Camera.Position[0]), int(p.Camera.Position[1])
	// 		},
	// 		Scale: func() (float64, float64) {
	// 			return p.Camera.ZoomFactor, p.Camera.ZoomFactor
	// 		},
	// 		Ship: gameManager.PlayerShip,
	// 	}
	// 	err = shipDoodad.Setup()
	// 	if err != nil {
	// 		slog.Error("Failed to setup ship doodad", "error", err)
	// 		loadingDoodad.SetMessage("Failed to load world map")
	// 		return
	// 	}
	// 	p.Children.Add(shipDoodad)

	// 	// --- Route Doodad

	// 	routeDoodad := &RouteDoodad{
	// 		SpaceTranslator: p.SpaceTranslator,
	// 		Gesturer:        p.Gesturer,
	// 		Ship:            gameManager.PlayerShip,
	// 	}
	// 	err = routeDoodad.Setup()
	// 	if err != nil {
	// 		slog.Error("Failed to setup route doodad", "error", err)
	// 		loadingDoodad.SetMessage("Failed to load world map")
	// 		return
	// 	}
	// 	p.Children.Add(routeDoodad)

	// 	// --- Cursor Doodad

	// 	cursorDoodad := &CursorDoodad{
	// 		Gesturer:        p.Gesturer,
	// 		SpaceTranslator: p.SpaceTranslator,
	// 	}
	// 	err = cursorDoodad.Setup()
	// 	if err != nil {
	// 		slog.Error("Failed to setup cursor doodad", "error", err)
	// 		loadingDoodad.SetMessage("Failed to load world map")
	// 		return
	// 	}
	// 	p.Children.Add(cursorDoodad)

	// 	timeControlDoodad := &TimeControlDoodad{
	// 		Gesturer: p.Gesturer,
	// 		Manager:  gameManager,
	// 		Position: func() doodad.Position {
	// 			return doodad.Position{
	// 				X: p.Width - 200,
	// 				Y: 0,
	// 			}
	// 		},
	// 	}
	// 	err = timeControlDoodad.Setup()
	// 	if err != nil {
	// 		slog.Error("Failed to setup time control doodad", "error", err)
	// 		loadingDoodad.SetMessage("Failed to load world map")
	// 		return
	// 	}
	// 	p.Children.Add(timeControlDoodad)

	// 	tileInfoDoodad := &TileInformationDoodad{
	// 		SpaceTranslator: p.SpaceTranslator,
	// 		Gesturer:        p.Gesturer,

	// 		HideTileCursor: func() {
	// 			cursorDoodad.Hidden = true
	// 		},
	// 		ShowTileCursor: func() {
	// 			cursorDoodad.Hidden = false
	// 		},
	// 	}
	// 	err = tileInfoDoodad.Setup()
	// 	if err != nil {
	// 		slog.Error("Failed to setup tile info doodad", "error", err)
	// 		loadingDoodad.SetMessage("Failed to load world map")
	// 		return
	// 	}
	// 	p.Children.Add(tileInfoDoodad)

	// 	informationTabsDoodad := &InformationTabsDoodad{
	// 		Gesturer: p.Gesturer,
	// 		Manager:  gameManager,
	// 		Position: doodad.ZeroZero,
	// 	}
	// 	err = informationTabsDoodad.Setup()
	// 	if err != nil {
	// 		slog.Error("Failed to setup information tabs doodad", "error", err)
	// 		loadingDoodad.SetMessage("Failed to load world map")
	// 		return
	// 	}
	// 	p.Children.Add(informationTabsDoodad)

	// }()

	return p
}

func (w *WorldMapPage) SetWidthAndHeight(width, height int) {
	w.Box.SetDimensions(width, height)
	w.Box.Recalculate()
}
