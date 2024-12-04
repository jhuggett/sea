package coordination

import "fmt"

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p Point) North() Point {
	return Point{
		X: p.X,
		Y: p.Y + 1,
	}
}

func (p Point) South() Point {
	return Point{
		X: p.X,
		Y: p.Y - 1,
	}
}

func (p Point) East() Point {
	return Point{
		X: p.X + 1,
		Y: p.Y,
	}
}

func (p Point) West() Point {
	return Point{
		X: p.X - 1,
		Y: p.Y,
	}
}

func (p Point) SameAs(to Point) bool {
	return p.X == to.X && p.Y == to.Y
}

type Pointable interface {
	Point() Point
}

func (p Point) Subtract(point Point) Point {
	p.X -= point.X
	p.Y -= point.Y
	return p
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func (p Point) Adjacent() []Point {
	return []Point{
		p.North(),
		p.East(),
		p.South(),
		p.West(),
	}
}

func (p Point) AdjacentDiagonal() []Point {
	return []Point{
		p.North().East(),
		p.South().East(),
		p.North().West(),
		p.South().West(),
	}
}

func (p Point) AdjacentDiagonalAndCardinal() []Point {
	return append(p.Adjacent(), p.AdjacentDiagonal()...)
}
