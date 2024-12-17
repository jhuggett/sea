package continent

import (
	"fmt"

	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/utils/coordination"
)

type Continent struct {
	Persistent models.Continent
}

// Continents should a random resource multiplier for raw resources that affect how much of a resource is produced.
// Maybe: Then ports can have a random resource multiplier between 0 and the continent's resource multiplier for further variance.

// A continents population acts in essence as a business. It has demands for it's people that it buys out of its wealth from ports.

var ErrNotInContinent = fmt.Errorf("point is not in continent")

func (c *Continent) Contains(point coordination.Point) (*models.Point, error) {
	for _, landPoint := range c.Persistent.Points {
		if landPoint.Point().SameAs(point) {
			return landPoint, nil
		}
	}

	return nil, ErrNotInContinent
}

func (c *Continent) GetCoastalPoints() []*models.Point {
	coastalPoints := []*models.Point{}

	for _, point := range c.Persistent.Points {
		if point.Coastal {
			coastalPoints = append(coastalPoints, point)
		}
	}

	return coastalPoints
}
