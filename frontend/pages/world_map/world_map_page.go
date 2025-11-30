package world_map

import (
	design_library "design-library"
	"design-library/doodad"
	"design-library/position/box"
	"design-library/reaction"
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/frontend/pages/port_map"
	"github.com/jhuggett/frontend/pages/world_map/bottom_bar"
	"github.com/jhuggett/frontend/pages/world_map/camera"
	"github.com/jhuggett/frontend/pages/world_map/pause_menu"
	"github.com/jhuggett/frontend/utils/space_translator"
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/outbound"
	"golang.org/x/image/math/f64"
)

type WorldMapPage struct {
	App *design_library.App

	Camera          *camera.Camera
	SpaceTranslator space_translator.SpaceTranslator

	GameSnapshot *game_context.Snapshot
	GameManager  *game.Manager

	doodad.Default

	dockConfirmation *DockConfirmation

	portMapPage *port_map.PortMapPage
}

func New(snapshot *game_context.Snapshot, app *design_library.App) *WorldMapPage {
	p := &WorldMapPage{
		GameSnapshot: snapshot,
		App:          app,
	}

	return p
}

func (w *WorldMapPage) LoadRequiredData() error {
	slog.Info("Loading World Map Page Required Data")

	w.Camera = &camera.Camera{
		ViewPort:   f64.Vec2{float64(w.Box.Width()), float64(w.Box.Height())},
		TileSize:   32,
		Position:   f64.Vec2{0, 0},
		ZoomFactor: 1,
	}

	w.SpaceTranslator = space_translator.New(w.Camera, float64(w.Camera.TileSize))

	gameManager := game.NewManager(w.GameSnapshot)
	gameManager.SetTileSize(int(w.Camera.TileSize))

	err := gameManager.Start()
	if err != nil {
		return fmt.Errorf("failed to start game manager: %w", err)
	}

	w.GameManager = gameManager

	worldMap := doodad.Stateful(w, "world-map", func() *WorldMapDoodad {
		return NewWorldMapDoodad(
			w.GameManager.WorldMap,
			w.SpaceTranslator,
		)
	})

	slog.Info("Loading World Map Data, worldMap.Load()")
	worldMap.Load()

	return nil
}

