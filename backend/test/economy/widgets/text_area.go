package widgets

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

type TextArea struct {
	Contents string
}

func (t *TextArea) Setup(root *widget.Container) {

	textAreaContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(0, 50),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(0, 0)),
		),
	)
	root.AddChild(textAreaContainer)

	textArea := newTextArea(t.Contents, widget.WidgetOpts.LayoutData(widget.GridLayoutData{
		MaxHeight: 220,
	}))
	textAreaContainer.AddChild(textArea)

	// c.AddChild(newSeparator(res, widget.RowLayoutData{
	// 	Stretch: true,
	// }))
}

func newTextArea(text string, widgetOpts ...widget.WidgetOpt) *widget.TextArea {
	face, _ := loadFont(14)

	return widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widgetOpts...)),
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Mask: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				// Set the track images
				&widget.SliderTrackImage{
					Idle:  image.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
					Hover: image.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
				},
				// Set the handle images
				&widget.ButtonImage{
					Idle:    image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Hover:   image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Pressed: image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				},
			),
		),
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		widget.TextAreaOpts.VerticalScrollMode(widget.PositionAtEnd),
		widget.TextAreaOpts.ProcessBBCode(true),
		widget.TextAreaOpts.FontFace(face),
		widget.TextAreaOpts.FontColor(color.NRGBA{R: 200, G: 100, B: 0, A: 255}),
		// widget.TextAreaOpts.TextPadding(res.textArea.entryPadding),
		widget.TextAreaOpts.Text(text),
	)
}
