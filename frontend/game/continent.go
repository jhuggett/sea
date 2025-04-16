package game

import "github.com/jhuggett/sea/inbound"

type Continent struct {
	Manager *Manager

	Ports []*Port

	rawData inbound.Continent

	LargestX int
	LargestY int

	SmallestX int
	SmallestY int
}

func NewContinent(manager *Manager, rawData inbound.Continent) *Continent {
	cont := &Continent{
		rawData: rawData,
		Manager: manager,
	}

	for _, p := range rawData.Points {
		if p.X > cont.LargestX {
			cont.LargestX = p.X
		}
		if p.Y > cont.LargestY {
			cont.LargestY = p.Y
		}
		if p.X < cont.SmallestX {
			cont.SmallestX = p.X
		}
		if p.Y < cont.SmallestY {
			cont.SmallestY = p.Y
		}
	}

	return cont
}

func (c *Continent) Points() []inbound.Point {
	return c.rawData.Points
}

func (c *Continent) TileSize() int {
	return c.Manager.TileSize()
}
