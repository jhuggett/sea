package doodad

import (
	"bytes"
	"image/color"
	"log"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

func NewLabel() *Label {
	label := &Label{}
	label.Setup()
	return label
}

type Label struct {
	background *ebiten.Image
	fontSource *text.GoTextFaceSource

	message string

	position   Position
	dimensions Rectangle
}

func (w *Label) Update() error {
	return nil
}

func (w *Label) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w.position.X), float64(w.position.Y))
	screen.DrawImage(w.background, op)
}

func (w *Label) Setup() error {
	var err error
	w.fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (w *Label) SetMessage(message string) {
	w.message = message

	textFace := &text.GoTextFace{
		Source: w.fontSource,
		Size:   48,
	}

	width, height := text.Measure(
		w.message,
		textFace,
		0,
	)

	w.dimensions.Width = int(width)
	w.dimensions.Height = int(height)

	slog.Debug("Label dimensions", "width", w.dimensions.Width, "height", w.dimensions.Height)

	w.background = ebiten.NewImage(int(width), int(height))
	w.background.Fill(color.RGBA{
		R: 56,
		G: 0,
		B: 23,
		A: 255,
	})

	op := &text.DrawOptions{}
	op.GeoM.Translate(0, 0)
	text.Draw(w.background, w.message, &text.GoTextFace{
		Source: textFace.Source,
		Size:   textFace.Size,
	}, op)
}

func (w *Label) SetPosition(position Position) {
	w.position = position
}

func (w *Label) Position() Position {
	return w.position
}

func (w *Label) Dimensions() Rectangle {
	return w.dimensions
}
