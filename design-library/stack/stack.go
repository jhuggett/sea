package stack

import (
	"design-library/doodad"
	"design-library/position/box"
	"design-library/reaction"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Type string

const (
	Horizontal Type = "horizontal"
	Vertical   Type = "vertical"
)

type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

type Config struct {
	Type         Type
	Layout       *box.Box
	Children     *doodad.Children
	SpaceBetween int
	Padding      Padding

	BackgroundColor color.Color
}

func New(config Config) *Stack {
	// stack := &Stack{
	// 	Type: config.Type,
	// 	Default: doodad.Default{
	// 		Box:      config.Layout,
	// 		Children: &doodad.Children{},
	// 	},
	// 	SpaceBetween:    config.SpaceBetween,
	// 	Padding:         config.Padding,
	// 	BackgroundColor: config.BackgroundColor,
	// }

	// if config.Type != Horizontal && config.Type != Vertical {
	// 	config.Type = Vertical // Default to vertical if an invalid type is provided
	// }

	// if config.Children != nil {
	// 	for _, child := range config.Children.Doodads {
	// 		stack.AddChild(child)
	// 	}
	// }

	s := &Stack{
		Config: config,
	}

	if s.Config.Layout != nil {
		s.Box = s.Config.Layout
	}

	return s
}

type Stack struct {
	Config Config

	doodad.Default
	Type         Type
	SpaceBetween int
	Padding      Padding

	BackgroundColor color.Color

	background *ebiten.Image
}

/*
  - Warning: when setup is called, we reposition all children!!!
    This means that if you want to apply changes to the child of a stack, you
    should call setup first, otherwise the changes will be discarded.
*/
func (s *Stack) Setup() {
	if s.Config.Type != Horizontal && s.Type != Vertical {
		s.Type = Vertical // Default to vertical if an invalid type is provided
	} else {
		s.Type = s.Config.Type
	}

	if s.Config.SpaceBetween == 0 {
		s.SpaceBetween = 0
	} else {
		s.SpaceBetween = s.Config.SpaceBetween
	}

	if s.Config.Padding.Top == 0 && s.Config.Padding.Right == 0 &&
		s.Config.Padding.Bottom == 0 && s.Config.Padding.Left == 0 {
		s.Padding = Padding{} // Default padding
	} else {
		s.Padding = s.Config.Padding
	}

	if s.Config.BackgroundColor != nil {
		s.BackgroundColor = s.Config.BackgroundColor
	}

	// if s.Config.Children != nil {
	// 	for _, child := range s.Config.Children.Doodads {
	// 		s.AddChild(child)
	// 	}
	// }

	if s.Config.Children != nil {
		s.AddChild(s.Config.Children.Doodads...)
	}

	s.Box.Computed(func(b *box.Box) {
		// Adjust box for padding

		// Handle positioning - we use IncreaseWidth/Height instead of Move for padding
		// to ensure the stack is sized correctly regardless of position
		b.IncreaseWidth(s.Padding.Left + s.Padding.Right)
		b.IncreaseHeight(s.Padding.Top + s.Padding.Bottom)

		// If we're at position Y=0 (top of screen), only move right
		// Otherwise apply normal positioning
		// if b.Y() == 0 {
		// 	// b.MoveRight(s.Padding.Left)
		// } else {
		// 	// b.MoveLeft(s.Padding.Left).MoveUp(s.Padding.Top)
		// }

	})

	s.Children().Setup()

	var previousChild doodad.Doodad
	for _, child := range s.Children().Doodads {
		previousChildReferenceCopy := previousChild

		if s.Type == Horizontal {
			if previousChildReferenceCopy == nil {
				child.Layout().Computed(func(b *box.Box) {
					b.CopyPositionOf(s.Box).
						MoveDown(s.Padding.Top).
						MoveRight(s.Padding.Left)
				})
			} else {
				child.Layout().Computed(func(b *box.Box) {
					b.CopyPositionOf(previousChildReferenceCopy.Layout()).
						MoveRightOf(previousChildReferenceCopy.Layout()).
						MoveRight(s.SpaceBetween)
				})
			}
		} else if s.Type == Vertical {
			if previousChildReferenceCopy == nil {
				child.Layout().Computed(func(b *box.Box) {
					b.CopyPositionOf(s.Box).
						MoveDown(s.Padding.Top).
						MoveRight(s.Padding.Left)
				})
			} else {
				child.Layout().Computed(func(b *box.Box) {
					b.CopyPositionOf(previousChildReferenceCopy.Layout()).
						MoveBelow(previousChildReferenceCopy.Layout()).
						MoveDown(s.SpaceBetween)
				})
			}
		}
		previousChild = child

	}

	s.Reactions().Add(
		reaction.NewMouseMovedReaction(
			doodad.MouseMovedWithin[*reaction.MouseMovedEvent](s),
			func(event *reaction.MouseMovedEvent) {
				event.StopPropagation()
			},
		),
		reaction.NewMouseMovedReaction(
			doodad.MouseMovedOutside[*reaction.MouseMovedEvent](s),
			func(event *reaction.MouseMovedEvent) {
			},
		),
		reaction.NewMouseUpReaction(
			doodad.MouseMovedWithin[*reaction.MouseUpEvent](s),
			func(event *reaction.MouseUpEvent) {
				event.StopPropagation()
			},
		),
		reaction.NewMouseDragReaction(
			doodad.MouseMovedWithin[*reaction.OnMouseDragEvent](s),
			func(event *reaction.OnMouseDragEvent) {
				event.StopPropagation()
			},
		),
		reaction.NewMouseWheelReaction(
			doodad.MouseMovedWithin[*reaction.MouseWheelEvent](s),
			func(event *reaction.MouseWheelEvent) {
				event.StopPropagation()
			},
		),
	)

	if s.BackgroundColor != nil && s.Box.Width() > 0 && s.Box.Height() > 0 {
		s.background = ebiten.NewImage(s.Box.Width(), s.Box.Height())
		s.background.Fill(s.BackgroundColor)
	}
}

func (s *Stack) Draw(screen *ebiten.Image) {
	if !s.IsVisible() {
		return
	}

	if s.background != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(s.Box.X()), float64(s.Box.Y()))
		screen.DrawImage(s.background, op)

		if s.background.Bounds().Dx() != s.Box.Width() ||
			s.background.Bounds().Dy() != s.Box.Height() {
			// If the background size doesn't match the box size, we need to recreate it
			s.background = ebiten.NewImage(s.Box.Width(), s.Box.Height())
			s.background.Fill(s.BackgroundColor)
		}
	} else if s.BackgroundColor != nil && s.Box.Width() > 0 && s.Box.Height() > 0 {
		s.background = ebiten.NewImage(s.Box.Width(), s.Box.Height())
		s.background.Fill(s.BackgroundColor)
	}

	s.Children().Draw(screen)
}
