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
}

func New(config Config) *Stack {
	stack := &Stack{
		Type: config.Type,
		Default: doodad.Default{
			Box:      config.Layout,
			Children: &doodad.Children{},
		},
		SpaceBetween: config.SpaceBetween,
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

	background *ebiten.Image
}

func (s *Stack) Setup() {
	var previousChild doodad.Doodad
	for _, child := range s.Children.Doodads {
		previousChildReferenceCopy := previousChild

		child.Setup()

		if previousChildReferenceCopy == nil {
			child.Layout().Computed(func(b *box.Box) *box.Box {
				return b.CopyPositionOf(s.Box)
			})
		} else {
			child.Layout().Computed(func(b *box.Box) *box.Box {
				return b.CopyPositionOf(previousChildReferenceCopy.Layout()).
					MoveBelow(previousChildReferenceCopy.Layout()).
					MoveDown(s.SpaceBetween)
			})
		}
		previousChild = child
	}

	s.background = ebiten.NewImage(s.Box.Width(), s.Box.Height())
	s.background.Fill(color.RGBA{0, 120, 100, 200})
}

func (s *Stack) Draw(screen *ebiten.Image) {
	if s.background != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(s.Box.X()), float64(s.Box.Y()))
		screen.DrawImage(s.background, op)
	}

	s.Children.Draw(screen)
}
