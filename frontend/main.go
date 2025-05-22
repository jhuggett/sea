package main

import "log/slog"

func main() {
	game := &Game{}
	game.Start()

	slog.Info("Game over.")
}
