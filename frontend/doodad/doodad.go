package doodad

import "github.com/hajimehoshi/ebiten/v2"

type Position struct {
	X int
	Y int
}

type Rectangle struct {
	Width  int
	Height int
}

type Doodad interface {
	Update() error
	Draw(screen *ebiten.Image)
	Setup() error
}

type Rectangular interface {
	Dimensions() Rectangle
}

type Positioned interface {
	Position() Position
}
