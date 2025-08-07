package world_map

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"fmt"

	"github.com/jhuggett/frontend/colors"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/outbound"
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
}

func (w *TimeControlDoodad) Setup() {

	currentTimeLabel := label.New(label.Config{
		Message: "Time Control",
		Layout:  box.Zeroed(),
	})

	increaseSpeedButton := button.New(button.Config{
		OnClick: func(*button.Button) {
			t := w.Manager.LastTimeChangedReq.TicksPerSecond + 1
			_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
				SetTicksPerSecondTo: &t,
			})
			if err != nil {
				fmt.Println("Error increasing time:", err)
				return
			}
		},
		Config: label.Config{
			Message: "+",
		},
	})

	decreaseSpeedButton := button.New(button.Config{
		OnClick: func(*button.Button) {
			t := w.Manager.LastTimeChangedReq.TicksPerSecond - 1
			_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
				SetTicksPerSecondTo: &t,
			})
			if err != nil {
				fmt.Println("Error decreasing time:", err)
				return
			}
		},
		Config: label.Config{
			Message: "-",
		},
	})

	pauseResumeButton := button.New(button.Config{
		OnClick: func(*button.Button) {
			shouldResume := w.Manager.LastTimeChangedReq.IsPaused
			shouldPause := !shouldResume
			_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
				Pause:  shouldPause,
				Resume: shouldResume,
			})
			if err != nil {
				fmt.Println("Error toggling time pause:", err)
				return
			}
		},
		Config: label.Config{
			Message: "Pause/Resume",
		},
	})

	children := doodad.NewChildren(
		w,
		[]doodad.Doodad{
			currentTimeLabel,
			increaseSpeedButton,
			decreaseSpeedButton,
			pauseResumeButton,
		})

	panel := stack.New(stack.Config{
		Type: stack.Horizontal,
		Layout: box.Computed(func(b *box.Box) {
			boundingBox := box.Bounding(children.Boxes())
			b.CopyDimensionsOf(boundingBox).CenterHorizontallyWithin(w.Layout())
		}),
		Children:        children,
		BackgroundColor: colors.SemiTransparent(colors.Panel),
		Padding: stack.Padding{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
		SpaceBetween: 20,
	})

	w.Manager.OnTimeChangedCallback.Add(func(tcr outbound.TimeChangedReq) error {
		// currentTimeLabelDoodad.SetMessage(fmt.Sprintf("Time Control: %d", tcr.CurrentTick))
		// currentTickSpeedLabel.SetMessage(fmt.Sprintf("S: %d, P: %v", tcr.TicksPerSecond, tcr.IsPaused))

		currentTimeLabel.SetMessage(fmt.Sprintf("Day %d of Year %d", tcr.CurrentDay, tcr.CurrentYear))
		return nil
	})

	w.AddChild(panel)
	w.Children().Setup()

	// panelDoodad, err := panel.New(panel.Config{
	// 	Position: w.Position,
	// 	Gesturer: w.Gesturer,
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to create panel doodad: %w", err)
	// }
	// w.Children.Add(panelDoodad)

	// currentTimeLabelDoodad, err := label.New(label.Config{
	// 	Layout:  position.NewBox(position.BoxConfig{}),
	// 	Message: "Time Control",
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to create current time label doodad: %w", err)
	// }
	// w.Children.Add(currentTimeLabelDoodad)

	// currentTickSpeedLabel, err := label.New(label.Config{
	// 	// Position: doodad.Below2(currentTimeLabelDoodad),
	// 	Layout: position.
	// 		NewBox(position.BoxConfig{}).
	// 		MoveBelow(currentTimeLabelDoodad.Layout()),
	// 	Message: "Tick Speed",
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to create current tick speed label doodad: %w", err)
	// }
	// w.Children.Add(currentTickSpeedLabel)

	// w.Manager.OnTimeChangedCallback.Add(func(tcr outbound.TimeChangedReq) error {
	// 	currentTimeLabelDoodad.SetMessage(fmt.Sprintf("Time Control: %d", tcr.CurrentTick))
	// 	currentTickSpeedLabel.SetMessage(fmt.Sprintf("S: %d, P: %v", tcr.TicksPerSecond, tcr.IsPaused))
	// 	return nil
	// })

	// increaseTimeButton, err := button.New(button.Config{
	// 	Message: "Increase Tick Speed",
	// 	OnClick: func() {
	// 		t := w.Manager.LastTimeChangedReq.TicksPerSecond + 1
	// 		_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
	// 			SetTicksPerSecondTo: &t,
	// 		})
	// 		if err != nil {
	// 			fmt.Println("Error increasing time:", err)
	// 			return
	// 		}
	// 	},
	// 	Gesturer: w.Gesturer,
	// 	// Position: doodad.Below2(currentTickSpeedLabel),

	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to create increase time button: %w", err)
	// }
	// w.Children.Add(increaseTimeButton)

	// decreaseTimeButton, err := button.New(button.Config{
	// 	Message: "Decrease Tick Speed",
	// 	OnClick: func() {
	// 		t := w.Manager.LastTimeChangedReq.TicksPerSecond - 1
	// 		_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
	// 			SetTicksPerSecondTo: &t,
	// 		})
	// 		if err != nil {
	// 			fmt.Println("Error decreasing time:", err)
	// 			return
	// 		}
	// 	},
	// 	Gesturer: w.Gesturer,
	// 	// Position: doodad.Below2(increaseTimeButton),
	// 	Layout: position.NewBox(position.BoxConfig{}).
	// 		Copy(increaseTimeButton.Layout()).
	// 		MoveBelow(increaseTimeButton.Layout()),
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to create decrease time button: %w", err)
	// }
	// w.Children.Add(decreaseTimeButton)

	// pauseTimeButton, err := button.New(button.Config{
	// 	Message: "Pause/Resume Time",
	// 	OnClick: func() {
	// 		shouldResume := w.Manager.LastTimeChangedReq.IsPaused
	// 		shouldPause := !shouldResume
	// 		_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
	// 			Pause:  shouldPause,
	// 			Resume: shouldResume,
	// 		})
	// 		if err != nil {
	// 			fmt.Println("Error toggling time pause:", err)
	// 			return
	// 		}
	// 	},
	// 	Gesturer: w.Gesturer,
	// 	// Position: doodad.Below2(decreaseTimeButton),
	// 	Layout: position.NewBox(position.BoxConfig{}).
	// 		Copy(decreaseTimeButton.Layout()).
	// 		MoveBelow(decreaseTimeButton.Layout()),
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to create pause time button: %w", err)
	// }
	// w.Children.Add(pauseTimeButton)

}
