package game

import (
	"fmt"
	"log/slog"

	"github.com/jhuggett/frontend/utils/callback"
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/outbound"
	"github.com/jhuggett/sea/start"
	"github.com/jhuggett/sea/timeline"
)

type OnTimeChangedCallback func(outbound.TimeChangedReq) error
type OnShipChangedCallback func(outbound.ShipChangedReq) error
type OnShipMovedCallback func(outbound.ShipMovedReq) error
type OnShipInventoryChangedCallback func(outbound.ShipInventoryChangedReq) error
type Manager struct {
	PlayerShip *Ship
	WorldMap   *WorldMap

	Snapshot *game_context.Snapshot
	Conn     *start.Connection

	tileSize int

	OnTimeChangedCallback          callback.CallbackRegistry[OnTimeChangedCallback]
	OnShipChangedCallback          callback.CallbackRegistry[OnShipChangedCallback]
	OnShipMovedCallback            callback.CallbackRegistry[OnShipMovedCallback]
	OnShipInventoryChangedCallback callback.CallbackRegistry[OnShipInventoryChangedCallback]

	LastTimeChangedReq outbound.TimeChangedReq
}

func NewManager(snapshot *game_context.Snapshot) *Manager {
	return &Manager{
		Snapshot:                       snapshot,
		OnTimeChangedCallback:          callback.CallbackRegistry[OnTimeChangedCallback]{},
		OnShipChangedCallback:          callback.CallbackRegistry[OnShipChangedCallback]{},
		OnShipMovedCallback:            callback.CallbackRegistry[OnShipMovedCallback]{},
		OnShipInventoryChangedCallback: callback.CallbackRegistry[OnShipInventoryChangedCallback]{},
	}
}

func (m *Manager) Start() error {

	conn := &start.Connection{
		Receiver: &outbound.Receiver{
			OnCrewInformation: func(cir outbound.CrewInformationReq) (outbound.CrewInformationResp, error) {
				slog.Info("CrewInformation called", "req", cir)
				return outbound.CrewInformationResp{}, nil
			},
			OnShipMoved: func(smr outbound.ShipMovedReq) (outbound.ShipMovedResp, error) {
				slog.Info("ShipMoved called", "req", smr)

				m.PlayerShip.RawData.X = float64(smr.Location.X)
				m.PlayerShip.RawData.Y = float64(smr.Location.Y)

				m.OnShipMovedCallback.InvokeEndToStart(func(oshmc OnShipMovedCallback) error {
					return oshmc(smr)
				})

				return outbound.ShipMovedResp{}, nil
			},
			OnShipDocked: func(sdr outbound.ShipDockedReq) (outbound.ShipDockedResp, error) {
				slog.Info("ShipDocked called", "req", sdr)
				return outbound.ShipDockedResp{}, nil
			},
			OnTimeChanged: func(tcr outbound.TimeChangedReq) (outbound.TimeChangedResp, error) {
				slog.Info("TimeChanged called", "req", tcr)

				m.LastTimeChangedReq = tcr

				m.OnTimeChangedCallback.InvokeEndToStart(func(otcc OnTimeChangedCallback) error {
					return otcc(tcr)
				})

				return outbound.TimeChangedResp{}, nil
			},
			OnShipInventoryChanged: func(sicr outbound.ShipInventoryChangedReq) (outbound.ShipInventoryChangedResp, error) {
				slog.Info("ShipInventoryChanged called", "req", sicr)

				m.OnShipInventoryChangedCallback.InvokeEndToStart(func(osicc OnShipInventoryChangedCallback) error {
					return osicc(sicr)
				})

				return outbound.ShipInventoryChangedResp{}, nil
			},
			OnShipChanged: func(scr outbound.ShipChangedReq) (outbound.ShipChangedResp, error) {
				slog.Info("ShipChanged called", "req", scr)

				// add a panel for this

				m.OnShipChangedCallback.InvokeEndToStart(func(otcc OnShipChangedCallback) error {
					return otcc(scr)
				})

				return outbound.ShipChangedResp{}, nil
			},
		},
		GameCtx: &game_context.GameContext{},
	}

	if m.Snapshot == nil {
		resp, err := inbound.Register(inbound.RegisterReq{})
		if err != nil {
			return err
		}

		m.Snapshot = &resp.GameCtx

		slog.Info("Register response", "resp", resp)
	}

	conn.GameCtx = game_context.New(*m.Snapshot)

	loginResp, err := inbound.Login(inbound.LoginReq{
		Snapshot: *m.Snapshot,
	}, conn)
	if err != nil {
		return err
	}

	m.PlayerShip = &Ship{
		Manager: m,
		RawData: loginResp.Ship,
	}

	slog.Info("Login response", "loginResp", loginResp)

	conn.GameCtx.Timeline = timeline.New()

	m.Conn = conn

	start.Game(conn)

	m.WorldMap = NewWorldMap(m)

	err = m.WorldMap.Load()
	if err != nil {
		return fmt.Errorf("failed to load world map: %w", err)
	}

	return nil
}

