package person

import (
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/name"
)

type Person struct {
	Persistent data.Person
}

func GeneratePeople(count int) []data.Person {
	people := []data.Person{}

	for i := 0; i < count; i++ {
		people = append(people, GeneratePerson())
	}

	return people
}

func GeneratePerson() data.Person {
	return data.Person{
		FirstName: name.Generate(2),
		LastName:  name.Generate(3),
		NickName:  name.GenerateNickName(),

		Age: uint(20),
	}
}

func Get(id uint) (*Person, error) {
	var s data.Person
	err := db.Conn().First(&s, id).Error
	if err != nil {
		return nil, err
	}

	return &Person{
		Persistent: s,
	}, nil
}
