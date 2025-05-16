package world_map

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/doodad"
	"github.com/jhuggett/frontend/game"
	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/outbound"
)

type TimeControlDoodad struct {
	Gesturer doodad.Gesturer

	Doodads []doodad.Doodad

	Manager *game.Manager

	Position func() doodad.Position
}

func (w *TimeControlDoodad) Update() error {
	return nil
}

func (w *TimeControlDoodad) Draw(screen *ebiten.Image) {
	for _, doodad := range w.Doodads {
		doodad.Draw(screen)
	}
}

func (w *TimeControlDoodad) Setup() error {

	panelDoodad := doodad.NewPanel()
	panelDoodad.SetPosition(func() doodad.Position {
		return w.Position()
	})
	w.Doodads = append(w.Doodads, panelDoodad)

	currentTimeLabelDoodad := doodad.NewLabel()
	currentTimeLabelDoodad.SetPosition(func() doodad.Position {
		return doodad.Position{
			X: panelDoodad.Position().X,
			Y: panelDoodad.Position().Y,
		}
	})
	currentTimeLabelDoodad.Setup()
	currentTimeLabelDoodad.SetMessage("Time Control")
	w.Doodads = append(w.Doodads, currentTimeLabelDoodad)

	currentTickSpeedLabel := doodad.NewLabel()

	currentTickSpeedLabel.SetPosition(func() doodad.Position {
		return doodad.Below(currentTimeLabelDoodad)
	})
	currentTickSpeedLabel.Setup()
	currentTickSpeedLabel.SetMessage("Time Control")
	w.Doodads = append(w.Doodads, currentTickSpeedLabel)

	w.Manager.OnTimeChangedCallback.Add(func(tcr outbound.TimeChangedReq) error {
		currentTimeLabelDoodad.SetMessage(fmt.Sprintf("Time Control: %d", tcr.CurrentTick))
		currentTickSpeedLabel.SetMessage(fmt.Sprintf("S: %d, P: %v", tcr.TicksPerSecond, tcr.IsPaused))
		return nil
	})

	increaseTimeButton := doodad.NewButton(
		"Increase",
		func() {
			t := w.Manager.LastTimeChangedReq.TicksPerSecond + 1
			_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
				SetTicksPerSecondTo: &t,
			})
			if err != nil {
				fmt.Println("Error increasing time:", err)
				return
			}
		},
		w.Gesturer,
	)
	increaseTimeButton.SetPosition(func() doodad.Position {
		return doodad.Below(currentTickSpeedLabel)
	})
	w.Doodads = append(w.Doodads, increaseTimeButton)

	decreaseTimeButton := doodad.NewButton(
		"Decrease",
		func() {
			t := w.Manager.LastTimeChangedReq.TicksPerSecond - 1

			_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
				SetTicksPerSecondTo: &t,
			})
			if err != nil {
				fmt.Println("Error increasing time:", err)
				return
			}
		},
		w.Gesturer,
	)
	decreaseTimeButton.SetPosition(func() doodad.Position {
		return doodad.Below(increaseTimeButton)
	})
	w.Doodads = append(w.Doodads, decreaseTimeButton)

	pauseTimeButton := doodad.NewButton(
		"Pause",
		func() {
			shouldResume := w.Manager.LastTimeChangedReq.IsPaused
			shouldPause := !shouldResume

			_, err := w.Manager.ControlTime(inbound.ControlTimeReq{
				Pause:  shouldPause,
				Resume: shouldResume,
			})
			if err != nil {
				fmt.Println("Error increasing time:", err)
				return
			}
		},
		w.Gesturer,
	)
	pauseTimeButton.SetPosition(func() doodad.Position {
		return doodad.Below(decreaseTimeButton)
	})
	w.Doodads = append(w.Doodads, pauseTimeButton)

	return nil
}
