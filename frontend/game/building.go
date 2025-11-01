package game

import "github.com/jhuggett/sea/inbound"

type Building struct {
	Port *Port

	Name string
	Type string

	X int
	Y int
}

type HireablePerson struct {
	FirstName string
	LastName  string

	NickName string

	Age int

	PlaceOfResidence string
}

func (b *Building) GetHireablePeople() ([]*HireablePerson, error) {
	resp, err := inbound.GetHirablePeopleAtPort(inbound.GetHirablePeopleAtPortReq{
		PortID: b.Port.RawData.ID,
	})
	if err != nil {
		return nil, err
	}

	var result []*HireablePerson
	for _, p := range resp.People {
		result = append(result, &HireablePerson{
			FirstName:        p.FirstName,
			LastName:         p.LastName,
			NickName:         p.NickName,
			Age:              p.Age,
			PlaceOfResidence: p.PlaceOfResidence,
		})
	}

	return result, nil
}

func (b *Building) HireCrewMember(person *HireablePerson) error {
	_, err := inbound.HireCrew(inbound.HireCrewReq{
		People: []inbound.Person{
			{
				FirstName:        person.FirstName,
				LastName:         person.LastName,
				NickName:         person.NickName,
				Age:              person.Age,
				PlaceOfResidence: person.PlaceOfResidence,
			},
		},
	}, b.Port.Manager.Conn)

	return err

}
