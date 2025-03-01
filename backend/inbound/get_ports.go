package inbound

import (
	"encoding/json"

	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/data/port"
)

type GetPortsReq struct {
}

type Port struct {
	ID uint `json:"id"`

	Point data.Point `json:"point"`

	Name string `json:"name"`
}

type GetPortsResp struct {
	Ports []Port `json:"ports"`
}

func GetPorts(r GetPortsReq, gameMapID uint) (GetPortsResp, error) {
	ports, err := port.All(gameMapID)
	if err != nil {
		return GetPortsResp{}, err
	}

	resp := GetPortsResp{
		Ports: []Port{},
	}

	for _, p := range ports {
		port := Port{}

		port.ID = p.Persistent.ID
		port.Point = *p.Persistent.Point
		port.Name = p.Persistent.Name

		resp.Ports = append(resp.Ports, port)
	}

	return resp, nil
}

func WSGetPorts(conn Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r GetPortsReq
		if err := json.Unmarshal(req, &r); err != nil {
			return nil, err
		}
		return GetPorts(r, conn.Context().GameMapID())
	}
}
