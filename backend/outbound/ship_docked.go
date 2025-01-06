package outbound

import (
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/models/point"
	"github.com/jhuggett/sea/models/port"
	"github.com/jhuggett/sea/utils/coordination"
)

type ShipDockedReq struct {
	ShipID   uint               `json:"ship_id,omitempty"`
	Location coordination.Point `json:"location,omitempty"`
	Undocked bool               `json:"undocked,omitempty"`

	Port Port `json:"port,omitempty"`
}

type Port struct {
	ID uint `json:"id,omitempty"`

	Inventory Inventory `json:"inventory,omitempty"`
}

type Inventory struct {
	ID uint `json:"id,omitempty"`

	Items []Item `json:"items,omitempty"`
}

type ShipDockedResp struct{}

func (s *Sender) ShipDocked(shipId uint, location coordination.Point, undocked bool) error {
	slog.Info("ShipDocked", "id", shipId, "location", location, "undocked", undocked)

	if !undocked {
		point, err := point.Find(location, s.gameContext.GameMapID())
		if err != nil {
			return fmt.Errorf("failed to find point: %w", err)
		}

		slog.Info("ShipDocked", "point", point)

		port, err := port.Find(*point)
		if err != nil {
			slog.Error("Failed to find port", "point", point, "error", err)
			return fmt.Errorf("failed to find port: %w", err)
		}

		slog.Info("ShipDocked", "port", port)

		items := make([]Item, 0, len(port.Persistent.Inventory.Items))

		for _, it := range port.Persistent.Inventory.Items {
			items = append(items, Item{
				ID:     it.ID,
				Name:   it.Name,
				Amount: it.Amount,
			})
		}
		slog.Info("ShipDocked", "items", items, "port", port, "point", point)

		_, err = s.rpc.Send("ShipDocked", ShipDockedReq{
			ShipID:   shipId,
			Location: location,
			Undocked: undocked,
			Port: Port{
				ID: port.Persistent.ID,
				Inventory: Inventory{
					ID:    port.Persistent.InventoryID,
					Items: items,
				},
			},
		})
		if err != nil {
			return err
		}
	} else {
		slog.Info("rpc.ShipDocked", "shipId", shipId, "location", location, "undocked", undocked)
		_, err := s.rpc.Send("ShipDocked", ShipDockedReq{
			ShipID:   shipId,
			Location: location,
			Undocked: undocked,
		})
		if err != nil {
			return err
		}
		slog.Info("rpc.ShipDocked - done", "shipId", shipId, "location", location, "undocked", undocked)
	}

	return nil
}
