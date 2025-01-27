package inbound

import (
	"encoding/json"

	"github.com/jhuggett/sea/data/person"
	"github.com/jhuggett/sea/data/port"
)

type GetHirablePeopleAtPortReq struct {
	PortID uint `json:"port_id"`
}

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	NickName string `json:"nick_name"`

	Age int `json:"age"`

	// Not sure what this should be called, basically where they are from
	PlaceOfResidence string `json:"place_of_residence"`
}

type GetHirablePeopleAtPortResp struct {
	People []Person `json:"people"`
}

func GetHirablePeopleAtPort(conn Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r GetHirablePeopleAtPortReq
		if err := json.Unmarshal(req, &r); err != nil {
			return nil, err
		}

		port, err := port.Get(r.PortID)
		if err != nil {
			return nil, err
		}

		people := person.GeneratePeople(10)

		resp := GetHirablePeopleAtPortResp{
			People: []Person{},
		}

		for _, p := range people {
			resp.People = append(resp.People, Person{
				FirstName:        p.FirstName,
				LastName:         p.LastName,
				NickName:         p.NickName,
				Age:              int(p.Age),
				PlaceOfResidence: port.Persistent.Name,
			})
		}
		return resp, nil
	}
}