func (w *WorldMapPage) Setup() {
	// worldMapDoodad := NewWorldMapDoodad(
	// 	w.GameManager.WorldMap,
	// 	w.SpaceTranslator,
	// )
	// w.AddChild(worldMapDoodad)

	// Should be pre-loaded in LoadRequiredData
	w.AddStatefulChild("world-map", nil)

	// w.AddChild(NewRouteDoodad(
	// 	w.SpaceTranslator,
	// 	w.GameManager.PlayerShip,
	// ))

	w.AddStatefulChild("route", func() doodad.Doodad {
		return NewRouteDoodad(
			w.SpaceTranslator,
			w.GameManager.PlayerShip,
		)
	})

	// w.AddChild(NewPlayerShipDoodad(
	// 	w.GameManager,
	// 	w.Camera,
	// ))

	w.AddStatefulChild("player-ship", func() doodad.Doodad {
		return NewPlayerShipDoodad(
			w.GameManager,
			w.Camera,
		)
	})

	// w.AddChild(NewTileInformationDoodad(
	// 	w.SpaceTranslator,
	// ))

	w.AddStatefulChild("tile-information", func() doodad.Doodad {
		return NewTileInformationDoodad(
			w.SpaceTranslator,
		)
	})

	// w.AddChild(NewCursorDoodad(
	// 	w.SpaceTranslator,
	// ))

	w.AddStatefulChild("cursor", func() doodad.Doodad {
		return NewCursorDoodad(
			w.SpaceTranslator,
		)
	})

	timeControlDoodad := doodad.Stateful(w, "time-control", func() *TimeControlDoodad {
		return NewTimeControlDoodad(
			w.GameManager,
		)
	})

	w.DoOnTeardown(w.GameManager.OnShipDockedCallback.Register(func(req outbound.ShipDockedReq) error {
		if !req.Undocked {
			w.dockConfirmation.Port = game.PortFromOutboundData(w.GameManager, req.Port)
			w.dockConfirmation.Show()
			doodad.ReSetup(w.dockConfirmation)
		}
		return nil
	}))
	w.AddChild(timeControlDoodad)
	// timeControlDoodad.Layout().Computed(func(b *box.Box) {
	// 	b.Copy(w.Box)
	// })

	bottomBar := doodad.Stateful(w, "bottom-bar", func() *bottom_bar.BottomBar {
		return bottom_bar.NewBottomBar(w.GameManager)
	})
	w.AddChild(bottomBar)

	routeInfoDoodad := doodad.Stateful(w, "route-info", func() *RouteInformationDoodad {
		return NewRouteInformationDoodad(
			w.GameManager.PlayerShip,
			func(b *box.Box) {
				b.MoveAbove(bottomBar.Box).MoveDown(100).CenterHorizontallyWithin(w.Box)
			},
		)
	})
	w.AddChild(routeInfoDoodad)

	pauseMenu := doodad.Stateful(w, "pause-menu", func() *pause_menu.PauseMenu {
		return pause_menu.NewPauseMenu(w.App)
	})
	w.AddChild(pauseMenu)
	// pauseMenu.Layout().Computed(func(b *box.Box) {
	// 	b.Copy(w.Box)
	// })
	pauseMenu.Hide()

	w.dockConfirmation = doodad.Stateful(w, "dock-confirmation", func() *DockConfirmation {
		return NewDockConfirmation(
			w.GameManager,
			nil,
			func(port *game.Port) {
				// On dock accepted
				w.dockConfirmation.Hide()

				w.GameManager.ControlTime(inbound.ControlTimeReq{
					Pause: true,
				})

				// portMap := port_map.New(
				// 	w.GameManager,
				// 	w.App,
				// 	port,
				// )

				// w.portMapPage.Show()

				portMapPage := port_map.New(
					w.GameManager,
					w.App,
					port,
				)

				w.AddChild(portMapPage)

				portMapPage.GoBack = func() {
					w.Children().Remove(portMapPage)
				}

				// portMapPage.Setup()

				doodad.Setup(portMapPage)

				// portMap.GoBack = func() {
				// 	w.App.Pop()
				// }

				// w.App.Push(portMap)
			},
			func() {
				// On dock declined
				w.dockConfirmation.Hide()
			},
		)
	})
	w.AddChild(w.dockConfirmation)
	// w.dockConfirmation.Layout().Computed(func(b *box.Box) {
	// 	b.Copy(w.Box)
	// })
	w.dockConfirmation.Hide()

	// w.portMapPage = port_map.New(
	// 	w.GameManager,
	// 	w.App,
	// 	nil,
	// )
	// w.AddChild(w.portMapPage)

	w.Reactions().Add(
		reaction.NewKeyDownReaction(reaction.SpecificKeyDown(ebiten.KeySpace),
			func(event *reaction.KeyDownEvent) {
				// Center the camera on the player's ship when space bar is pressed
				w.CenterCameraOnPlayerShip()
			},
		),
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

				// Scale zoom by event.YOffset (so trackpads produce smaller changes)
				zoomDelta := event.YOffset * 0.1 // reduce sensitivity; tune this constant
				w.Camera.ZoomFactor += zoomDelta

				// Clamp zoom
				if w.Camera.ZoomFactor < 0.1 {
					w.Camera.ZoomFactor = 0.1
				} else if w.Camera.ZoomFactor > 10 {
					w.Camera.ZoomFactor = 10
				}

				// Adjust camera position so the world point under the mouse stays fixed
				newZoom := w.Camera.ZoomFactor
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

	w.Box.Recalculate()
}
