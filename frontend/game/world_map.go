package game

import (
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/inbound"
)

type WorldMap struct {
	Manager *Manager

	gameMap *inbound.GetWorldMapResp

	Continents []*Continent
	Ports      []*Port
}

func NewWorldMap(manager *Manager) *WorldMap {
	return &WorldMap{
		Manager: manager,
	}
}

func (w *WorldMap) Load() error {

	gameMap, err := inbound.GetWorldMap(inbound.GetWorldMapReq{}, w.Manager.Snapshot.GameMapID)
	if err != nil {
		return fmt.Errorf("failed to get world map: %w", err)
	}

	portsData, err := inbound.GetPorts(inbound.GetPortsReq{}, w.Manager.Snapshot.GameMapID)
	if err != nil {
		return err
	}
	for _, portData := range portsData.Ports {
		w.Ports = append(w.Ports, &Port{
			Manager: w.Manager,
			RawData: portData,
		})
	}

	w.gameMap = &gameMap

	for _, continent := range gameMap.Continents {

		c := NewContinent(w.Manager, *continent)

		for _, port := range w.Ports {
			slog.Debug("Port", "port", port.RawData.ContinentID, "continent", continent.ID)
			if port.RawData.ContinentID == continent.ID {
				c.Ports = append(c.Ports, port)
			}
		}

		w.Continents = append(w.Continents, c)
	}

	return nil
}
