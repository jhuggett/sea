package settings

import (
	design_library "design-library"
	"design-library/button"
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
)

type SettingsPage struct {
	doodad.Default

	App *design_library.App
}

func New(app *design_library.App) *SettingsPage {
	return &SettingsPage{
		App: app,
	}
}

func (m *SettingsPage) Setup() {
	backButton := button.New(button.Config{
		OnClick: func(b *button.Button) {
			m.App.Pop()
		},
		Config: label.Config{
			Message: "Back",
		},
	})

	m.AddChild(backButton)

	titleLabel := label.New(label.Config{
		Message: "Settings",
	})

	m.AddChild(titleLabel)

	m.Children().Setup()

	titleLabel.Layout().Computed(func(b *box.Box) {
		b.MoveBelow(backButton.Box)
	})

	m.Box.Recalculate()
}
