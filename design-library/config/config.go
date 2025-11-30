package config

type Flow int

const (
	TopToBottom Flow = iota
	LeftToRight
)

type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

func EqualPadding(amount int) Padding {
	return Padding{
		Top:    amount,
		Right:  amount,
		Bottom: amount,
		Left:   amount,
	}
}

func SymmetricPadding(vertical int, horizontal int) Padding {
	return Padding{
		Top:    vertical,
		Right:  horizontal,
		Bottom: vertical,
		Left:   horizontal,
	}
}

type VerticalAlignment int

const (
	VerticalAlignmentTop VerticalAlignment = iota
	VerticalAlignmentCenter
	VerticalAlignmentBottom
)

type HorizontalAlignment int

const (
	HorizontalAlignmentLeft HorizontalAlignment = iota
	HorizontalAlignmentCenter
	HorizontalAlignmentRight
)
