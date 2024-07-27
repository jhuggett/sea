package world_map

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p *Point) North() *Point {
	return &Point{
		X: p.X,
		Y: p.Y + 1,
	}
}

func (p *Point) South() *Point {
	return &Point{
		X: p.X,
		Y: p.Y - 1,
	}
}

func (p *Point) East() *Point {
	return &Point{
		X: p.X + 1,
		Y: p.Y,
	}
}

func (p *Point) West() *Point {
	return &Point{
		X: p.X - 1,
		Y: p.Y,
	}
}

func (p *Point) SameAs(to *Point) bool {
	return p.X == to.X && p.Y == to.Y
}

type Pointable interface {
	Point() *Point
}

func (p *Point) Subtract(point Point) *Point {
	p.X -= point.X
	p.Y -= point.Y
	return p
}
