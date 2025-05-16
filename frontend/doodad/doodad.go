package doodad

import "github.com/hajimehoshi/ebiten/v2"

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
