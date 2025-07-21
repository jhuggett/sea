package box

import (
	"fmt"
	"log/slog"
)

type Box struct {
	x int
	y int

	width  int
	height int

	calculationSteps   []func(*Box)
	needsRecalculation bool

	dependents []*Box

	recalculationCount int
}

func (b *Box) ZeroOut() *Box {
	b.x = 0
	b.y = 0
	b.width = 0
	b.height = 0

	return b
}

func (b *Box) FlagNeedsRecalculation() *Box {
	b.needsRecalculation = true

	for _, dependent := range b.dependents {
		dependent.FlagNeedsRecalculation()
	}

	return b
}

func (b *Box) Recalculate() *Box {
	b.needsRecalculation = true
	b.recalculateIfNeeded()
	for _, dependent := range b.dependents {
		dependent.Recalculate()
	}
	return b
}

func (b *Box) AddDependent(dependent *Box) {
	slog.Debug("Adding dependent Box", "dependent", fmt.Sprintf("%p", dependent), "to", fmt.Sprintf("%p", b))

	if b == dependent {
		// panic("A Box cannot depend on itself")
		slog.Warn("A Box cannot depend on itself, ignoring dependency", "box", b)
		return
	}
	if dependent == nil {
		panic("Cannot add a nil Box as a dependent")
	}
	for _, existingDependent := range b.dependents {
		if existingDependent == dependent {
			panic("A Box cannot depend on the same Box multiple times")
		}
	}
	for _, existingDependent := range dependent.dependents {
		if existingDependent == b {
			panic("A Box cannot depend on a Box that already depends on it")
		}
	}

	b.dependents = append(b.dependents, dependent)
}

func (b *Box) ClearDependents() {
	b.dependents = []*Box{}
}

type Config struct {
	X      int
	Y      int
	Width  int
	Height int
}

func New(config Config) *Box {
	return &Box{
		x:                config.X,
		y:                config.Y,
		width:            config.Width,
		height:           config.Height,
		dependents:       []*Box{},
		calculationSteps: []func(*Box){},
	}
}

func (b *Box) Computed(calculateFn func(*Box)) *Box {
	b.calculationSteps = append(b.calculationSteps, calculateFn)
	b.FlagNeedsRecalculation()
	return b
}

func (b *Box) MoveBelow(other *Box) *Box {
	b.y = other.Y() + other.Height()

	return b
}

func (b *Box) MoveAbove(other *Box) *Box {
	b.y = other.Y() - b.height
	return b
}

func (b *Box) MoveLeftOf(other *Box) *Box {
	b.x = other.X() - b.width

	return b
}

func (b *Box) MoveRightOf(other *Box) *Box {
	b.x = other.X() + other.Width()

	return b
}

func (b *Box) CopyPositionOf(other *Box) *Box {
	b.x = other.X()
	b.y = other.Y()

	return b
}

func (b *Box) CopyDimensionsOf(other *Box) *Box {
	b.width = other.Width()
	b.height = other.Height()

	return b
}

func (b *Box) Copy(other *Box) *Box {
	b.x = other.X()
	b.y = other.Y()
	b.width = other.Width()
	b.height = other.Height()

	return b
}

func (b *Box) SetPosition(x, y int) {
	b.x = x
	b.y = y
}

func (b *Box) Contains(other *Box) bool {
	return b.x <= other.X() && b.y <= other.Y() && (b.x+b.width) >= (other.X()+other.Width()) && (b.y+b.height) >= (other.Y()+other.Height())
}

func (b *Box) recalculateIfNeeded() {
	if b.needsRecalculation && len(b.calculationSteps) > 0 {
		b.recalculationCount++
		b.needsRecalculation = false
		b.ZeroOut()
		for _, step := range b.calculationSteps {
			step(b)
		}
	}
}

func (b *Box) X() int {
	b.recalculateIfNeeded()
	return b.x
}

func (b *Box) Y() int {
	b.recalculateIfNeeded()
	return b.y
}

func (b *Box) Width() int {
	b.recalculateIfNeeded()
	return b.width
}

func (b *Box) Height() int {
	b.recalculateIfNeeded()
	return b.height
}

func (b *Box) SetWidth(width int) *Box {
	b.width = width
	return b
}

func (b *Box) SetHeight(height int) *Box {
	b.height = height
	return b
}

func (b *Box) SetDimensions(width, height int) *Box {
	b.width, b.height = width, height
	return b
}

func (b *Box) SetX(x int) *Box {
	b.x = x
	return b
}

func (b *Box) SetY(y int) *Box {
	b.y = y
	return b
}

func (b *Box) SetOrigin(x, y int) *Box {
	b.x, b.y = x, y
	return b
}

func (b *Box) CenterVerticallyWithin(other *Box) *Box {
	otherCenterY := other.Y() + other.Height()/2
	y := otherCenterY - b.height/2

	b.SetOrigin(b.x, y)

	return b
}

