package button

import (
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/reaction"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	OnClick func(*Button)
	label.Config
}

func New(config Config) *Button {
	button := &Button{}
	button.message = config.Message
	button.OnClick = config.OnClick

	button.labelConfig = config.Config

	if config.Layout != nil {
		button.Box = config.Layout
	} else {
		// button.Box = box.New(box.Config{})
	}

	button.labelConfig.BackgroundColor = color.RGBA{50, 50, 50, 100}
	button.labelConfig.Padding = label.Padding{Top: 5, Right: 10, Bottom: 5, Left: 10}

	return button
}

type Button struct {
	message string

	labelConfig label.Config

	onSetMessage func(message string)

	OnClick func(*Button)

	hovered bool

	doodad.Default

	showHoveringLabel   func()
	showNonHoveredLabel func()
}

func (w *Button) Setup() {
	nonHoveredLabel := label.New(label.Config{
		BackgroundColor: w.labelConfig.BackgroundColor,
		ForegroundColor: w.labelConfig.ForegroundColor,
		Message:         w.message,
		FontSize:        w.labelConfig.FontSize,
		Padding:         w.labelConfig.Padding,
	})
	w.AddChild(nonHoveredLabel)

	w.showHoveringLabel = func() {
		nonHoveredLabel.Config = label.Config{
			BackgroundColor: color.RGBA{50, 50, 50, 25},
			ForegroundColor: color.RGBA{
				R: 230,
				G: 255,
				B: 240,
				A: 255,
			},
			Message:  w.message,
			FontSize: w.labelConfig.FontSize,
			Padding:  w.labelConfig.Padding,
		}

		doodad.ReSetup(nonHoveredLabel)
	}
	w.showNonHoveredLabel = func() {
		nonHoveredLabel.Config = label.Config{
			BackgroundColor: w.labelConfig.BackgroundColor,
			ForegroundColor: w.labelConfig.ForegroundColor,
			Message:         w.message,
			FontSize:        w.labelConfig.FontSize,
			Padding:         w.labelConfig.Padding,
		}

		doodad.ReSetup(nonHoveredLabel)
	}

	w.onSetMessage = func(message string) {
		nonHoveredLabel.SetMessage(message)
	}

	w.showNonHoveredLabel() // Show the non-hovered label by default

	w.Children().Setup()

	w.Box.Computed(func(b *box.Box) {
		boundingBox := box.Bounding(w.Children().Boxes())
		b.CopyDimensionsOf(boundingBox)
	})

	w.Reactions().Add(
		reaction.NewMouseUpReaction(
			doodad.MouseMovedWithin[*reaction.MouseUpEvent](w),
			func(event *reaction.MouseUpEvent) {
				w.OnClick(w)
				event.StopPropagation()
			},
		),
		reaction.NewMouseMovedReaction(
			doodad.MouseMovedWithin[*reaction.MouseMovedEvent](w),
			func(event *reaction.MouseMovedEvent) {
				w.hovering()
				event.StopPropagation()
			},
		),
		reaction.NewMouseMovedReaction(
			doodad.MouseMovedOutside[*reaction.MouseMovedEvent](w),
			func(event *reaction.MouseMovedEvent) {
				w.stoppedHovering()
			},
		),
	)
}

func (w *Button) SetMessage(message string) {
	w.message = message
	w.onSetMessage(message)
}

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

// func (w *Button) Gestures(gesturer doodad.Gesturer) []func() {
// 	slog.Debug("Registering button gestures", "button", w, "gesturer", gesturer)

// 	withinBounds := func(x, y int) bool {
// 		return x >= w.Box.X() && x <= w.Box.X()+w.Box.Width() &&
// 			y >= w.Box.Y() && y <= w.Box.Y()+w.Box.Height()
// 	}

// 	return []func(){
// 		gesturer.OnMouseMove(func(x, y int) error {
// 			if withinBounds(x, y) {
// 				w.hovering()
// 				return doodad.ErrStopPropagation
// 			} else {
// 				w.stoppedHovering()
// 			}
// 			return nil
// 		}),
// 		gesturer.OnMouseUp(func(event doodad.MouseUpEvent) error {
// 			if event.Button != ebiten.MouseButtonLeft {
// 				return nil
// 			}
// 			x, y := event.X, event.Y
// 			if withinBounds(x, y) {
// 				if w.OnClick != nil {
// 					w.OnClick(w)
// 				}
// 				return doodad.ErrStopPropagation
// 			}
// 			return nil
// 		}),
// 	}
// }
