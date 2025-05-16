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
	label.SetPosition(func() Position {
		return Position{
			X: 0,
			Y: 0,
		}
	})
	label.FontSize = 16
	label.BackgroundColor = color.RGBA{0, 0, 0, 0}
	label.ForegroundColor = color.White

	return label
}

type Label struct {
	background *ebiten.Image
	fontSource *text.GoTextFaceSource

	BackgroundColor color.Color
	ForegroundColor color.Color

	message string

	position   func() Position
	dimensions Rectangle

	FontSize int
}

func (w *Label) Update() error {
	return nil
}

func (w *Label) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w.position().X), float64(w.position().Y))
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
		Size:   float64(w.FontSize),
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
	w.background.Fill(w.BackgroundColor)

	op := &text.DrawOptions{}
	colorScale := (&ebiten.ColorScale{})
	colorScale.ScaleWithColor(w.ForegroundColor)
	op.ColorScale = *colorScale
	op.GeoM.Translate(0, 0)
	text.Draw(w.background, w.message, &text.GoTextFace{
		Source: textFace.Source,
		Size:   textFace.Size,
	}, op)
}

func (w *Label) SetPosition(position func() Position) {
	w.position = position
}

func (w *Label) Position() Position {
	return w.position()
}

func (w *Label) Dimensions() Rectangle {
	return w.dimensions
}
