package start

import (
	"log/slog"
	"math"
	"math/rand"

	"github.com/jhuggett/sea/constructs/items"
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/data/crew"
	"github.com/jhuggett/sea/data/inventory"
	"github.com/jhuggett/sea/data/ship"
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/outbound"
	"github.com/jhuggett/sea/timeline"
)

type Connection struct {
	Receiver *outbound.Receiver
	GameCtx  *game_context.GameContext
}

func (c *Connection) Context() *game_context.GameContext {
	return c.GameCtx
}

func (c *Connection) Sender() *outbound.Sender {
	return outbound.NewSender(c.GameCtx, *c.Receiver)
}

func Game(conn *Connection) func() {
	Timeline := conn.Context().Timeline

	bitsToCleanUp := []func(){}

	// Timeline.OnTicksPerCycleChangedDo(func(data. timeline.TicksPerCycleChangedEvent) {
	// 	conn.Sender().TimeChanged(data..CurrentTick, data..NewTicksPerCycle)
	// })

	bitsToCleanUp = append(bitsToCleanUp, Timeline.Do(func() timeline.Tick {
		conn.Sender().TimeChanged()

		return timeline.Day
	}, timeline.Day))

	Timeline.Start()

	s, err := conn.Context().Ship()
	if err != nil {
		slog.Error("Error getting ship", "err", err)
		panic(err)
	}

	bitsToCleanUp = append(bitsToCleanUp,
		Timeline.Do(func() timeline.Tick {
			slog.Info("A day has passed")

			// Pay wages (probably should be payed in a different way later)

			crew, err := s.Crew()
			if err != nil {
				slog.Error("Error getting crew", "err", err)
				return timeline.Day
			}

			slog.Info("Got crew", "crew", crew)

			// inventory, err := s.Inventory().Fetch()
			// if err != nil {
			// 	slog.Error("Error fetching inventory", "err", err)
			// 	return timeline.Day
			// }

			// err = inventory.RemoveItem(data.Item{
			// 	Name:   string(constructs.PieceOfEight),
			// 	Amount: float32(crew.Persistent.Wage),
			// })

			// if err != nil {
			// 	slog.Error("Error removing item", "err", err)
			// }

			// Hand out rations

			// rations, err := inventory.Rations()
			// if err != nil {
			// 	slog.Error("Error getting rations", "err", err)
			// 	return timeline.Day
			// }

			// amountOwed := float64(crew.Persistent.Size) * crew.Persistent.Rations
			// slog.Info("Amount owed", "amountOwed", amountOwed, "rations", rations)

			// for _, ration := range rations {
			// 	slog.Info("Ration", "ration", ration, "amountOwed", amountOwed)
			// 	if amountOwed <= 0 {
			// 		break
			// 	}

			// 	if float64(ration.Persistent.Amount) >= amountOwed { // remove some

			// 		slog.Info("Removing ration", "ration", ration, "amountOwed", float32(amountOwed), "rationAmount", ration.Persistent.Amount)

			// 		err := inventory.RemoveItem(data.Item{
			// 			Name:   ration.Persistent.Name,
			// 			Amount: float32(amountOwed),
			// 		})
			// 		amountOwed = 0

			// 		if err != nil {
			// 			slog.Error("Error removing ration", "err", err)
			// 		}
			// 	} else { // remove all
			// 		amountOwed -= float64(ration.Persistent.Amount)

			// 		err := inventory.RemoveItem(data.Item{
			// 			Name:   ration.Persistent.Name,
			// 			Amount: float32(ration.Persistent.Amount),
			// 		})

			// 		if err != nil {
			// 			slog.Error("Error removing ration", "err", err)
			// 		}
			// 	}
			// }

			// if amountOwed > 0 {
			// 	crew.Persistent.Morale = math.Max(0, crew.Persistent.Morale-.1)
			// 	err := crew.Save()
			// 	if err != nil {
			// 		slog.Error("Error saving crew", "err", err)
			// 	}
			// } else {
			// 	crew.Persistent.Morale = math.Min(1, crew.Persistent.Morale+.025)
			// 	err := crew.Save()
			// 	if err != nil {
			// 		slog.Error("Error saving crew", "err", err)
			// 	}
			// }

			// Random ship damage

			s, err = s.Fetch()
			if err != nil {
				slog.Error("Error fetching ship", "err", err)
				return timeline.Day
			}

			amount := rand.Float64() * 0.01
			s.Persistent.StateOfRepair = math.Max(0, s.Persistent.StateOfRepair-amount)
			err = s.Save()
			if err != nil {
				slog.Error("Error saving ship", "err", err)
			}

			err = conn.Sender().CrewInformation()
			if err != nil {
				slog.Error("Error sending crew information", "err", err)
			}

			err = conn.Sender().ShipChanged(s.Persistent.ID)
			if err != nil {
				slog.Error("Error sending ship changed", "err", err)
			}

			return timeline.Day
		}, timeline.Day))

	// Timeline.Do(func() uint64 {
	// 	conn.Sender().TimeChanged(Timeline.CurrentTick(), Timeline.TicksPerCycle())

	// 	return 1
	// }, 1)

	bitsToCleanUp = append(bitsToCleanUp, s.Inventory().OnChangedDo(func(data inventory.OnChangedEvent) {
		conn.Sender().ShipInventoryChanged(s.Persistent.ID, data.Inventory)
	}))

	playerCrew, err := s.Crew()
	if err != nil {
		slog.Error("Error getting crew", "err", err)
		panic(err)
	}

	bitsToCleanUp = append(bitsToCleanUp, playerCrew.OnChangedDo(func(data crew.OnChangedEvent) {
		conn.Sender().CrewInformation()
	}))

	var stopFishing func()

	bitsToCleanUp = append(bitsToCleanUp, s.OnAnchorRaisedDo(func(data ship.AnchorRaisedEvent) {
		if stopFishing != nil {
			stopFishing()
		}
	}))

	bitsToCleanUp = append(bitsToCleanUp, s.OnAnchorLoweredDo(func(e ship.AnchorLoweredEvent) {
		if e.Location == ship.AnchorLoweredLocationOpenSea {
			stopFishing = Timeline.Do(func() timeline.Tick {
				inventory, err := s.Inventory().Fetch()
				if err != nil {
					slog.Error("Error fetching inventory", "err", err)
					return timeline.Day

				}

				perishableDate := uint(Timeline.CurrentTick()) + uint(timeline.Day*3)

				crew, err := s.Crew()
				if err != nil {
					slog.Error("Error getting crew", "err", err)
					return timeline.Day
				}

				crewSize, err := crew.Size()
				if err != nil {
					return timeline.Day
				}

				err = inventory.AddItem(data.Item{
					Name:           string(items.Fish),
					Amount:         float32(crewSize) * 2,
					PerishDate:     &perishableDate, // this doesn't really do anything yet. Need to figure out how this stuff merges
					MarkedAsRation: true,            // this should be controllable by the player in future
				})

				if err != nil {
					slog.Error("Error adding item", "err", err)
				}

				return timeline.Day
			}, timeline.Day)
		}
	}))

	// TODO: optimize this

	// ports, err := port.All(conn.gameCtx.GameMapID())
	// if err != nil {
	// 	slog.Error("Error getting ports", "err", err)
	// 	return
	// }

	// for _, p := range ports {
	// 	slog.Info("Port", "id", p.Persistent.ID)
	// 	producers, err := producer.All(p.Persistent.ID)
	// 	if err != nil {
	// 		slog.Error("Error getting producers", "err", err)
	// 		return
	// 	}

	// 	productionRate := timeline.Day

	// 	for _, pr := range producers {
	// 		Timeline.Do(func() uint64 {
	// 			inventory, err := p.Inventory().Fetch()
	// 			if err != nil {
	// 				slog.Error("Error fetching inventory", "err", err)
	// 				return productionRate
	// 			}

	// 			for _, product := range pr.Products() {
	// 				err = inventory.AddItem(data.Item{
	// 					Name:   product,
	// 					Amount: 1,
	// 				})
	// 				if err != nil {
	// 					slog.Error("Error adding item", "err", err)
	// 				}
	// 			}

	// 			return productionRate
	// 		}, productionRate)
	// 	}
	// }

	return func() {
		for _, fn := range bitsToCleanUp {
			fn()
		}
	}
}
