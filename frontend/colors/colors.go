package colors

import "image/color"

var Primary = color.RGBA{R: 0x00, G: 0x80, B: 0xFF, A: 0xFF}
var Secondary = color.RGBA{R: 0xFF, G: 0xA5, B: 0x00, A: 0xFF}
var Text = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
var Panel = color.RGBA{
	R: 26 + 50,
	G: 13 + 50,
	B: 0 + 50,
	A: 200,
}

func SemiTransparent(c color.Color) color.Color {
	if c == nil {
		return nil
	}
	r, g, b, _ := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: 230,
	}
}
