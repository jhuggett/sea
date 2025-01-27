package outbound

import "github.com/jhuggett/sea/data/inventory"

type Item struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
}

type ShipInventoryChangedReq struct {
	ShipID uint `json:"ship_id"`

	Inventory Inventory `json:"inventory"`
}

type ShipInventoryChangedResp struct{}

func (s *Sender) ShipInventoryChanged(shipId uint, inventory inventory.Inventory) error {

	items := make([]Item, 0, len(inventory.Items()))

	for _, it := range inventory.Items() {
		items = append(items, Item{
			ID:     it.Persistent.ID,
			Name:   it.Persistent.Name,
			Amount: it.Persistent.Amount,
		})
	}

	_, err := s.rpc.Send("ShipInventoryChanged", ShipInventoryChangedReq{
		ShipID: shipId,
		Inventory: Inventory{
			ID:    inventory.Persistent.ID,
			Items: items,
		},
	})

	return err
}
