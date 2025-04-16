package doodad

type PositionedRectangle interface {
	Positioned
	Rectangular
}

func Below(doodad PositionedRectangle) Position {
	return Position{
		X: doodad.Position().X,
		Y: doodad.Position().Y + doodad.Dimensions().Height,
	}
}
