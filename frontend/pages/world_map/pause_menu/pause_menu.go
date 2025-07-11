package pause_menu

import (
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/colors"
)

func NewPauseMenu(
	gesturer doodad.Gesturer,
) *PauseMenu {
	pauseMenu := &PauseMenu{
		Default: *doodad.NewDefault(),
	}

	pauseMenu.Gesturer = gesturer

	return pauseMenu
}

type PauseMenu struct {
	doodad.Default

	isHidden bool
}

func (w *PauseMenu) Show() {
	w.isHidden = false
}

func (w *PauseMenu) Hide() {
	w.isHidden = true
}

func (w *PauseMenu) IsHidden() bool {
	return w.isHidden
}

func (w *PauseMenu) Setup() {
	w.isHidden = true

	panelChildren := doodad.NewChildren(
		[]doodad.Doodad{
			label.New(label.Config{
				Message:  "Pause Menu",
				FontSize: 24,
			}),
			button.New(button.Config{
				OnClick: func() {
					w.Hide()
				},
				Gesturer: w.Gesturer,
				Config: label.Config{
					Message: "Resume",
				},
			}),
			button.New(button.Config{
				OnClick: func() {
					os.Exit(0)
				},
				Gesturer: w.Gesturer,
				Config: label.Config{
					Message: "Quit to Desktop",
				},
			}),
		},
	)

	panel := stack.New(stack.Config{
		Children: panelChildren,
		Type:     stack.Vertical,
		Layout: box.Computed(func(b *box.Box) *box.Box {
			boundingBox := box.Bounding(panelChildren.Boxes())

			return b.CopyDimensionsOf(boundingBox).CenterWithin(w.Box)
		}),

		BackgroundColor: colors.Panel,
		SpaceBetween:    10,
		Padding: stack.Padding{
			Top:    20,
			Bottom: 20,
			Left:   20,
			Right:  20,
		},
	})

	w.AddChild(panel)

	w.Children.Setup()
}

func (w *PauseMenu) Draw(screen *ebiten.Image) {
	if w.isHidden {
		return
	}

	// Apply blur effect to the background
	background := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	background.DrawImage(screen, nil)
	options := &ebiten.DrawImageOptions{}
	options.ColorM.Scale(0.2, 0.2, 0.2, 1) // Reduce brightness for blur effect
	screen.DrawImage(background, options)

	w.Children.Draw(screen)
}
