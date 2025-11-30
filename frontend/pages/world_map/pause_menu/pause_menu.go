package pause_menu

import (
	design_library "design-library"
	"design-library/button"
	"design-library/config"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jhuggett/frontend/colors"
)

func NewPauseMenu(
	app *design_library.App,
) *PauseMenu {
	pauseMenu := &PauseMenu{
		App: app,
	}
	return pauseMenu
}

type PauseMenu struct {
	App *design_library.App
	doodad.Default
}

func (w *PauseMenu) Setup() {

	// Root panel catches all mouse input
	rootPanel := stack.New(stack.Config{
		LayoutRule: stack.Fill,
	})

	w.AddChild(rootPanel)

	panel := stack.New(stack.Config{

		BackgroundColor: colors.Panel,
		SpaceBetween:    10,
		Padding: config.Padding{
			Top:    20,
			Bottom: 20,
			Left:   20,
			Right:  40,
		},
	})

	w.AddChild(panel)

	panel.AddChild(label.New(label.Config{
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
				w.App.PopToRoot()
			},
			Config: label.Config{
				Message: "Quit to Main Menu",
			},
		}),
		button.New(button.Config{
			OnClick: func(*button.Button) {
				os.Exit(0)
			},
			Config: label.Config{
				Message: "Quit to Desktop",
			},
		}))

	w.Children().Setup()

	panel.Box.Computed(func(b *box.Box) {
		b.CenterWithin(rootPanel.Layout())
	})

	// Apply blur effect to the background
	background := ebiten.NewImage(w.Layout().Width(), w.Layout().Height())
	colorScale := ebiten.ColorScale{}
	colorScale.Scale(0.2, 0.2, 0.2, 1) // Reduce brightness for blur effect
	options := &ebiten.DrawImageOptions{
		ColorScale: colorScale,
	}

	w.SetCachedDraw(
		&doodad.CachedDraw{
			Image: background,
			Op:    options,
		},
	)
}
