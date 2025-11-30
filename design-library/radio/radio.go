package radio

import (
	"design-library/config"
	"design-library/doodad"
	"design-library/position/box"
	"design-library/stack"
)

type Option struct {
	Label    string
	OnSelect func()

	onSelect func()
}

func (o *Option) Select() {
	if o.onSelect != nil {
		o.onSelect()
	} else {
		panic("Radio option selected but no onSelect function defined")
	}
}

type Config struct {
	Flow         config.Flow
	Padding      config.Padding
	SpaceBetween int

	DefaultOptionIndex *int
	Options            []*Option

	UnselectedDoodad func(option *Option) doodad.Doodad
	SelectedDoodad   func(option *Option) doodad.Doodad

	OnSelect func(option *Option)
}

type Radio struct {
	doodad.Default

	CurrentIndex *int

	Config Config
}

func New(config Config) *Radio {
	radio := &Radio{
		Config: config,
	}

	radio.CurrentIndex = config.DefaultOptionIndex

	return radio
}

func (r *Radio) Setup() {
	mainStack := stack.New(stack.Config{
		Flow:         r.Config.Flow,
		Padding:      r.Config.Padding,
		SpaceBetween: r.Config.SpaceBetween,
		LayoutRule:   stack.FitContents,
	})

	r.AddChild(mainStack)

	for i, option := range r.Config.Options {

		option.onSelect = func() {
			r.CurrentIndex = &i
			doodad.ReSetup(r)
			if r.Config.OnSelect != nil {
				r.Config.OnSelect(option)
			}
			if option.OnSelect != nil {
				option.OnSelect()
			}
		}

		var optionDoodad doodad.Doodad
		if r.CurrentIndex != nil && *r.CurrentIndex == i {
			optionDoodad = r.Config.SelectedDoodad(option)
		} else {
			optionDoodad = r.Config.UnselectedDoodad(option)
		}

		mainStack.AddChild(optionDoodad)
	}

	r.Layout().Computed(func(b *box.Box) {
		b.CopyDimensionsOf(mainStack.Box)
	})

	r.Children().Setup()
}
