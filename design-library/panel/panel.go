package panel

import (
	"design-library/doodad"
	"design-library/position/box"

	"github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	Gesturer doodad.Gesturer
	Layout   *box.Box
	Children *doodad.Children
}

func New(config Config) *Panel {
	panel := &Panel{}

	panel.teardownCBs = []func(){}

	panel.Gesturer = config.Gesturer

	panel.Box = config.Layout

	if panel.Layout() == nil { // Don't trigger setup if a layout is not provided
		panel.Box = box.New(box.Config{})
	}

	panel.Children = config.Children

	return panel
}

type Panel struct {
	Gesturer doodad.Gesturer

	bg *ebiten.Image

	doodad.Default

	teardownCBs []func()
}

func (w *Panel) Draw(screen *ebiten.Image) {
	if w.Layout() == nil || w.bg == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w.Layout().X()), float64(w.Layout().Y()))
	screen.DrawImage(w.bg, op)

	w.Children.Draw(screen)
}

func (w *Panel) Setup() {
	width, height := w.Layout().Width(), w.Layout().Height()

	if width <= 0 || height <= 0 {
		return // No need to create a background image if the panel has no size
	}

	w.bg = ebiten.NewImage(width, height)
	// w.bg.Fill(color.RGBA{0, 120, 0, 10})

	w.teardownCBs = append(w.teardownCBs, w.Gesturer.OnMouseMove(func(x, y int) error {
		if x >= w.Layout().X() && x <= w.Layout().X()+w.Layout().Width() &&
			y >= w.Layout().Y() && y <= w.Layout().Y()+w.Layout().Height() {
			return doodad.ErrStopPropagation
		}
		return nil
	}))
}

func (w *Panel) Teardown() error {
	for _, cb := range w.teardownCBs {
		cb()
	}
	return nil
}
