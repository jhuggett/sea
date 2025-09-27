package game

import (
	"fmt"

	"github.com/jhuggett/sea/inbound"
)

type Ship struct {
	Manager *Manager
	RawData inbound.ShipInfo

	route *ShipRoute
}

func (s *Ship) Location() inbound.Point {
	return inbound.Point{
		X: int(s.RawData.X),
		Y: int(s.RawData.Y),
	}
}

type ShipRoute struct {
	Points []inbound.Coordinate
	Active bool
}

func (s *Ship) HasRoute() bool {
	return s.route != nil && len(s.route.Points) > 0
}

func (s *Ship) IsRouteActive() bool {
	return s.route != nil && s.route.Active
}

func (s *Ship) PlotRoute(x, y int) (*ShipRoute, error) {
	resp, err := s.Manager.PlotRoute(x, y)
	if err != nil {
		return nil, fmt.Errorf("failed to plot route: %w", err)
	}

	route := &ShipRoute{
		Points: resp.Coordinates,
	}

	s.route = route

	return route, nil
}

func (s *Ship) SetSail() (*inbound.MoveShipResp, error) {
	if s.route == nil {
		return nil, fmt.Errorf("no route set")
	}

	endTileX := int(s.route.Points[len(s.route.Points)-1].X)
	endTileY := int(s.route.Points[len(s.route.Points)-1].Y)

	resp, err := inbound.MoveShip(inbound.MoveShipReq{
		X: float64(endTileX),
		Y: float64(endTileY),
	}, s.Manager.Conn)

	if err != nil {
		return nil, fmt.Errorf("failed to move ship: %w", err)
	}

	s.route.Active = true

	return resp, nil
}

func (s *Ship) Repair() (inbound.RepairShipResp, error) {
	resp, err := inbound.RepairShip(inbound.RepairShipReq{
		ShipID: s.RawData.ID,
	})
	if err != nil {
		return inbound.RepairShipResp{}, fmt.Errorf("failed to repair ship: %w", err)
	}
	return resp, nil
}
