package stack

import (
	"design-library/config"
	"design-library/doodad"
	"design-library/position/box"
	"design-library/reaction"
	"log/slog"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type LayoutRule int

const (
	FitContents LayoutRule = iota
	Fill
)

type Config struct {
	Flow         config.Flow
	SpaceBetween int
	Padding      config.Padding

	BackgroundColor color.Color

	LayoutRule LayoutRule

	HorizontalAlignment config.HorizontalAlignment
	VerticalAlignment   config.VerticalAlignment
}

func New(config Config) *Stack {
	s := &Stack{
		Config: config,
	}

	return s
}

type Stack struct {
	Config Config

	doodad.Default
}

/*
  - Warning: when setup is called, we reposition all children!!!
    This means that if you want to apply changes to the child of a stack, you
    should call setup first, otherwise the changes will be discarded.
*/
func (s *Stack) Setup() {

	slog.Info("Setting up stack", "z", s.Z())

	var previousChild doodad.Doodad
	for _, child := range s.Children().Doodads {
		previousChildReferenceCopy := previousChild

		switch s.Config.Flow {
		case config.LeftToRight:
			if previousChildReferenceCopy == nil {
				child.Layout().Computed(func(b *box.Box) {
					switch s.Config.VerticalAlignment {
					case config.VerticalAlignmentTop:
						b.SetY(s.Box.Y() + s.Config.Padding.Top)
					case config.VerticalAlignmentCenter:
						b.SetY(s.Box.Y() + (s.Box.Height() / 2) - (b.Height() / 2))
					case config.VerticalAlignmentBottom:
						b.SetY(s.Box.Y() + s.Box.Height() - b.Height() - s.Config.Padding.Bottom)
					}

					b.SetX(s.Box.X() + s.Config.Padding.Left)
				})
			} else {
				child.Layout().Computed(func(b *box.Box) {
					switch s.Config.VerticalAlignment {
					case config.VerticalAlignmentTop:
						b.SetY(s.Box.Y() + s.Config.Padding.Top)
					case config.VerticalAlignmentCenter:
						b.SetY(s.Box.Y() + (s.Box.Height() / 2) - (b.Height() / 2))
					case config.VerticalAlignmentBottom:
						b.SetY(s.Box.Y() + s.Box.Height() - b.Height() - s.Config.Padding.Bottom)
					}

					b.SetX(previousChildReferenceCopy.Layout().X() +
						previousChildReferenceCopy.Layout().Width() +
						s.Config.SpaceBetween)
				})
			}
		case config.TopToBottom:
			if previousChildReferenceCopy == nil {
				child.Layout().Computed(func(b *box.Box) {
					switch s.Config.HorizontalAlignment {
					case config.HorizontalAlignmentLeft:
						b.SetX(s.Box.X() + s.Config.Padding.Left)
					case config.HorizontalAlignmentCenter:
						b.SetX(s.Box.X() + (s.Box.Width() / 2) - (b.Width() / 2))
					case config.HorizontalAlignmentRight:
						b.SetX(s.Box.X() + s.Box.Width() - b.Width() - s.Config.Padding.Right)
					}

					b.SetY(s.Box.Y() + s.Config.Padding.Top)
				})
			} else {
				child.Layout().Computed(func(b *box.Box) {
					switch s.Config.HorizontalAlignment {
					case config.HorizontalAlignmentLeft:
						b.SetX(s.Box.X() + s.Config.Padding.Left)
					case config.HorizontalAlignmentCenter:
						b.SetX(s.Box.X() + (s.Box.Width() / 2) - (b.Width() / 2))
					case config.HorizontalAlignmentRight:
						b.SetX(s.Box.X() + s.Box.Width() - b.Width() - s.Config.Padding.Right)
					}

					b.SetY(previousChildReferenceCopy.Layout().Y() +
						previousChildReferenceCopy.Layout().Height() +
						s.Config.SpaceBetween)
				})
			}
		}
		previousChild = child

	}

	s.Reactions().Add(
		reaction.NewMouseMovedReaction(
			doodad.MouseIsWithin[*reaction.MouseMovedEvent](s),
			func(event *reaction.MouseMovedEvent) {
				event.StopPropagation()
			},
		),
		reaction.NewMouseMovedReaction(
			doodad.MouseIsOutside[*reaction.MouseMovedEvent](s),
			func(event *reaction.MouseMovedEvent) {
			},
		),
		reaction.NewMouseUpReaction(
			doodad.MouseIsWithin[*reaction.MouseUpEvent](s),
			func(event *reaction.MouseUpEvent) {
				event.StopPropagation()
			},
		),
		reaction.NewMouseDragReaction(
			doodad.MouseIsWithin[*reaction.OnMouseDragEvent](s),
			func(event *reaction.OnMouseDragEvent) {
				event.StopPropagation()
			},
		),
		reaction.NewMouseWheelReaction(
			doodad.MouseIsWithin[*reaction.MouseWheelEvent](s),
			func(event *reaction.MouseWheelEvent) {
				event.StopPropagation()
			},
		),
	)
	s.Children().Setup()

	s.Box.Computed(func(b *box.Box) {
		switch s.Config.LayoutRule {
		case FitContents:
			bounding := box.Bounding(s.Children().Boxes())
			b.SetWidth(bounding.Width() + s.Config.Padding.Left + s.Config.Padding.Right)
			b.SetHeight(bounding.Height() + s.Config.Padding.Top + s.Config.Padding.Bottom)
		case Fill:
			// Do nothing, we fill the available space
		}
	})

	// if s.Config.FitContents {
	// 	s.Layout().Computed(func(b *box.Box) {
	// 		// bounding := box.Bounding(s.Children().Boxes())
	// 		b.SetWidth(paddedBounds.Width() + s.Config.Padding.Left + s.Config.Padding.Right)
	// 		b.SetHeight(paddedBounds.Height() + s.Config.Padding.Top + s.Config.Padding.Bottom)
	// 	})
	// }

	s.Layout().Recalculate()

	if s.Config.BackgroundColor != nil && s.Box.Width() > 0 && s.Box.Height() > 0 {
		background := ebiten.NewImage(s.Box.Width(), s.Box.Height())
		background.Fill(s.Config.BackgroundColor)

		s.SetCachedDraw(&doodad.CachedDraw{
			Image: background,
		})
	}
}
