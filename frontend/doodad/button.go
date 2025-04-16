package doodad

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
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

	position   Position
	dimensions Rectangle

	hovered bool
}

func (w *Button) Update() error {

	// // Get mouse position
	// mouseX, mouseY := ebiten.CursorPosition()
	// if mouseX >= w.position.X && mouseX <= w.position.X+w.dimensions.Width &&
	// 	mouseY >= w.position.Y && mouseY <= w.position.Y+w.dimensions.Height {

	// 	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	// 		if w.OnClick != nil {
	// 			w.OnClick()
	// 		}
	// 	} else {
	// 		if !w.hovered {
	// 			w.hovered = true
	// 			w.hovering()
	// 		}
	// 	}
	// } else {
	// 	if w.hovered {
	// 		w.hovered = false
	// 		w.stoppedHovering()
	// 	}
	// }

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
	w.nonHoveredLabel.SetMessage(w.message)

	w.dimensions = w.nonHoveredLabel.Dimensions()

	w.hoveredLabel = NewLabel()
	w.hoveredLabel.SetMessage("Click me!")

	withinBounds := func(x, y int) bool {
		return x >= w.position.X && x <= w.position.X+w.dimensions.Width &&
			y >= w.position.Y && y <= w.position.Y+w.dimensions.Height
	}

	w.Gesturer.OnMouseMove(func(x, y int) bool {
		if withinBounds(x, y) {
			w.hovering()
			// return true
		} else {
			w.stoppedHovering()
		}
		return false
	})

	w.Gesturer.OnMouseUp(func(x, y int) bool {
		fmt.Println("Button.OnClick", x, y)
		if withinBounds(x, y) {
			if w.OnClick != nil {
				w.OnClick()
			}
			return true
		}
		return false
	})

	return nil
}

func (w *Button) SetMessage(message string) {
	w.message = message
	w.nonHoveredLabel.SetMessage(message)
	w.dimensions = w.nonHoveredLabel.Dimensions()
}

func (w *Button) SetPosition(position Position) {
	w.position = position
	w.nonHoveredLabel.SetPosition(position)
	w.hoveredLabel.SetPosition(position)
}

func (w *Button) Position() Position {
	return w.position
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
