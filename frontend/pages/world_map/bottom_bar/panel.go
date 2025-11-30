package bottom_bar

import (
	"design-library/config"
	"design-library/doodad"
	"design-library/stack"

	"github.com/jhuggett/frontend/colors"
)

func NewPanel() *Panel {
	panel := &Panel{}
	return panel
}

type Panel struct {
	doodad.Default

	Contents []doodad.Doodad
}

func (p *Panel) Setup() {
	panelStack := stack.New(stack.Config{
		Flow:            config.LeftToRight,
		BackgroundColor: colors.SemiTransparent(colors.Panel),
	})
	p.AddChild(panelStack)

	for _, content := range p.Contents {
		panelStack.AddChild(content)
	}

	p.Children().Setup()
}

func (p *Panel) SetContents(contents []doodad.Doodad) {
	p.Contents = contents
	p.Children().Clear()
	p.Setup()
}
