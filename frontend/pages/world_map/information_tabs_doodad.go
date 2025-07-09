package world_map

// type InformationTabsDoodad struct {
// 	Gesturer doodad.Gesturer

// 	Children doodad.Children

// 	Manager *game.Manager

// 	Position func() doodad.Position
// }

// func (w *InformationTabsDoodad) Teardown() error {
// 	w.Children.Teardown()
// 	return nil
// }

// func (w *InformationTabsDoodad) Update() error {
// 	return nil
// }

// func (w *InformationTabsDoodad) Draw(screen *ebiten.Image) {
// 	w.Children.Draw(screen)
// }

// func (w *InformationTabsDoodad) Setup() error {
// 	shipTabDoodad := &ShipTabDoodad{
// 		Gesturer: w.Gesturer,
// 		Children: doodad.Children{},
// 		Manager:  w.Manager,
// 		Position: w.Position,
// 	}
// 	if err := shipTabDoodad.Setup(); err != nil {
// 		return fmt.Errorf("failed to setup ship tab doodad: %w", err)
// 	}

// 	inventoryTabDoddad := &InventoryTabDoodad{
// 		Default: doodad.Default{
// 			Gesturer: w.Gesturer,
// 			Children: doodad.Children{},
// 			Position: w.Position,
// 		},
// 	}
// 	if err := inventoryTabDoddad.Setup(); err != nil {
// 		return fmt.Errorf("failed to setup inventory tab doodad: %w", err)
// 	}

// 	tabsDoodad, err := tabs.New(tabs.TabsConfig{
// 		Gesturer: w.Gesturer,
// 		Position: w.Position,
// 		Dimensions: doodad.Rectangle{
// 			Width:  200,
// 			Height: 300,
// 		},
// 		Tabs: map[string]doodad.Doodad{
// 			"Ship":      shipTabDoodad,
// 			"Inventory": inventoryTabDoddad,
// 		},
// 		InitialTab: "Ship",
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create tabs doodad: %w", err)
// 	}
// 	w.Children.Add(tabsDoodad)

// 	return nil
// }

// type InventoryTabDoodad struct {
// 	doodad.Default
// 	Manager *game.Manager
// }

// func (i *InventoryTabDoodad) Setup() error {
// 	return nil
// }

// type ShipTabDoodad struct {
// 	Gesturer doodad.Gesturer

// 	Children doodad.Children

// 	Manager *game.Manager

// 	Position func() doodad.Position
// }

// func (w *ShipTabDoodad) Teardown() error {
// 	w.Children.Teardown()
// 	return nil
// }

// func (w *ShipTabDoodad) Update() error {
// 	return nil
// }

// func (w *ShipTabDoodad) Draw(screen *ebiten.Image) {
// 	w.Children.Draw(screen)
// }

// func (w *ShipTabDoodad) Setup() error {
// 	currentTimeLabelDoodad, err := label.New(label.Config{
// 		Position: w.Position,
// 		Message:  "Information Tabs",
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create current time label doodad: %w", err)
// 	}
// 	w.Children.Add(currentTimeLabelDoodad)

// 	stateOfRepairLabel, err := label.New(label.Config{
// 		Position: doodad.Below2(currentTimeLabelDoodad),
// 		Message:  "State of Repair: ---",
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create state of repair label: %w", err)
// 	}
// 	w.Children.Add(stateOfRepairLabel)

// 	repairButton, err := button.New(button.Config{
// 		Message: "Repair Ship",
// 		OnClick: func() {
// 			_, err := w.Manager.PlayerShip.Repair()
// 			if err != nil {
// 				fmt.Printf("Failed to repair ship: %v\n", err)
// 				return
// 			}
// 		},
// 		Gesturer: w.Gesturer,
// 		Position: doodad.Below2(stateOfRepairLabel),
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create repair button: %w", err)
// 	}
// 	w.Children.Add(repairButton)

// 	estimatedSailingSpeedLabel, err := label.New(label.Config{
// 		Position: doodad.Below2(repairButton),
// 		Message:  "Estimated Sailing Speed: ---",
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create estimated sailing speed label: %w", err)
// 	}
// 	w.Children.Add(estimatedSailingSpeedLabel)

// 	w.Manager.OnShipChangedCallback.Register(func(scr outbound.ShipChangedReq) error {
// 		stateOfRepairLabel.SetMessage(
// 			fmt.Sprintf(
// 				"State of Repair: %f",
// 				scr.StateOfRepair,
// 			),
// 		)

// 		estimatedSailingSpeedLabel.SetMessage(
// 			fmt.Sprintf(
// 				"Speed: %f",
// 				scr.EstimatedSailingSpeed,
// 			),
// 		)

// 		return nil
// 	})

// 	return nil
// }
