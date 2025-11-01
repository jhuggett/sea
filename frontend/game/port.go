package game

import (
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/outbound"
)

type RawPortData struct {
	ID uint `json:"id"`

	Point data.Point `json:"point"`

	Name string `json:"name"`

	ContinentID uint `json:"continent_id"`
}

type Port struct {
	Manager *Manager

	RawData RawPortData
}

func PortFromInboundData(m *Manager, d inbound.Port) *Port {
	return &Port{
		Manager: m,
		RawData: RawPortData{
			ID: d.ID,
			Point: data.Point{
				X: d.Point.X,
				Y: d.Point.Y,
			},
			Name:        d.Name,
			ContinentID: d.ContinentID,
		},
	}
}

func PortFromOutboundData(m *Manager, d outbound.Port) *Port {
	return &Port{
		Manager: m,
		RawData: RawPortData{
			ID: d.ID,

			Point: data.Point{
				X: d.Point.X,
				Y: d.Point.Y,
			},

			Name:        d.Name,
			ContinentID: d.ContinentID,
		},
	}
}

func (p *Port) Location() inbound.Point {
	return inbound.Point{
		X: p.RawData.Point.X,
		Y: p.RawData.Point.Y,
	}
}

func (p *Port) GetBuildings() ([]*Building, error) {
	resp, err := p.Manager.GetBuildingsByPortID(p.RawData.ID)
	if err != nil {
		return nil, err
	}

	var result []*Building
	for _, b := range resp.Buildings {
		result = append(result, &Building{
			Name: b.Name,
			Type: b.Type,
			X:    b.X,
			Y:    b.Y,
			Port: p,
		})
	}

	return result, nil
}
