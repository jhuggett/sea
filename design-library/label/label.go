package label

import (
	"bytes"
	"design-library/doodad"
	"design-library/position/box"
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

type Config struct {
	BackgroundColor color.Color
	ForegroundColor color.Color
	FontSize        int
	Message         string
	Layout          *box.Box
	Padding         Padding
}

func New(config Config) *Label {
	label := &Label{
		Config: config,
	}

	return label
}

type Label struct {
	background *ebiten.Image
	fontSource *text.GoTextFaceSource

	BackgroundColor color.Color
	ForegroundColor color.Color

	message string

	FontSize int

	padding Padding

	doodad.Default

	Hidden bool

	Config Config

	OriginalBox *box.Box
}

func (w *Label) Hide() {
	w.Hidden = true
}

func (w *Label) Show() {
	w.Hidden = false
}

func (w *Label) Draw(screen *ebiten.Image) {
	if w.Hidden {
		return
	}

	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(float64(w.position().X), float64(w.position().Y))
	op.GeoM.Translate(float64(w.Box.X()), float64(w.Box.Y()))

	if w.background != nil {
		screen.DrawImage(w.background, op)
	}
}

func (w *Label) Setup() {
	if w.Config.Layout != nil {
		w.Box = w.Config.Layout
	}

	// if w.OriginalBox == nil { // this resets the box to the original box
	// 	w.OriginalBox = w.Box
	// }

	// w.Box = box.Computed(func(b *box.Box) {
	// 	b.Copy(w.OriginalBox)
	// })
	// w.OriginalBox.AddDependent(w.Box)

	if w.Config.BackgroundColor == nil {
		w.BackgroundColor = color.RGBA{0, 0, 0, 0}
	} else {
		w.BackgroundColor = w.Config.BackgroundColor
	}

	if w.Config.ForegroundColor == nil {
		w.ForegroundColor = color.White
	} else {
		w.ForegroundColor = w.Config.ForegroundColor
	}

	if w.Config.FontSize <= 0 {
		w.FontSize = 16
	} else {
		w.FontSize = w.Config.FontSize
	}

	if w.Config.Message == "" {
		w.Config.Message = "Label"
	} else {
		w.message = w.Config.Message
	}

	w.padding = w.Config.Padding

	var err error
	w.fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic("failed to load font: " + err.Error())
	}

	if w.message != "" {
		w.SetMessage(w.message)
	}
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

	// w.dimensions.Width = int(width)
	// w.dimensions.Height = int(height)
	// w.Layout().Width = int(width)
	// w.Layout().Height = int(height)

	w.Layout().Computed(func(b *box.Box) {
		b.SetDimensions(
			int(width)+w.padding.Left+w.padding.Right,
			int(height)+w.padding.Top+w.padding.Bottom,
		)
	})

	slog.Debug("(SetMessage) Updated Label dimensions", "width", w.Layout().Width, "height", w.Layout().Height)

	// slog.Debug("Label dimensions", "width", w.dimensions.Width, "height", w.dimensions.Height)

	w.background = ebiten.NewImage(w.Layout().Width(), w.Layout().Height())
	w.background.Fill(w.BackgroundColor)

	op := &text.DrawOptions{}
	colorScale := (&ebiten.ColorScale{})
	colorScale.ScaleWithColor(w.ForegroundColor)
	op.ColorScale = *colorScale
	op.GeoM.Translate(float64(w.padding.Left), float64(w.padding.Top))
	text.Draw(w.background, w.message, &text.GoTextFace{
		Source: textFace.Source,
		Size:   textFace.Size,
	}, op)
}
