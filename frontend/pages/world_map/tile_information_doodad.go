package world_map

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/reaction"
	"design-library/stack"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/utils/space_translator"
)

func NewTileInformationDoodad(
	spaceTranslator space_translator.SpaceTranslator,
) *TileInformationDoodad {
	doodad := &TileInformationDoodad{
		SpaceTranslator: spaceTranslator,
	}

	return doodad
}

type TileInformationDoodad struct {
	space_translator.SpaceTranslator

	Hidden bool

	doodad.Default
}

func (w *TileInformationDoodad) Setup() {
	w.Hidden = true

	infoStack := stack.New(stack.Config{
		BackgroundColor: colors.Panel,
	})

	w.AddChild(infoStack)

	w.Reactions().Add()

	w.Reactions().Add(reaction.NewMouseUpReaction(
		doodad.MouseMovedWithin[*reaction.MouseUpEvent](w),
		func(event *reaction.MouseUpEvent) {
			if event.Button != ebiten.MouseButtonRight {
				return
			}
			event.StopPropagation()

			x, y := w.SpaceTranslator.FromWorldToData(w.SpaceTranslator.FromScreenToWorld(float64(event.X), float64(event.Y)))

			slog.Debug("WorldMapDoodad.OnClick", "x", x, "y", y)

			stackX, stackY := w.SpaceTranslator.FromWorldToScreen(w.SpaceTranslator.FromDataToWorld(x, y))
			scaleX, scaleY := w.SpaceTranslator.ScreenScale()

			infoStack.Layout().Computed(func(b *box.Box) {
				b.SetX(int(stackX * scaleX))
				b.SetY(int(stackY * scaleY))
				b.SetDimensions(200, 300)
			})

			infoStack.Show()
		}))

	closeButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			infoStack.Hide()
		},
		Config: label.Config{
			Message: "X",
			Layout: box.Computed(func(b *box.Box) {
				b.SetX(infoStack.Layout().X() + 10)
				b.SetY(infoStack.Layout().Y() + 10)
			}),
		},
	})
	infoStack.AddChild(closeButton)

	w.Children().Setup()

	infoStack.Hide()
}
