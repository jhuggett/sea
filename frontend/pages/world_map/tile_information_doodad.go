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
)

func NewTileInformationDoodad(
	spaceTranslator SpaceTranslator,
) *TileInformationDoodad {
	doodad := &TileInformationDoodad{
		SpaceTranslator: spaceTranslator,
	}

	return doodad
}

type TileInformationDoodad struct {
	SpaceTranslator

	Hidden bool

	HideTileCursor func()
	ShowTileCursor func()

	doodad.Default
}

func (w *TileInformationDoodad) Setup() {
	w.Hidden = true

	infoStack := stack.New(stack.Config{
		BackgroundColor: colors.Panel,
	})

	w.AddChild(infoStack)

	w.Reactions().Add()

	// w.Gesturer.OnMouseUp(func(event doodad.MouseUpEvent) error {
	// 	if !w.Hidden {
	// 		return doodad.ErrStopPropagation
	// 	}

	// 	if event.Button != ebiten.MouseButtonRight {
	// 		return nil
	// 	}

	// x, y := w.SpaceTranslator.FromWorldToData(w.SpaceTranslator.FromScreenToWorld(float64(event.X), float64(event.Y)))

	// // 	slog.Debug("WorldMapDoodad.OnClick", "x", x, "y", y)

	// stackX, stackY := w.SpaceTranslator.FromWorldToScreen(w.SpaceTranslator.FromDataToWorld(x, y))
	// scaleX, scaleY := w.SpaceTranslator.ScreenScale()

	// // 	// infoStack.SetPosition(int(stackX*scaleX), int(stackY*scaleY))

	// infoStack.Layout().Computed(func(b *box.Box) {
	// 	b.SetX(int(stackX * scaleX))
	// 	b.SetY(int(stackY * scaleY))
	// })

	// 	// infoStack.SetPosition(func() doodad.Position {

	// 	// 	return doodad.Position{
	// 	// 		X: int(stackX * scaleX),
	// 	// 		Y: int(stackY * scaleY),
	// 	// 	}
	// 	// })

	// 	w.Hidden = false
	// 	w.HideTileCursor()

	// 	return nil
	// })

	w.Reactions().Add(reaction.NewMouseUpReaction(
		doodad.MouseMovedWithin[*reaction.MouseUpEvent](w),
		func(event *reaction.MouseUpEvent) {
			// if !w.Hidden {
			// 	return doodad.ErrStopPropagation
			// }

			if event.Button != ebiten.MouseButtonRight {
				return
			}
			event.StopPropagation()

			x, y := w.SpaceTranslator.FromWorldToData(w.SpaceTranslator.FromScreenToWorld(float64(event.X), float64(event.Y)))

			slog.Debug("WorldMapDoodad.OnClick", "x", x, "y", y)

			stackX, stackY := w.SpaceTranslator.FromWorldToScreen(w.SpaceTranslator.FromDataToWorld(x, y))
			scaleX, scaleY := w.SpaceTranslator.ScreenScale()

			// infoStack.SetPosition(int(stackX*scaleX), int(stackY*scaleY))

			infoStack.Layout().Computed(func(b *box.Box) {
				b.SetX(int(stackX * scaleX))
				b.SetY(int(stackY * scaleY))
			})

			// infoStack.SetPosition(func() doodad.Position {

			// 	return doodad.Position{
			// 		X: int(stackX * scaleX),
			// 		Y: int(stackY * scaleY),
			// 	}
			// })

			w.Hidden = false
			// w.HideTileCursor()
			return
		}))

	// s.Reactions().Add(
	// 	reaction.NewMouseMovedReaction(
	// 		doodad.MouseMovedWithin[*reaction.MouseMovedEvent](s),
	// 		func(event *reaction.MouseMovedEvent) {
	// 			event.StopPropagation()
	// 		},
	// 	),
	// 	reaction.NewMouseMovedReaction(
	// 		doodad.MouseMovedOutside[*reaction.MouseMovedEvent](s),
	// 		func(event *reaction.MouseMovedEvent) {
	// 		},
	// 	),
	// 	reaction.NewMouseUpReaction(
	// 		doodad.MouseMovedWithin[*reaction.MouseUpEvent](s),
	// 		func(event *reaction.MouseUpEvent) {
	// 			event.StopPropagation()
	// 		},
	// 	),
	// 	reaction.NewMouseDragReaction(
	// 		doodad.MouseMovedWithin[*reaction.OnMouseDragEvent](s),
	// 		func(event *reaction.OnMouseDragEvent) {
	// 			event.StopPropagation()
	// 		},
	// 	),
	// 	reaction.NewMouseWheelReaction(
	// 		doodad.MouseMovedWithin[*reaction.MouseWheelEvent](s),
	// 		func(event *reaction.MouseWheelEvent) {
	// 			event.StopPropagation()
	// 		},
	// 	),
	// )

	// w.Gesturer.OnMouseMove(func(x, y int) error {
	// 	if !w.Hidden {
	// 		return doodad.ErrStopPropagation
	// 	}
	// 	return nil
	// })

	// closeButton, err := button.New(button.Config{
	// 	Message: "X",
	// 	OnClick: func() {
	// 		slog.Debug("WorldMapDoodad.CloseButton.OnClick")
	// 		w.Hidden = true
	// 		w.ShowTileCursor()
	// 	},
	// 	Gesturer: w.Gesturer,
	// 	Position: func() doodad.Position {
	// 		if infoStack.Position == nil {
	// 			return doodad.ZeroZero()
	// 		}
	// 		return doodad.Position{
	// 			X: infoStack.Position().X + 10,
	// 			Y: infoStack.Position().Y + 10,
	// 		}
	// 	},
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to create close button: %w", err)
	// }
	// w.Children.Add(closeButton)

	closeButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			w.Hide()
			w.ShowTileCursor()
		},
		Config: label.Config{
			Message: "X",
			Layout: box.Computed(func(b *box.Box) {
				b.SetX(infoStack.Layout().X() + 10)
				b.SetY(infoStack.Layout().Y() + 10)
			}),
		},
	})
	w.AddChild(closeButton)

	w.Children().Setup()

	infoStack.Hide()
}
