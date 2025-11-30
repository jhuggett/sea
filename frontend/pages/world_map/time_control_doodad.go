package world_map

import (
	"design-library/button"
	"design-library/config"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/radio"
	"design-library/stack"
	"fmt"
	"image/color"

	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/outbound"
	"github.com/jhuggett/sea/timeline"
)

func NewTimeControlDoodad(
	manager *game.Manager,
) *TimeControlDoodad {
	timeControlDoodad := &TimeControlDoodad{
		Manager: manager,
	}

	return timeControlDoodad
}

type TimeControlDoodad struct {
	// Gesturer doodad.Gesturer

	// Children doodad.Children

	doodad.Default

	Manager *game.Manager

	// Position func() doodad.Position

	tcr outbound.TimeChangedReq
}

func (w *TimeControlDoodad) Setup() {

	panel := stack.New(stack.Config{
		BackgroundColor: colors.SemiTransparent(colors.Panel),
		Padding: config.Padding{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
		SpaceBetween: 20,
	})

	w.AddChild(panel)

	timeControlRadio := doodad.Stateful(w, "radio", func() *radio.Radio {
		return radio.New(radio.Config{
			Flow: config.LeftToRight,
			Options: []*radio.Option{
				{Label: "x0", OnSelect: func() {
					var t timeline.Tick = 0
					_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
						SetTicksPerSecondTo: &t,
					})
					if err != nil {
						fmt.Println("Error setting time to 0:", err)
						return
					}
				}},
				{Label: "x1", OnSelect: func() {
					var t timeline.Tick = 1
					_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
						SetTicksPerSecondTo: &t,
					})
					if err != nil {
						fmt.Println("Error setting time to 1:", err)
						return
					}
				}},
				{Label: "x3", OnSelect: func() {
					var t timeline.Tick = 3
					_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
						SetTicksPerSecondTo: &t,
					})
					if err != nil {
						fmt.Println("Error setting time to 3:", err)
						return
					}
				}},
				{Label: "x9", OnSelect: func() {
					var t timeline.Tick = 9
					_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
						SetTicksPerSecondTo: &t,
					})
					if err != nil {
						fmt.Println("Error setting time to 9:", err)
						return
					}
				}},
			},
			SelectedDoodad: func(option *radio.Option) doodad.Doodad {
				return label.New(label.Config{
					Message:         option.Label,
					BackgroundColor: color.RGBA{50, 120, 140, 250},
				})
			},
			UnselectedDoodad: func(option *radio.Option) doodad.Doodad {
				return button.New(button.Config{
					OnClick: func(b *button.Button) {
						option.Select()
					},
					Config: label.Config{
						Message: option.Label,
					},
				})
			},
		})
	})
	panel.AddChild(timeControlRadio)

	bottomRow := stack.New(stack.Config{
		Flow: config.LeftToRight,
	})

	panel.AddChild(bottomRow)

	currentTimeLabel := label.New(label.Config{
		Message: fmt.Sprintf("Day %d of Year %d", w.tcr.CurrentDay, w.tcr.CurrentYear),
	})

	bottomRow.AddChild(currentTimeLabel)

	w.DoOnTeardown(w.Manager.OnTimeChangedCallback.Register(func(tcr outbound.TimeChangedReq) error {
		w.tcr = tcr
		doodad.ReSetup(w)

		return nil
	}))

	w.Children().Setup()

	w.Layout().Recalculate()

	panel.Layout().Computed(func(b *box.Box) {
		b.CenterHorizontallyWithin(w.Box)
	})
}
