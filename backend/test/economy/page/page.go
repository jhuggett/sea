package page

import "github.com/hajimehoshi/ebiten/v2"

type Page interface {
	Update() error
	Draw(screen *ebiten.Image)

	SetWidthAndHeight(width, height int)
}

type PageControls interface {
	Push(Page)
	Pop()
}
