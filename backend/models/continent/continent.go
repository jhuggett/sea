package continent

import (
	"log/slog"

	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/coastal_point"
	"github.com/jhuggett/sea/utils/coordination"
	"github.com/soniakeys/raycast"
)

type Continent struct {
	Persistent models.Continent
}

// Continents should a random resource multiplier for raw resources that affect how much of a resource is produced.
// Maybe: Then ports can have a random resource multiplier between 0 and the continent's resource multiplier for further variance.

// A continents population acts in essence as a business. It has demands for it's people that it buys out of its wealth from ports.

type PointInformation struct {
	CoastalPoint *coastal_point.CoastalPoint
}

func (c *Continent) Contains(point coordination.Point) (isWithin bool, information PointInformation, err error) {

	pointXY := raycast.XY{
		X: float64(point.X),
		Y: float64(point.Y),
	}

	poly := raycast.Poly{}

	coordination.Sort(c.Persistent.CoastalPoints)

	for _, coastalPoint := range c.Persistent.CoastalPoints {
		if point.SameAs(coastalPoint.Point()) {
			return true, PointInformation{
				CoastalPoint: &coastal_point.CoastalPoint{
					Persistent: *coastalPoint,
				},
			}, nil
		}
		poly = append(poly, raycast.XY{
			X: float64(coastalPoint.X),
			Y: float64(coastalPoint.Y),
		})
	}
	if pointXY.In(poly) {
		slog.Info("point point is in a continent")
		return true, PointInformation{}, nil
	}

	return false, PointInformation{}, nil
}