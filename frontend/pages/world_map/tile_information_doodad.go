package world_map

// type TileInformationDoodad struct {
// 	SpaceTranslator
// 	Gesturer doodad.Gesturer

// 	Children doodad.Children

// 	Hidden bool

// 	HideTileCursor func()
// 	ShowTileCursor func()
// }

// func (w *TileInformationDoodad) Teardown() error {
// 	w.Children.Teardown()
// 	return nil
// }

// func (w *TileInformationDoodad) Update() error {
// 	return nil
// }

// func (w *TileInformationDoodad) Draw(screen *ebiten.Image) {
// 	if w.Hidden {
// 		return
// 	}

// 	w.Children.Draw(screen)
// }

// func (w *TileInformationDoodad) Setup() error {

// 	w.Hidden = true

// 	infoPanel, err := panel.New(panel.Config{
// 		Position: doodad.ZeroZero,
// 		Gesturer: w.Gesturer,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create info panel: %w", err)
// 	}
// 	w.Children.Add(infoPanel)

// 	w.Gesturer.OnMouseUp(func(event doodad.MouseUpEvent) error {
// 		if !w.Hidden {
// 			return doodad.ErrStopPropagation
// 		}

// 		if event.Button != ebiten.MouseButtonRight {
// 			return nil
// 		}

// 		x, y := w.SpaceTranslator.FromWorldToData(w.SpaceTranslator.FromScreenToWorld(float64(event.X), float64(event.Y)))

// 		slog.Debug("WorldMapDoodad.OnClick", "x", x, "y", y)

// 		panelX, panelY := w.SpaceTranslator.FromWorldToScreen(w.SpaceTranslator.FromDataToWorld(x, y))
// 		scaleX, scaleY := w.SpaceTranslator.ScreenScale()

// 		infoPanel.SetPosition(int(panelX*scaleX), int(panelY*scaleY))

// 		// infoPanel.SetPosition(func() doodad.Position {

// 		// 	return doodad.Position{
// 		// 		X: int(panelX * scaleX),
// 		// 		Y: int(panelY * scaleY),
// 		// 	}
// 		// })

// 		w.Hidden = false
// 		w.HideTileCursor()

// 		return nil
// 	})

// 	w.Gesturer.OnMouseMove(func(x, y int) error {
// 		if !w.Hidden {
// 			return doodad.ErrStopPropagation
// 		}
// 		return nil
// 	})

// 	closeButton, err := button.New(button.Config{
// 		Message: "X",
// 		OnClick: func() {
// 			slog.Debug("WorldMapDoodad.CloseButton.OnClick")
// 			w.Hidden = true
// 			w.ShowTileCursor()
// 		},
// 		Gesturer: w.Gesturer,
// 		Position: func() doodad.Position {
// 			if infoPanel.Position == nil {
// 				return doodad.ZeroZero()
// 			}
// 			return doodad.Position{
// 				X: infoPanel.Position().X + 10,
// 				Y: infoPanel.Position().Y + 10,
// 			}
// 		},
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create close button: %w", err)
// 	}
// 	w.Children.Add(closeButton)

// 	return nil
// }
