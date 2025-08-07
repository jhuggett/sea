package world_map

import (
	"design-library/doodad"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/colors"
)

func NewCursorDoodad(spaceTranslator SpaceTranslator) *CursorDoodad {
	cursorDoodad := &CursorDoodad{
		SpaceTranslator: spaceTranslator,
	}

	return cursorDoodad
}

type CursorDoodad struct {
	MouseX, MouseY int

	SpaceTranslator SpaceTranslator

	img *ebiten.Image

	doodad.Default
}

func (w *CursorDoodad) Teardown() error {
	return nil
}

func (w *CursorDoodad) Update() error {
	return nil
}

func (w *CursorDoodad) Draw(screen *ebiten.Image) {
	if !w.IsVisible() {
		return
	}

	op := &ebiten.DrawImageOptions{}

	x, y := w.SpaceTranslator.FromScreenToWorld(
		float64(w.MouseX),
		float64(w.MouseY),
	)

	x, y = w.SpaceTranslator.FromWorldToData(x, y)

	x, y = Floor(x, y)

	x, y = w.SpaceTranslator.FromDataToWorld(x, y)

	x, y = w.SpaceTranslator.FromWorldToScreen(x, y)

	xScale, yScale := w.SpaceTranslator.ScreenScale()

	op.GeoM.Translate(x, y)
	op.GeoM.Scale(xScale, yScale)

	screen.DrawImage(w.img, op)
}

func (w *CursorDoodad) Setup() {
	width, height := w.SpaceTranslator.TileSize()
	w.img = ebiten.NewImage(int(width), int(height))

	// Define the thickness of the outline
	outlineThickness := 2
	clr := colors.Primary

	// Draw a thicker outline of a square
	for i := 0; i < int(width); i++ {
		for t := 0; t < outlineThickness; t++ {
			if t < int(height) {
				w.img.Set(i, t, clr)               // Top edge
				w.img.Set(i, int(height)-1-t, clr) // Bottom edge
			}
		}
	}
	for j := 0; j < int(height); j++ {
		for t := 0; t < outlineThickness; t++ {
			if t < int(width) {
				w.img.Set(t, j, clr)              // Left edge
				w.img.Set(int(width)-1-t, j, clr) // Right edge
			}
		}
	}

	w.Reactions().Add(
		doodad.NewMouseMovedWithinReaction(w, func(mm doodad.MouseMoved) {
			w.MouseX = mm.X
			w.MouseY = mm.Y
		}),
	)
}
