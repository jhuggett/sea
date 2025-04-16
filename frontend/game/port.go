package game

import "github.com/jhuggett/sea/inbound"

type Port struct {
	Manager *Manager

	RawData inbound.Port
}

func (p *Port) Location() inbound.Point {
	return inbound.Point{
		X: p.RawData.Point.X,
		Y: p.RawData.Point.Y,
	}
}
