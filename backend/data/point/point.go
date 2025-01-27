package point

import (
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/utils/coordination"
)

func Find(location coordination.Point, worldMapId uint) (*data.Point, error) {

	slog.Info("Find", "location", location)

	var point data.Point
	err := db.Conn().Where("x = ?", location.X).Where("y = ?", location.Y).Where("world_map_id = ?", worldMapId).First(&point).Error
	if err != nil {
		slog.Error("Failed to find point", "location", location, "error", err)
		return nil, fmt.Errorf("failed to find point: %w", err)
	}

	slog.Info("Found point", "point", point)

	return &point, nil
}
