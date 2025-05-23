package doodad

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewPanel(gesturer Gesturer) *Panel {
	panel := &Panel{
		Gesturer: gesturer,
	}
	panel.Setup()
	return panel
}

type Panel struct {
	position   func() Position
	dimensions Rectangle

	Gesturer Gesturer

	bg *ebiten.Image
}

func (w *Panel) Update() error {
	return nil
}

func (w *Panel) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w.position().X), float64(w.position().Y))
	screen.DrawImage(w.bg, op)
}

func (w *Panel) Setup() error {

	w.bg = ebiten.NewImage(200, 200)
	w.bg.Fill(color.RGBA{0, 0, 0, 50})
	w.dimensions = Rectangle{
		Width:  200,
		Height: 200,
	}

	w.Gesturer.OnMouseMove(func(x, y int) error {
		if x >= w.position().X && x <= w.position().X+w.dimensions.Width &&
			y >= w.position().Y && y <= w.position().Y+w.dimensions.Height {
			return ErrStopPropagation
		}
		return nil
	})

	return nil
}

func (w *Panel) SetPosition(position func() Position) {
	w.position = position
}

func (w *Panel) Position() Position {
	return w.position()
}

func (w *Panel) Dimensions() Rectangle {
	return w.dimensions
}
