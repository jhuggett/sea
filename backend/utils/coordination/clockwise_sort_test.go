package coordination

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClockwiseSort(t *testing.T) {
	points := []Point{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
		{X: -1, Y: 0},
		{X: 0, Y: -1},
	}

	center := Point{X: 0, Y: 0}

	sort.SliceStable(points, func(i, j int) bool {
		return comparePoints(center, points[i], points[j])
	})

	assert.Equal(t, []Point{
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 0, Y: 1},
		{X: -1, Y: 1},
		{X: -1, Y: 0},
		{X: -1, Y: -1},
		{X: 0, Y: -1},
		{X: 1, Y: -1},
	}, points)
}

func TestClockwiseSort2(t *testing.T) {
	points := []Point{
		{X: 0, Y: -1},
		{X: -1, Y: -1},
		{X: 1, Y: 1},
	}

	center := Point{X: 0, Y: 0}

	sort.SliceStable(points, func(i, j int) bool {
		return comparePoints(center, points[i], points[j])
	})

	assert.Equal(t, []Point{
		{X: 1, Y: 1},
		{X: -1, Y: -1},
		{X: 0, Y: -1},
	}, points)
}

func TestClockwiseSort3(t *testing.T) {
	points := []Point{
		{X: -1, Y: -1},
		{X: 1, Y: 1},
	}

	center := Point{X: 0, Y: 0}

	sort.SliceStable(points, func(i, j int) bool {
		return comparePoints(center, points[i], points[j])
	})

	assert.Equal(t, []Point{
		{X: 1, Y: 1},
		{X: -1, Y: -1},
	}, points)
}
