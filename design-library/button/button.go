package button

import (
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"image/color"
	"log/slog"

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
		button.Box = box.New(box.Config{})
	}

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
	})
	w.AddChild(hoveredLabel)

	w.Children().Setup()

	w.Box.Computed(func(b *box.Box) {
		boundingBox := box.Bounding(w.Children().Boxes())
		b.CopyDimensionsOf(boundingBox)
	})

	w.showHoveringLabel = func() {
		hoveredLabel.Show()
		nonHoveredLabel.Hide()
	}
	w.showNonHoveredLabel = func() {
		hoveredLabel.Hide()
		nonHoveredLabel.Show()
	}

	w.hovered = true // so that it'll trigger the non-hovered event first

	w.onSetMessage = func(message string) {
		nonHoveredLabel.SetMessage(message)
		hoveredLabel.SetMessage(message)
	}
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

func (w *Button) Gestures(gesturer doodad.Gesturer) []func() {
	slog.Debug("Registering button gestures", "button", w, "gesturer", gesturer)

	withinBounds := func(x, y int) bool {
		return x >= w.Box.X() && x <= w.Box.X()+w.Box.Width() &&
			y >= w.Box.Y() && y <= w.Box.Y()+w.Box.Height()
	}

	return []func(){
		gesturer.OnMouseMove(func(x, y int) error {
			if withinBounds(x, y) {
				w.hovering()
				// return callback.ErrStopPropagation
			} else {
				w.stoppedHovering()
			}
			return nil
		}),
		gesturer.OnMouseUp(func(event doodad.MouseUpEvent) error {
			if event.Button != ebiten.MouseButtonLeft {
				return nil
			}
			x, y := event.X, event.Y
			if withinBounds(x, y) {
				if w.OnClick != nil {
					w.OnClick(w)
				}
				return doodad.ErrStopPropagation
			}
			return nil
		}),
	}
}