func (m *Manager) TileSize() int {
	return m.tileSize
}

func (m *Manager) SetTileSize(tileSize int) {
	m.tileSize = tileSize
}

func (s *Manager) PlotRoute(x, y int) (*inbound.PlotRouteResp, error) {
	resp, err := inbound.PlotRoute(inbound.PlotRouteReq{
		Coordinate: inbound.Coordinate{
			X: float64(x),
			Y: float64(y),
		},
	}, s.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed to plot route: %w", err)
	}

	return &resp, nil
}

func (s *Manager) ControlTime(req inbound.ControlTimeReq) (inbound.ControlTimeResp, error) {
	return inbound.ControlTime(s.Conn, req)
}

/*

Ship ---
shipImg := ebiten.NewImage(
		int(page.TileSize),
		int(page.TileSize),
	)
	vector.DrawFilledRect(
		shipImg,
		0,
		0,
		float32(page.TileSize),
		float32(page.TileSize),
		color.RGBA{
			B: 255,
			G: 255,
			A: 255,
		},
		false,
	)

	page.Ship.Image = shipImg
	---

	WorldMap ---
	gameMap, err := inbound.GetWorldMap(inbound.GetWorldMapReq{}, snapshot.GameMapID)
	if err != nil {
		return nil, err
	}

	slog.Info("GetWorldMap response", "gameMap", gameMap)

	a := uint8(50)
	b := uint8(100)
	lowestColor := color.RGBA{
		R: 23 + a,
		G: 18 + a,
		B: 13 + a,
		A: 255,
	}
	highestColor := color.RGBA{
		R: 58 + b,
		G: 45 + b,
		B: 25 + b,
		A: 255,
	}

	page.Continents = []Continent{}
	for _, continent := range gameMap.Continents {
		smallestX := 0
		smallestY := 0
		largestX := 0
		largestY := 0

		for _, point := range continent.Points {
			if point.X < smallestX {
				smallestX = point.X
			}
			if point.Y < smallestY {
				smallestY = point.Y
			}
			if point.X > largestX {
				largestX = point.X
			}
			if point.Y > largestY {
				largestY = point.Y
			}

			if page.SmallestPointX > point.X {
				page.SmallestPointX = point.X
			}
			if page.SmallestPointY > point.Y {
				page.SmallestPointY = point.Y
			}
			if page.LargestPointX < point.X {
				page.LargestPointX = point.X
			}
			if page.LargestPointY < point.Y {
				page.LargestPointY = point.Y
			}
		}

		img := ebiten.NewImage(
			(largestX-smallestX+1)*int(page.TileSize),
			(largestY-smallestY+1)*int(page.TileSize),
		)
		for _, point := range continent.Points {
			e := (point.Elevation - .5) * 2
			vector.DrawFilledRect(
				img,
				float32(float64(point.X-smallestX)*page.TileSize),
				float32(float64(point.Y-smallestY)*page.TileSize),
				float32(page.TileSize),
				float32(page.TileSize),
				color.RGBA{
					R: uint8(float64(highestColor.R-lowestColor.R)*e) + lowestColor.R,
					G: uint8(float64(highestColor.G-lowestColor.G)*e) + lowestColor.G,
					B: uint8(float64(highestColor.B-lowestColor.B)*e) + lowestColor.B,
					A: 255,
				},
				false,
			)
		}
		page.Continents = append(page.Continents, Continent{
			OriginX: smallestX,
			OriginY: smallestY,
			Image:   img,
		})
	}

	portsData, err := inbound.GetPorts(inbound.GetPortsReq{}, snapshot.GameMapID)
	if err != nil {
		return nil, err
	}

	slog.Info("GetPorts response", "portsData", portsData)

	page.Ports = []Port{}

	for _, port := range portsData.Ports {
		img := ebiten.NewImage(int(page.TileSize), int(page.TileSize))
		vector.DrawFilledRect(
			img,
			float32(page.TileSize/4),
			float32(page.TileSize/4),
			float32(page.TileSize/2),
			float32(page.TileSize/2),
			color.RGBA{
				R: 255,
				A: 255,
			},
			false,
		)
		page.Ports = append(page.Ports, Port{
			X:     port.Point.X,
			Y:     port.Point.Y,
			Image: img,
		})
	}

	// should add border
	mapImg := ebiten.NewImage(
		int(float64(page.LargestPointX-page.SmallestPointX+1)*page.TileSize),
		int(float64(page.LargestPointY-page.SmallestPointY+1)*page.TileSize),
	)

	mapImg.Fill(color.RGBA{
		R: 220,
		G: 202,
		B: 127,
		A: 255,
	})

	page.MapImage = mapImg
	---



*/
