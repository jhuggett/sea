package stack

import (
	"design-library/doodad"
	"design-library/position/box"
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
	stack := &Stack{
		Type: config.Type,
		Default: doodad.Default{
			Box:      config.Layout,
			Children: &doodad.Children{},
		},
		SpaceBetween: config.SpaceBetween,
		Padding:      config.Padding,
	}

	if config.Type != Horizontal && config.Type != Vertical {
		config.Type = Vertical // Default to vertical if an invalid type is provided
	}

	for _, child := range config.Children.Doodads {
		stack.AddChild(child)
	}

	return stack
}

type Stack struct {
	doodad.Default
	Type         Type
	SpaceBetween int
	Padding      Padding

	BackgroundColor color.Color

	background *ebiten.Image
}

func (s *Stack) Setup() {

	s.Box.Computed(func(b *box.Box) *box.Box {
		return b.MoveLeft(s.Padding.Left).
			MoveUp(s.Padding.Top).
			IncreaseWidth(s.Padding.Left + s.Padding.Right).
			IncreaseHeight(s.Padding.Top + s.Padding.Bottom)
	})

	var previousChild doodad.Doodad
	for _, child := range s.Children.Doodads {
		previousChildReferenceCopy := previousChild

		child.Setup()

		if s.Type == Horizontal {
			if previousChildReferenceCopy == nil {
				child.Layout().Computed(func(b *box.Box) *box.Box {
					return b.CopyPositionOf(s.Box).
						MoveDown(s.Padding.Top).
						MoveRight(s.Padding.Left)
				})
			} else {
				child.Layout().Computed(func(b *box.Box) *box.Box {
					return b.CopyPositionOf(previousChildReferenceCopy.Layout()).
						MoveRightOf(previousChildReferenceCopy.Layout()).
						MoveRight(s.SpaceBetween)
				})
			}
		} else if s.Type == Vertical {
			if previousChildReferenceCopy == nil {
				child.Layout().Computed(func(b *box.Box) *box.Box {
					return b.CopyPositionOf(s.Box).
						MoveDown(s.Padding.Top).
						MoveRight(s.Padding.Left)
				})
			} else {
				child.Layout().Computed(func(b *box.Box) *box.Box {
					return b.CopyPositionOf(previousChildReferenceCopy.Layout()).
						MoveBelow(previousChildReferenceCopy.Layout()).
						MoveDown(s.SpaceBetween)
				})
			}
		}
		previousChild = child
	}

	if s.BackgroundColor != nil {
		s.background = ebiten.NewImage(s.Box.Width(), s.Box.Height())
		s.background.Fill(s.BackgroundColor)
	}
}

func (s *Stack) Draw(screen *ebiten.Image) {
	if s.background != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(s.Box.X()), float64(s.Box.Y()))
		screen.DrawImage(s.background, op)
	}

	s.Children.Draw(screen)
}
