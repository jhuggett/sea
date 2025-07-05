package button

import (
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"fmt"
	"image/color"

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

	OnClick func()

	Gesturer doodad.Gesturer

	hovered bool

	doodad.Default

	showHoveringLabel   func()
	showNonHoveredLabel func()
}

func (w *Button) Update() error {
	return nil
}

func (w *Button) Setup() {
	nonHoveredLabel := label.New(label.Config{
		BackgroundColor: w.labelConfig.BackgroundColor,
		ForegroundColor: w.labelConfig.ForegroundColor,
		Message:         w.message,
		FontSize:        w.labelConfig.FontSize,
		Padding:         w.labelConfig.Padding,
		Layout:          w.Box,
	})
	nonHoveredLabel.Setup()
	w.AddChild(nonHoveredLabel)

	hoveredLabel := label.New(label.Config{
		BackgroundColor: w.labelConfig.BackgroundColor,
		ForegroundColor: color.RGBA{
			R: 255,
			G: 0,
			B: 100,
			A: 255,
		},
		Message:  w.message,
		FontSize: w.labelConfig.FontSize,
		Padding:  w.labelConfig.Padding,
		Layout:   w.Box,
	})
	hoveredLabel.Setup()
	w.AddChild(hoveredLabel)

	w.showHoveringLabel = func() {
		hoveredLabel.Show()
		nonHoveredLabel.Hide()
	}
	w.showNonHoveredLabel = func() {
		hoveredLabel.Hide()
		nonHoveredLabel.Show()
	}

	w.hovered = true // so that it'll trigger the non-hovered event first

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

	err := w.Teardown()
	if err != nil {
		return fmt.Errorf("failed to teardown button: %w", err)
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
	if w.hovered {
		return
	}

	w.hovered = true
	w.showHoveringLabel()

	ebiten.SetCursorShape(ebiten.CursorShapePointer) // Change cursor to text mode when hovering over the button
}

func (w *Button) stoppedHovering() {
	if !w.hovered {
		return
	}

	w.hovered = false
	w.showNonHoveredLabel()
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
}