func (b *Box) CenterHorizontallyWithin(other *Box) *Box {
	otherCenterX := other.X() + other.Width()/2
	x := otherCenterX - b.width/2

	b.SetOrigin(x, b.y)

	return b
}

func (b *Box) CenterWithin(other *Box) *Box {
	b.CenterHorizontallyWithin(other)
	b.CenterVerticallyWithin(other)

	return b
}

func (b *Box) AlignLeftWithin(other *Box) *Box {
	b.SetX(other.X())
	return b
}

func (b *Box) AlignRightWithin(other *Box) *Box {
	b.SetX(other.X() + other.Width() - b.Width())
	return b
}

func (b *Box) AlignTopWithin(other *Box) *Box {
	b.SetY(other.Y())
	return b
}

func (b *Box) AlignBottomWithin(other *Box) *Box {
	b.SetY(other.Y() + other.Height() - b.Height())
	return b
}

func Copy(b *Box) *Box {
	return &Box{
		x:      b.x,
		y:      b.y,
		width:  b.width,
		height: b.height,
	}
}

func Zeroed() *Box {
	return &Box{
		x:      0,
		y:      0,
		width:  0,
		height: 0,
	}
}

func Computed(calculateFn func(*Box)) *Box {
	return Zeroed().Computed(calculateFn)
}

func Bounding(boxes []*Box) *Box {
	boundingBox := Zeroed()

	// Find the minimum x and y coordinates
	for i, box := range boxes {
		if i == 0 {
			boundingBox.SetX(box.X())
			boundingBox.SetY(box.Y())
			boundingBox.SetWidth(box.Width())
			boundingBox.SetHeight(box.Height())
			continue
		}

		if box.X() < boundingBox.X() {
			boundingBox.SetX(box.X())
		}
		if box.Y() < boundingBox.Y() {
			boundingBox.SetY(box.Y())
		}
	}

	for _, box := range boxes {
		// boundingBoxX, boundingBoxY := boundingBox.X(), boundingBox.Y() // Which is zero zero
		boundingBoxWidth, boundingBoxHeight := boundingBox.Width(), boundingBox.Height()
		childX, childY, childWidth, childHeight := box.X(), box.Y(), box.Width(), box.Height()

		if childX+childWidth > boundingBoxWidth+boundingBox.X() {
			boundingBox.SetWidth(childX + childWidth - boundingBox.X())
		}

		if childY+childHeight > boundingBoxHeight+boundingBox.Y() {
			boundingBox.SetHeight(childY + childHeight - boundingBox.Y())
		}

		// if x < boundingBoxX {
		// 	boundingBox.SetX(func(previousX int) int {
		// 		return x
		// 	})
		// }
		// if y < boundingBoxY {
		// 	boundingBox.SetY(func(previousY int) int {
		// 		return y
		// 	})
		// }
		// if x+width > boundingBoxX+boundingBoxWidth {
		// 	boundingBox.SetWidth(func(previousWidth int) int {
		// 		return x + width - boundingBoxX
		// 	})
		// }
		// if y+height > boundingBoxY+boundingBoxHeight {
		// 	boundingBox.SetHeight(func(previousHeight int) int {
		// 		return y + height - boundingBoxY
		// 	})
		// }

	}

	return boundingBox
}

func (b *Box) MoveDown(amount int) *Box {
	b.SetY(b.Y() + amount)
	return b
}

func (b *Box) MoveUp(amount int) *Box {
	b.SetY(b.Y() - amount)
	return b
}

func (b *Box) MoveLeft(amount int) *Box {
	b.SetX(b.X() - amount)
	return b
}

func (b *Box) MoveRight(amount int) *Box {
	b.SetX(b.X() + amount)
	return b
}

func (b *Box) IncreaseWidth(amount int) *Box {
	b.SetWidth(b.Width() + amount)
	return b
}

func (b *Box) IncreaseHeight(amount int) *Box {
	b.SetHeight(b.Height() + amount)
	return b
}

func (b *Box) DecreaseWidth(amount int) *Box {
	b.SetWidth(b.Width() - amount)
	return b
}

func (b *Box) DecreaseHeight(amount int) *Box {
	b.SetHeight(b.Height() - amount)
	return b
}

func (b *Box) AlignRight(other *Box) *Box {
	b.SetX(other.X() + other.Width() - b.Width())
	return b
}

func (b *Box) AlignLeft(other *Box) *Box {
	b.SetX(other.X())
	return b
}

func (b *Box) AlignTop(other *Box) *Box {
	b.SetY(other.Y())
	return b
}

func (b *Box) AlignBottom(other *Box) *Box {
	b.SetY(other.Y() + other.Height() - b.Height())
	return b
}
