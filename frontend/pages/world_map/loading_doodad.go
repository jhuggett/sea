package world_map

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/doodad"
)

type LoadingDoodad struct {
	label *doodad.Label

	hidden bool
}

func (w *LoadingDoodad) Update() error {
	return nil
}

func (w *LoadingDoodad) Draw(screen *ebiten.Image) {

	if w.hidden {
		return
	}

	w.label.Draw(screen)
}

func (w *LoadingDoodad) Setup() error {

	w.label = doodad.NewLabel()

	w.hidden = true

	return nil
}

func (w *LoadingDoodad) Show() {
	w.hidden = false
}

func (w *LoadingDoodad) Hide() {
	w.hidden = true
}

func (w *LoadingDoodad) SetMessage(message string) {
	w.label.SetMessage(message)
}
