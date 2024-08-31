package world_map

import (
	"errors"
	"log/slog"

	"github.com/soniakeys/raycast"
	"gorm.io/gorm"
)

type Continent struct {
	gorm.Model

	WorldMapID uint

	CoastalPoints []*CoastalPoint `gorm:"foreignKey:ContinentID"`
}

// Continents should a random resource multiplier for raw resources that affect how much of a resource is produced.
// Maybe: Then ports can have a random resource multiplier between 0 and the continent's resource multiplier for further variance.

// A continents population acts in essence as a business. It has demands for it's people that it buys out of its wealth from ports.

type PointInformation struct {
	CoastalPoint *CoastalPoint
}

func (c *Continent) ContainsPoint(point Point) (isWithin bool, information PointInformation, err error) {

	pointXY := raycast.XY{
		X: float64(point.X),
		Y: float64(point.Y),
	}

	poly := raycast.Poly{}

	Sort(c.CoastalPoints)

	for _, coastalPoint := range c.CoastalPoints {

		slog.Info("Coastal Point", "x", coastalPoint.X, "y", coastalPoint.Y)
		if point.SameAs(coastalPoint.Point()) {
			slog.Info("point point is a coastal point")
			return nil, errors.New("point point is a coastal point")
		}
		poly = append(poly, raycast.XY{
			X: float64(coastalPoint.X),
			Y: float64(coastalPoint.Y),
		})
	}
	if pointXY.In(poly) {
		slog.Info("point point is in a continent")
		return nil, errors.New("point point is in a continent")
	}

	return false, PointInformation{}, nil
}
