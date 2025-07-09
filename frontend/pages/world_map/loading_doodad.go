package world_map

// type LoadingDoodad struct {
// 	Children doodad.Children

// 	hidden bool

// 	setMessage func(message string)
// }

// func (w *LoadingDoodad) Teardown() error {
// 	w.Children.Teardown()
// 	return nil
// }

// func (w *LoadingDoodad) Update() error {
// 	return nil
// }

// func (w *LoadingDoodad) Draw(screen *ebiten.Image) {
// 	if w.hidden {
// 		return
// 	}
// 	w.Children.Draw(screen)
// }

// func (w *LoadingDoodad) Setup() error {

// 	labelDoodad, err := label.New(label.Config{
// 		Position: doodad.ZeroZero,
// 		Message:  "Loading...",
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create label doodad: %w", err)
// 	}
// 	w.Children.Add(labelDoodad)

// 	w.setMessage = labelDoodad.SetMessage

// 	w.hidden = true

// 	return nil
// }

// func (w *LoadingDoodad) Show() {
// 	w.hidden = false
// }

// func (w *LoadingDoodad) Hide() {
// 	w.hidden = true
// }

// func (w *LoadingDoodad) SetMessage(message string) {
// 	w.setMessage(message)
// }
