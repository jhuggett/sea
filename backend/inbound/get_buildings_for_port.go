package inbound

import (
	"github.com/jhuggett/sea/data/building"
)

type GetBuildingsForPortReq struct {
	PortID uint
}

type Building struct {
	PortID uint

	Name string
	Type string

	X int
	Y int
}

type GetBuildingsForPortResp struct {
	Buildings []Building
}

func GetBuildingsForPort(conn Connection, req GetBuildingsForPortReq) (GetBuildingsForPortResp, error) {
	buildings, err := building.Where(req.PortID)
	if err != nil {
		return GetBuildingsForPortResp{}, err
	}

	buildingsResp := []Building{}
	for _, building := range buildings {
		buildingsResp = append(buildingsResp, Building{
			PortID: building.Persistent.PortID,
			Name:   building.Persistent.Name,
			Type:   building.Persistent.Type,
			X:      building.Persistent.X,
			Y:      building.Persistent.Y,
		})
	}

	return GetBuildingsForPortResp{
		Buildings: buildingsResp,
	}, nil
}
