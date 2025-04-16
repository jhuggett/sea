package world_map

func (w *WorldMapPage) Update() error {
	for _, d := range w.Doodads {
		d.Update()
	}

	w.Gesturer.Update()

	return nil
}
