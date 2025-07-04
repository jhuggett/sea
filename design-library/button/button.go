package button

import (
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	OnClick  func()
	Gesturer doodad.Gesturer
	label.Config
}

func New(config Config) *Button {
	button := &Button{
		Gesturer: config.Gesturer,
	}
	button.message = config.Message
	button.OnClick = config.OnClick

	button.labelConfig = config.Config

	if config.Layout != nil {
		button.Box = config.Layout
	} else {
		button.Box = box.New(box.Config{})
	}

	return button
}

type Button struct {
	message string

	labelConfig label.Config

	onSetMessage func(message string)

	nonHoveredChildren doodad.Children
	hoveredChildren    doodad.Children

	OnClick func()

	Gesturer doodad.Gesturer

	hovered bool

	doodad.Default
}

func (w *Button) Update() error {
	return nil
}

func (w *Button) Draw(screen *ebiten.Image) {
	if w.hovered {
		w.hoveredChildren.Draw(screen)
	} else {
		w.nonHoveredChildren.Draw(screen)
	}
}

func (w *Button) Setup() {
	nonHoveredLabel := label.New(w.labelConfig)
	nonHoveredLabel.Setup()
	w.nonHoveredChildren.Add(nonHoveredLabel)

	hoveredLabel := label.New(w.labelConfig)
	hoveredLabel.Setup()
	w.hoveredChildren.Add(hoveredLabel)

	withinBounds := func(x, y int) bool {
		return x >= w.Box.X() && x <= w.Box.X()+w.Box.Width() &&
			y >= w.Box.Y() && y <= w.Box.Y()+w.Box.Height()
	}

	w.Gesturer.OnMouseMove(func(x, y int) error {
		if withinBounds(x, y) {
			w.hovering()
			// return callback.ErrStopPropagation
		} else {
			w.stoppedHovering()
		}
		return nil
	})

	w.Gesturer.OnMouseUp(func(event doodad.MouseUpEvent) error {
		if event.Button != ebiten.MouseButtonLeft {
			return nil
		}
		x, y := event.X, event.Y
		if withinBounds(x, y) {
			if w.OnClick != nil {
				w.OnClick()
			}
			return doodad.ErrStopPropagation
		}
		return nil
	})

	w.onSetMessage = func(message string) {
		nonHoveredLabel.SetMessage(message)
		hoveredLabel.SetMessage(message)
	}
}

func (w *Button) SetMessage(message string) {
	w.message = message
	w.onSetMessage(message)

}

// func (w *Button) SetPosition(position func() doodad.Position) {
// 	// w.position = position
// 	// w.nonHoveredLabel.SetPosition(position)
// 	// w.hoveredLabel.SetPosition(position)

// 	// Instead of setting these positions, they should be already dependent on the button's position, we shouldn't have to be able to access them

// }

func (w *Button) Teardown() error {

	// TODO: Implement button teardown logic

	err := w.hoveredChildren.Teardown()
	if err != nil {
		return fmt.Errorf("failed to teardown hovered children: %w", err)
	}
	err = w.nonHoveredChildren.Teardown()
	if err != nil {
		return fmt.Errorf("failed to teardown non-hovered children: %w", err)
	}

	return nil
}

// func (w *Button) Position() doodad.Position {
// 	return w.position()
// }

// func (w *Button) Dimensions() doodad.Rectangle {
// 	return w.dimensions
// }

func (w *Button) hovering() {
	w.hovered = true
}

func (w *Button) stoppedHovering() {
	w.hovered = false
}
