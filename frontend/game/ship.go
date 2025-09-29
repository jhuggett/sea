package game

import (
	"fmt"

	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/outbound"
)

type Ship struct {
	Manager *Manager
	RawData inbound.ShipInfo

	Route *ShipRoute
}

func (s *Ship) Location() inbound.Point {
	return inbound.Point{
		X: int(s.RawData.X),
		Y: int(s.RawData.Y),
	}
}

type ShipRoute struct {
	Points       []inbound.Coordinate
	Active       bool
	ShipMovedReq *outbound.ShipMovedReq
}

func (s *Ship) HasRoute() bool {
	return s.Route != nil && len(s.Route.Points) > 0
}

func (s *Ship) IsRouteActive() bool {
	return s.Route != nil && s.Route.Active
}

func (s *Ship) PlotRoute(x, y int) (*ShipRoute, error) {
	resp, err := s.Manager.PlotRoute(x, y)
	if err != nil {
		return nil, fmt.Errorf("failed to plot route: %w", err)
	}

	route := &ShipRoute{
		Points: resp.Coordinates,
	}

	s.Route = route

	return route, nil
}

func (s *Ship) SetSail() (*inbound.MoveShipResp, error) {
	if s.Route == nil {
		return nil, fmt.Errorf("no route set")
	}

	s.Route.Active = true

	endTileX := int(s.Route.Points[len(s.Route.Points)-1].X)
	endTileY := int(s.Route.Points[len(s.Route.Points)-1].Y)

	resp, err := inbound.MoveShip(inbound.MoveShipReq{
		X: float64(endTileX),
		Y: float64(endTileY),
	}, s.Manager.Conn)

	if err != nil {
		s.Route.Active = false
		return nil, fmt.Errorf("failed to move ship: %w", err)
	}

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

func (s *Ship) TriggerShipInfoRequest() error {
	err := s.Manager.RequestShipInfo(int(s.RawData.ID))
	return err
}

func (s *Ship) TriggerShipInventoryRequest() error {
	err := s.Manager.RequestShipInventoryInfo()
	return err
}

func (s *Ship) TriggerCrewInfoRequest() error {
	err := s.Manager.RequestCrewInfo()
	return err
}
