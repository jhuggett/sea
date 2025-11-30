package label

import (
	"bytes"
	"design-library/doodad"
	"design-library/position/box"
	"image/color"

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
	if config.ForegroundColor == nil {
		config.ForegroundColor = color.White
	}

	if config.FontSize == 0 {
		config.FontSize = 16
	}

	if config.Message == "" {
		config.Message = "Missing Label Message"
	}

	label := &Label{
		Config: config,
	}

	return label
}

type Label struct {
	fontSource *text.GoTextFaceSource

	doodad.Default

	Config Config

	OriginalBox *box.Box
}

// func (w *Label) Draw(screen *ebiten.Image) {
// 	if w.Hidden {
// 		return
// 	}

// 	op := &ebiten.DrawImageOptions{}
// 	// op.GeoM.Translate(float64(w.position().X), float64(w.position().Y))
// 	op.GeoM.Translate(float64(w.Box.X()), float64(w.Box.Y()))

// 	if w.background != nil {
// 		screen.DrawImage(w.background, op)
// 	}
// }

func (w *Label) Setup() {
	if w.fontSource == nil {
		var err error
		w.fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
		if err != nil {
			panic("failed to load font: " + err.Error())
		}
	}

	textFace := &text.GoTextFace{
		Source: w.fontSource,
		Size:   float64(w.Config.FontSize),
	}

	width, height := text.Measure(
		w.Config.Message,
		textFace,
		0,
	)

	w.Layout().Computed(func(b *box.Box) {
		b.SetDimensions(
			int(width)+w.Config.Padding.Left+w.Config.Padding.Right,
			int(height)+w.Config.Padding.Top+w.Config.Padding.Bottom,
		)
	})

	img := ebiten.NewImage(w.Layout().Width(), w.Layout().Height())

	if w.Config.BackgroundColor != nil {
		img.Fill(w.Config.BackgroundColor)
	}

	op := &text.DrawOptions{}
	colorScale := (&ebiten.ColorScale{})

	colorScale.ScaleWithColor(w.Config.ForegroundColor)
	op.ColorScale = *colorScale
	op.GeoM.Translate(float64(w.Config.Padding.Left), float64(w.Config.Padding.Top))
	text.Draw(img, w.Config.Message, &text.GoTextFace{
		Source: textFace.Source,
		Size:   textFace.Size,
	}, op)

	w.SetCachedDraw(&doodad.CachedDraw{
		Image: img,
	})
}

// func (w *Label) SetMessage(message string) {
// 	w.message = message

// 	textFace := &text.GoTextFace{
// 		Source: w.fontSource,
// 		Size:   float64(w.FontSize),
// 	}

// 	width, height := text.Measure(
// 		w.message,
// 		textFace,
// 		0,
// 	)

// 	// w.dimensions.Width = int(width)
// 	// w.dimensions.Height = int(height)
// 	// w.Layout().Width = int(width)
// 	// w.Layout().Height = int(height)

// 	w.Layout().Computed(func(b *box.Box) {
// 		b.SetDimensions(
// 			int(width)+w.padding.Left+w.padding.Right,
// 			int(height)+w.padding.Top+w.padding.Bottom,
// 		)
// 	})

// 	slog.Debug("(SetMessage) Updated Label dimensions", "width", w.Layout().Width, "height", w.Layout().Height)

// 	// slog.Debug("Label dimensions", "width", w.dimensions.Width, "height", w.dimensions.Height)

// 	w.background = ebiten.NewImage(w.Layout().Width(), w.Layout().Height())
// 	w.background.Fill(w.BackgroundColor)

// 	op := &text.DrawOptions{}
// 	colorScale := (&ebiten.ColorScale{})
// 	colorScale.ScaleWithColor(w.ForegroundColor)
// 	op.ColorScale = *colorScale
// 	op.GeoM.Translate(float64(w.padding.Left), float64(w.padding.Top))
// 	text.Draw(w.background, w.message, &text.GoTextFace{
// 		Source: textFace.Source,
// 		Size:   textFace.Size,
// 	}, op)
// }
