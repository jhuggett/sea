package tabs

import (
	"design-library/doodad"
)

type Tab struct {
	Name   string
	Doodad doodad.Doodad
}

// type TabsConfig struct {
// 	Gesturer doodad.Gesturer
// 	// Position   func() doodad.Position
// 	Dimensions doodad.Rectangle

// 	Tabs       map[string]doodad.Doodad
// 	InitialTab string
// }

// func New(config TabsConfig) (*Tabs, error) {
// 	tabs := &Tabs{}

// 	tabs.Gesturer = config.Gesturer
// 	// tabs.Position = config.Position
// 	tabs.Dimensions = config.Dimensions
// 	tabs.Tabs = config.Tabs
// 	tabs.currentTab = config.InitialTab

// 	if err := tabs.Setup(); err != nil {
// 		return nil, err
// 	}

// 	return tabs, nil
// }

// type Tabs struct {
// 	doodad.Default

// 	Tabs       map[string]doodad.Doodad
// 	currentTab string
// }

// func (t *Tabs) Setup() error {
// 	// Setup logic for the tabs can go here

// 	mainPanel, err := panel.New(panel.Config{
// 		Gesturer: t.Gesturer,
// 		// Position: t.Position,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create main panel: %w", err)
// 	}

// 	t.Children.Add(mainPanel)

// 	t.Children.Add(t.Tabs[t.currentTab])

// 	tabButtons := doodad.Children{
// 		Doodads: []doodad.Doodad{},
// 	}
// 	for name := range t.Tabs {
// 		button, err := button.New(button.Config{
// 			Gesturer: t.Gesturer,
// 			// Message:  name,
// 			OnClick: func() {
// 			},
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to create button for tab %s: %w", name, err)
// 		}
// 		tabButtons.Add(button)
// 	}

// 	stack, err := stack.New(stack.Config{
// 		Type: stack.Horizontal,
// 		// Position: t.Position,
// 		Children: tabButtons,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create stack: %w", err)
// 	}
// 	t.Children.Add(stack)

// 	return nil
// }
