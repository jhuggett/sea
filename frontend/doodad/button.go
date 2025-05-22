package doodad

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/utils/callback"
)

func NewButton(
	message string,
	onClick func(),
	gesturer Gesturer,
) *Button {
	button := &Button{
		Gesturer: gesturer,
	}
	button.message = message
	button.OnClick = onClick
	button.Setup()
	return button
}

type Button struct {
	message string

	nonHoveredLabel *Label
	hoveredLabel    *Label

	OnClick func()

	Gesturer Gesturer

	position   func() Position
	dimensions Rectangle

	hovered bool
}

func (w *Button) Update() error {
	return nil
}

func (w *Button) Draw(screen *ebiten.Image) {
	if w.hovered {
		w.hoveredLabel.Draw(screen)
	} else {
		w.nonHoveredLabel.Draw(screen)
	}
}

func (w *Button) Setup() error {
	w.nonHoveredLabel = NewLabel()
	w.nonHoveredLabel.Setup()
	w.nonHoveredLabel.SetMessage(w.message)

	w.dimensions = w.nonHoveredLabel.Dimensions()

	w.hoveredLabel = NewLabel()
	w.hoveredLabel.FontSize = w.nonHoveredLabel.FontSize + 1
	w.hoveredLabel.BackgroundColor = color.RGBA{20, 10, 2, 50}
	w.hoveredLabel.Setup()
	w.hoveredLabel.SetMessage(w.message)

	withinBounds := func(x, y int) bool {
		return x >= w.position().X && x <= w.position().X+w.dimensions.Width &&
			y >= w.position().Y && y <= w.position().Y+w.dimensions.Height
	}

	w.Gesturer.OnMouseMove(func(x, y int) error {
		if withinBounds(x, y) {
			w.hovering()
			return callback.ErrStopPropagation
		} else {
			w.stoppedHovering()
		}
		return nil
	})

	w.Gesturer.OnMouseUp(func(event MouseUpEvent) error {
		if event.Button != ebiten.MouseButtonLeft {
			return nil
		}
		x, y := event.X, event.Y
		fmt.Println("Button.OnClick", x, y)
		if withinBounds(x, y) {
			if w.OnClick != nil {
				w.OnClick()
			}
			return ErrStopPropagation
		}
		return nil
	})

	return nil
}

func (w *Button) SetMessage(message string) {
	w.message = message
	w.nonHoveredLabel.SetMessage(message)
	w.dimensions = w.nonHoveredLabel.Dimensions()
}

func (w *Button) SetPosition(position func() Position) {
	w.position = position
	w.nonHoveredLabel.SetPosition(position)
	w.hoveredLabel.SetPosition(position)
}

func (w *Button) Position() Position {
	return w.position()
}

func (w *Button) Dimensions() Rectangle {
	return w.dimensions
}

func (w *Button) hovering() {
	w.hovered = true
}

func (w *Button) stoppedHovering() {
	w.hovered = false
}
