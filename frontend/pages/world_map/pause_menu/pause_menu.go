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

func NewPauseMenu() *PauseMenu {
	pauseMenu := &PauseMenu{}
	return pauseMenu
}

type PauseMenu struct {
	doodad.Default
}

func (w *PauseMenu) Setup() {
	panelChildren := doodad.NewChildren(
		w,
		[]doodad.Doodad{
			label.New(label.Config{
				Message:  "Pause Menu",
				FontSize: 24,
			}),
			button.New(button.Config{
				OnClick: func(*button.Button) {
					w.Hide()
				},
				Config: label.Config{
					Message: "Resume",
				},
			}),
			button.New(button.Config{
				OnClick: func(*button.Button) {
					os.Exit(0)
				},
				Config: label.Config{
					Message: "Quit to Desktop",
				},
			}),
		},
	)

	panel := stack.New(stack.Config{
		Children: panelChildren,
		Type:     stack.Vertical,
		Layout: box.Computed(func(b *box.Box) {
			boundingBox := box.Bounding(panelChildren.Boxes())
			b.CopyDimensionsOf(boundingBox).CenterWithin(w.Box)
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

	w.Children().Setup()
}

func (w *PauseMenu) Draw(screen *ebiten.Image) {
	if !w.IsVisible() {
		return
	}

	// Apply blur effect to the background
	background := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	background.DrawImage(screen, nil)
	options := &ebiten.DrawImageOptions{}
	options.ColorM.Scale(0.2, 0.2, 0.2, 1) // Reduce brightness for blur effect
	screen.DrawImage(background, options)

	w.Children().Draw(screen)
}
