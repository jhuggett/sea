package camera

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
)

type Camera struct {
	ViewPort   f64.Vec2
	Position   f64.Vec2
	ZoomFactor float64
	Rotation   float64

	TileSize float64

	CameraIsMoving bool
	LastPosition   f64.Vec2
}

func (c *Camera) String() string {
	return fmt.Sprintf("Camera{ViewPort: %v, Position: %v, ZoomFactor: %v, Rotation: %v}", c.ViewPort, c.Position, c.ZoomFactor, c.Rotation)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] / 2,
		c.ViewPort[1] / 2,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])

	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])

	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	}
	screen.DrawImage(
		world,
		op,
	)
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) Reset() {
	c.Position[0] = 0
	c.Position[1] = 0
}
