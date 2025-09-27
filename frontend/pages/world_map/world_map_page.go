package world_map

import (
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/reaction"
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/frontend/pages/world_map/bottom_bar"
	"github.com/jhuggett/frontend/pages/world_map/camera"
	"github.com/jhuggett/frontend/pages/world_map/pause_menu"
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
	Camera          *camera.Camera
	SpaceTranslator SpaceTranslator

	GameSnapshot *game_context.Snapshot
	GameManager  *game.Manager

	doodad.Default

	loaded bool
}

func (w *WorldMapPage) SetupLoadingScreen() {
	loadingMessage := label.New(label.Config{
		Message:  "Loading...",
		Layout:   box.Zeroed(),
		FontSize: 24,
	})

	w.AddChild(loadingMessage)

	w.Children().Setup()

	loadingMessage.Box.Computed(func(b *box.Box) {
		b.CenterWithin(w.Box)
	})
}

func (w *WorldMapPage) LoadRequiredData() error {
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

	err := gameManager.Start()
	if err != nil {
		return fmt.Errorf("failed to start game manager: %w", err)
	}

	w.GameManager = gameManager

	w.loaded = true
	return nil
}

func (w *WorldMapPage) SetupMainScreen() {
	worldMapDoodad := NewWorldMapDoodad(
		w.GameManager.WorldMap,
		w.SpaceTranslator,
	)
	w.AddChild(worldMapDoodad)

	// Need to preload the things we need before we get to this method
	// then we can just render

	w.AddChild(NewPlayerShipDoodad(
		w.GameManager,
		w.Camera,
	))

	w.AddChild(NewTileInformationDoodad(
		w.SpaceTranslator,
	))

	w.AddChild(NewCursorDoodad(
		w.SpaceTranslator,
	))

	w.AddChild(NewRouteDoodad(
		w.SpaceTranslator,
		w.GameManager.PlayerShip,
	))

	timeControlDoodad := NewTimeControlDoodad(
		w.GameManager,
	)
	w.AddChild(timeControlDoodad)
	timeControlDoodad.Layout().Computed(func(b *box.Box) {
		b.Copy(w.Box)
	})

	bottomBar := bottom_bar.NewBottomBar(w.GameManager)
	w.AddChild(bottomBar)

	routeInfoDoodad := NewRouteInformationDoodad(
		w.GameManager.PlayerShip,
		func(b *box.Box) {
			b.MoveAbove(bottomBar.Box).MoveDown(100).CenterHorizontallyWithin(w.Box)
		},
	)

	w.AddChild(routeInfoDoodad)
	routeInfoDoodad.Hide()

	pauseMenu := pause_menu.NewPauseMenu()
	w.AddChild(pauseMenu)
	pauseMenu.Layout().Computed(func(b *box.Box) {

		b.Copy(w.Box)
	})
	pauseMenu.Hide()

	w.Reactions().Add(
		reaction.NewMouseDragReaction(
			func(event *reaction.OnMouseDragEvent) bool {
				return true
			},
			func(event *reaction.OnMouseDragEvent) {
				w.Camera.Position[0] += float64(event.StartX-event.X) / w.Camera.ZoomFactor
				w.Camera.Position[1] += float64(event.StartY-event.Y) / w.Camera.ZoomFactor
			},
		),
		reaction.NewMouseWheelReaction(
			func(event *reaction.MouseWheelEvent) bool {
				return true
			},
			func(event *reaction.MouseWheelEvent) {
				mouseX, mouseY := w.Gesturer().CurrentMouseLocation()

				// Convert mouse position to world coordinates before zoom
				worldX := (float64(mouseX) / w.Camera.ZoomFactor) + w.Camera.Position[0]
				worldY := (float64(mouseY) / w.Camera.ZoomFactor) + w.Camera.Position[1]

				// oldZoom := w.Camera.ZoomFactor
				if event.YOffset > 0 {
					w.Camera.ZoomFactor += 0.1
				} else {
					w.Camera.ZoomFactor -= 0.1
					if w.Camera.ZoomFactor < 0.1 {
						w.Camera.ZoomFactor = 0.1
					}
				}
				newZoom := w.Camera.ZoomFactor

				// Adjust camera position so the world point under the mouse stays under the mouse
				w.Camera.Position[0] = worldX - float64(mouseX)/newZoom
				w.Camera.Position[1] = worldY - float64(mouseY)/newZoom
			},
		),
		reaction.NewKeyDownReaction(
			func(event *reaction.KeyDownEvent) bool {
				return true
			},
			func(event *reaction.KeyDownEvent) {
				if event.Key == ebiten.KeyEscape {
					if !pauseMenu.IsVisible() {
						pauseMenu.Show()
					} else {
						pauseMenu.Hide()
					}
				}
			},
		),
	)
	w.Reactions().Register(w.Gesturer(), w.Z())

	w.Children().Setup()
}

func (w *WorldMapPage) Setup() {
	if len(w.Children().Doodads) > 0 {
		w.Children().Clear()
	}

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

	w.Box.Recalculate()
}

func New(snapshot *game_context.Snapshot) *WorldMapPage {
	p := &WorldMapPage{
		GameSnapshot: snapshot,
	}

	return p
}

func (w *WorldMapPage) SetWidthAndHeight(width, height int) {
	w.Box.SetDimensions(width, height)
	w.Box.Recalculate()
}
