package port_map

import (
	design_library "design-library"
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/reaction"
	"log/slog"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/frontend/pages/world_map/camera"
	"github.com/jhuggett/frontend/pages/world_map/pause_menu"
	"github.com/jhuggett/frontend/utils/space_translator"
	"golang.org/x/image/math/f64"
)

type PortMapPage struct {
	App *design_library.App

	GoBack func()

	Camera          *camera.Camera
	SpaceTranslator space_translator.SpaceTranslator

	GameManager *game.Manager

	Port      *game.Port
	Buildings []*game.Building

	PointToBuilding map[f64.Vec2]*game.Building

	doodad.Default

	loaded bool
}

func New(manager *game.Manager, app *design_library.App, port *game.Port) *PortMapPage {
	p := &PortMapPage{
		GameManager: manager,
		App:         app,
		Port:        port,
		Default: doodad.Default{
			Box: nil,
		},
	}

	buildings, err := p.Port.GetBuildings()
	if err != nil {
		panic(err)
	}
	p.Buildings = buildings

	return p
}

func (p *PortMapPage) Setup() {

	tileSize := 128

	p.Camera = &camera.Camera{
		ViewPort:   f64.Vec2{float64(p.Box.Width()), float64(p.Box.Height())},
		TileSize:   float64(tileSize),
		Position:   f64.Vec2{0, 0},
		ZoomFactor: 1,
	}

	p.SpaceTranslator = space_translator.New(p.Camera, float64(tileSize))

	// For lookups on mouse event
	p.PointToBuilding = make(map[f64.Vec2]*game.Building)
	for _, building := range p.Buildings {
		p.PointToBuilding[f64.Vec2{
			float64(building.X),
			float64(building.Y),
		}] = building
	}

	pauseMenu := pause_menu.NewPauseMenu()
	p.AddChild(pauseMenu)
	pauseMenu.Hide()

	backbutton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			p.GoBack()
		},
		Config: label.Config{
			Message: "Back to World Map",
		},
	})

	p.AddChild(backbutton)

	portMap := NewPortMapDoodad(
		p.Camera,
		p.SpaceTranslator,
		p.Port,
		p.Buildings,
		tileSize,
	)
	p.AddChild(portMap)

	buildingMenu := NewBuildingMenu()
	p.AddChild(buildingMenu)
	buildingMenu.Hide()

	p.Reactions().Add(
		reaction.NewMouseDragReaction(
			func(event *reaction.OnMouseDragEvent) bool {
				return true
			},
			func(event *reaction.OnMouseDragEvent) {
				p.Camera.Position[0] += float64(event.StartX-event.X) / p.Camera.ZoomFactor
				p.Camera.Position[1] += float64(event.StartY-event.Y) / p.Camera.ZoomFactor
			},
		),
		reaction.NewMouseWheelReaction(
			func(event *reaction.MouseWheelEvent) bool {
				return true
			},
			func(event *reaction.MouseWheelEvent) {
				mouseX, mouseY := p.Gesturer().CurrentMouseLocation()

				// Convert mouse position to world coordinates before zoom
				worldX := (float64(mouseX) / p.Camera.ZoomFactor) + p.Camera.Position[0]
				worldY := (float64(mouseY) / p.Camera.ZoomFactor) + p.Camera.Position[1]

				// Scale zoom by event.YOffset (so trackpads produce smaller changes)
				zoomDelta := event.YOffset * 0.1 // reduce sensitivity; tune this constant
				p.Camera.ZoomFactor += zoomDelta

				// Clamp zoom
				if p.Camera.ZoomFactor < 0.1 {
					p.Camera.ZoomFactor = 0.1
				} else if p.Camera.ZoomFactor > 10 {
					p.Camera.ZoomFactor = 10
				}

				// Adjust camera position so the world point under the mouse stays fixed
				newZoom := p.Camera.ZoomFactor
				p.Camera.Position[0] = worldX - float64(mouseX)/newZoom
				p.Camera.Position[1] = worldY - float64(mouseY)/newZoom
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
		reaction.NewMouseUpReaction(func(event *reaction.MouseUpEvent) bool {
			return true
		}, func(event *reaction.MouseUpEvent) {
			mouseX, mouseY := p.Gesturer().CurrentMouseLocation()

			x, y := p.SpaceTranslator.FromWorldToData(p.SpaceTranslator.FromScreenToWorld(float64(mouseX), float64(mouseY)))
			x = math.Floor(x)
			y = math.Floor(y)
			building := p.PointToBuilding[f64.Vec2{x, y}]
			if building != nil {
				buildingMenu.Building = building
				doodad.ReSetup(buildingMenu)
				buildingMenu.Show()
			}

			slog.Info("MouseUp", "mouseX", mouseX, "mouseY", mouseY, "x", x, "y", y, "building", building)
		}),
	)

	p.Children().Setup()
}
