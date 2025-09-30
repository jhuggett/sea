package world_map

import (
	"log/slog"
)

// CenterCameraOnPlayerShip centers the camera view on the player's ship
func (w *WorldMapPage) CenterCameraOnPlayerShip() {
	// Check if we have a valid player ship and game manager
	if w.GameManager == nil || w.GameManager.PlayerShip == nil {
		slog.Error("Cannot center camera: game manager or player ship is nil")
		return
	}

	slog.Info("Centering camera on player ship")

	// Get the ship's position in data coordinates (tile coordinates)
	shipDataX := float64(w.GameManager.PlayerShip.Location().X)
	shipDataY := float64(w.GameManager.PlayerShip.Location().Y)

	// Convert to world coordinates (pixel coordinates)
	shipWorldX, shipWorldY := w.SpaceTranslator.FromDataToWorld(shipDataX, shipDataY)

	// Center the camera on the ship (accounting for viewport size and zoom)
	w.Camera.Position[0] = shipWorldX - (w.Camera.ViewPort[0]/2)/w.Camera.ZoomFactor
	w.Camera.Position[1] = shipWorldY - (w.Camera.ViewPort[1]/2)/w.Camera.ZoomFactor

	slog.Info("Camera centered on ship",
		"shipX", w.GameManager.PlayerShip.Location().X,
		"shipY", w.GameManager.PlayerShip.Location().Y,
		"cameraX", w.Camera.Position[0],
		"cameraY", w.Camera.Position[1])
}
