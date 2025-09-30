package world_map

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"fmt"
	"log/slog"

	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
)

func NewRouteInformationDoodad(
	ship *game.Ship,
	positioner func(*box.Box),
) *RouteInformationDoodad {
	d := &RouteInformationDoodad{
		Ship:      ship,
		Postioner: positioner,
	}

	ship.Manager.RouteEventCallback.Add(func(data game.RouteEventCallbackData) error {
		doodad.ReSetup(d)
		return nil
	})

	return d
}

type RouteInformationDoodad struct {
	SpaceTranslator SpaceTranslator

	Postioner func(*box.Box)

	Ship *game.Ship

	doodad.Default
}

func (w *RouteInformationDoodad) Setup() {
	panel := stack.New(stack.Config{
		Layout: box.Computed(func(b *box.Box) {
			// boundingBox := box.Bounding(panelChildren.Boxes())
			// w.Postioner(b.CopyDimensionsOf(boundingBox))

			b.SetPosition(0, 0)
		}),
		FitContents: true,
		Padding: stack.Padding{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
		SpaceBetween:    10,
		Type:            stack.Vertical, // Changed to vertical for better information display
		BackgroundColor: colors.Panel,
	})

	w.AddChild(panel)

	// var doodads []doodad.Doodad

	// Title label is always shown
	titleLabel := label.New(label.Config{
		Message: "Route Information",
	})
	// doodads = append(doodads, titleLabel)

	panel.AddChild(titleLabel)

	// Always ensure we show the information when a route exists
	if w.Ship.HasRoute() {
		// Show Set Sail button for inactive routes
		if !w.Ship.IsRouteActive() {
			setSailButton := button.New(button.Config{
				OnClick: func(*button.Button) {
					slog.Info("Set Sail button clicked")
					_, err := w.Ship.SetSail()
					if err != nil {
						slog.Error("Failed to set sail", "error", err)
						return
					}
				},
				Config: label.Config{
					Message: "Set Sail",
				},
			})
			panel.AddChild(setSailButton)
		} else {
			// For active routes, show detailed information and controls
			var progressMessage, timeLeftMessage, statusMessage, destinationMessage string

			// Prepare information messages if we have movement data
			if w.Ship.Route.ShipMovedReq != nil {
				routeInfo := w.Ship.Route.ShipMovedReq.RouteInfo

				// Progress information
				progress := float64(routeInfo.TotalTilesMoved) / float64(routeInfo.TilesInRoute) * 100
				progressMessage = fmt.Sprintf("Progress: %.0f%% (%d/%d tiles)",
					progress, routeInfo.TotalTilesMoved, routeInfo.TilesInRoute)

				// ETA information
				if routeInfo.EstimatedTimeLeft > 0 {
					timeLeftMessage = fmt.Sprintf("ETA: %d days", int(routeInfo.EstimatedTimeLeft))
				} else {
					timeLeftMessage = "Arriving soon"
				}

				// Status information
				status := "Active"
				if routeInfo.IsPaused {
					status = "Paused"
				} else if routeInfo.IsCancelled {
					status = "Cancelled"
				} else if routeInfo.ReachedDestination {
					status = "Arrived"
				}
				statusMessage = fmt.Sprintf("Status: %s", status)

				// Destination information
				if routeInfo.HeadedToPort {
					destinationMessage = "Destination: Port"
				} else {
					destinationMessage = "Destination: Open sea"
				}
			} else {
				// Default messages if no movement data yet
				progressMessage = "Progress: --"
				timeLeftMessage = "ETA: --"
				statusMessage = "Status: Active"
				destinationMessage = "Destination: --"
			}

			// Create and add labels
			progressLabel := label.New(label.Config{Message: progressMessage})
			timeLeftLabel := label.New(label.Config{Message: timeLeftMessage})
			statusLabel := label.New(label.Config{Message: statusMessage})
			destLabel := label.New(label.Config{Message: destinationMessage})

			panel.AddChild(progressLabel)
			panel.AddChild(timeLeftLabel)
			panel.AddChild(statusLabel)
			panel.AddChild(destLabel)

			// Add route control buttons
			cancelButton := button.New(button.Config{
				OnClick: func(*button.Button) {

					w.Ship.CancelCurrentRoute()

				},
				Config: label.Config{
					Message: "Cancel Route",
				},
			})
			panel.AddChild(cancelButton)

			// Create pause button with appropriate label based on state
			pauseButtonLabel := "Pause Journey"
			if w.Ship.Route.ShipMovedReq != nil && w.Ship.Route.ShipMovedReq.RouteInfo.IsPaused {
				pauseButtonLabel = "Resume Journey"
			}

			pauseButton := button.New(button.Config{
				OnClick: func(*button.Button) {

					if w.Ship.Route.ShipMovedReq != nil && w.Ship.Route.ShipMovedReq.RouteInfo.IsPaused {
						w.Ship.ResumeCurrentRoute()
					} else {
						w.Ship.PauseCurrentRoute()
					}

				},
				Config: label.Config{
					Message: pauseButtonLabel,
				},
			})
			panel.AddChild(pauseButton)
		}
	}

	w.Children().Setup()

	w.Layout().Recalculate()
}
