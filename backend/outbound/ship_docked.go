package outbound

import (
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/constructs/items"
	"github.com/jhuggett/sea/data/continent"
	"github.com/jhuggett/sea/data/point"
	"github.com/jhuggett/sea/data/population"
	"github.com/jhuggett/sea/data/port"
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

	ItemValuation map[string]float64 `json:"item_valuation,omitempty"`
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

		continentData, err := port.Continent()
		if err != nil {
			slog.Error("Failed to find continent", "port", port, "error", err)
			return fmt.Errorf("failed to find continent: %w", err)
		}

		populationData, err := continent.Using(*continentData).Population()
		if err != nil {
			slog.Error("Failed to find population", "continent", continentData, "error", err)
			return fmt.Errorf("failed to find population: %w", err)
		}

		populationModel := population.Using(populationData)

		itemSummaries := make([]Item, 0, len(port.Persistent.Inventory.Items))

		valueMap := make(map[string]float64)
		minValue := 0.0
		maxValue := 0.0

		for _, it := range port.Persistent.Inventory.Items {

			itemValue, err := populationModel.Value(items.LookupItem(it.Name))
			if err != nil {
				slog.Error("Failed to get item value", "item", it.Name, "error", err)
				return fmt.Errorf("failed to get item value: %w", err)
			}

			valueMap[it.Name] = itemValue
			if itemValue < minValue {
				minValue = itemValue
			}

			if itemValue > maxValue {
				maxValue = itemValue
			}

			itemSummaries = append(itemSummaries, Item{
				ID:     it.ID,
				Name:   it.Name,
				Amount: it.Amount,
			})
		}

		playerShipModel, err := s.gameContext.Ship()
		if err != nil {
			slog.Error("Failed to get player ship", "error", err)
			return fmt.Errorf("failed to get player ship: %w", err)
		}

		for _, it := range playerShipModel.Persistent.Inventory.Items {
			itemValue, err := populationModel.Value(items.LookupItem(it.Name))
			if err != nil {
				slog.Error("Failed to get item value", "item", it.Name, "error", err)
				return fmt.Errorf("failed to get item value: %w", err)
			}

			valueMap[it.Name] = itemValue
			if itemValue < minValue {
				minValue = itemValue
			}
			if itemValue > maxValue {
				maxValue = itemValue
			}
		}

		for i, v := range valueMap {
			valueMap[i] = (v - minValue) / (maxValue - minValue)
		}

		slog.Info("ShipDocked", "items", itemSummaries, "port", port, "point", point)

		_, err = s.rpc.Send("ShipDocked", ShipDockedReq{
			ShipID:   shipId,
			Location: location,
			Undocked: undocked,
			Port: Port{
				ID: port.Persistent.ID,
				Inventory: Inventory{
					ID:    port.Persistent.InventoryID,
					Items: itemSummaries,
				},
				ItemValuation: valueMap,
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
