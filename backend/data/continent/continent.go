package continent

import (
	"fmt"

	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/utils/coordination"
)

type Continent struct {
	Persistent data.Continent
}

// Continents should a random resource multiplier for raw resources that affect how much of a resource is produced.
// Maybe: Then ports can have a random resource multiplier between 0 and the continent's resource multiplier for further variance.

// A continents population acts in essence as a business. It has demands for it's people that it buys out of its wealth from ports.

var ErrNotInContinent = fmt.Errorf("point is not in continent")

func (c *Continent) Contains(point coordination.Point) (*data.Point, error) {
	for _, landPoint := range c.Persistent.Points {
		if landPoint.Point().SameAs(point) {
			return landPoint, nil
		}
	}

	return nil, ErrNotInContinent
}

// Get rid of this
func (c *Continent) GetCoastalPoints() []*data.Point {
	coastalPoints := []*data.Point{}

	for _, point := range c.Persistent.Points {
		if point.Coastal {
			coastalPoints = append(coastalPoints, point)
		}
	}

	return coastalPoints
}

// Use this in future
func (c *Continent) CoastalPoints() ([]*data.Point, error) {
	var points []*data.Point
	err := db.Conn().Where("continent_id = ?", c.Persistent.ID).Where("coastal = ?", true).Find(&points).Error

	if err != nil {
		return nil, err
	}

	return points, nil
}

func (c *Continent) LoadPoints() ([]*data.Point, error) {
	points := []*data.Point{}
	err := db.Conn().Where("continent_id = ?", c.Persistent.ID).Find(&points).Error
	if err != nil {
		return nil, err
	}
	c.Persistent.Points = points
	return points, nil
}

func (c *Continent) Population() (data.Population, error) {
	var p data.Population
	err := db.Conn().Where("continent_id = ?", c.Persistent.ID).First(&p).Error
	return p, err
}

func Using(continent data.Continent) *Continent {
	return &Continent{
		Persistent: continent,
	}
}
