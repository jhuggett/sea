package inbound

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/inventory"
)

type Item struct {
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
}

type TradeAction struct {
	Item          Item `json:"item"`
	FromInventory uint `json:"from"`
	ToInventory   uint `json:"to"`
}

type TradeReq struct {
	Actions []TradeAction `json:"actions"`
}

type TradeResp struct {
}

func Trade(conn Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r TradeReq
		if err := json.Unmarshal(req, &r); err != nil {
			return nil, err
		}

		for _, a := range r.Actions {
			slog.Info("Trade", "from", a.FromInventory, "to", a.ToInventory, "item", a.Item)

			from, err := inventory.Fetch(a.FromInventory)
			if err != nil {
				return nil, err
			}

			to, err := inventory.Fetch(a.ToInventory)
			if err != nil {
				return nil, err
			}

			err = from.RemoveItem(models.Item{
				Name:   a.Item.Name,
				Amount: a.Item.Amount,
			})

			if err != nil {
				slog.Warn("Failed to remove item from inventory", "items", from.Items(), "err", err, "inventory", from, "item", a.Item)
				return nil, fmt.Errorf("failed to remove item from inventory: %w", err)
			}

			err = to.AddItem(models.Item{
				Name:   a.Item.Name,
				Amount: a.Item.Amount,
			})

			if err != nil {
				slog.Warn("Failed to add item to inventory", "items", to.Items(), "err", err, "inventory", to, "item", a.Item)
				return nil, fmt.Errorf("failed to add item to inventory: %w", err)
			}

		}

		resp := TradeResp{}

		return resp, nil
	}
}
