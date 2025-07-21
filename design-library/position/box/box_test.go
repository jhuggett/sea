package box

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {

	// mainBox := NewBox(BoxConfig{
	// 	width:  200,
	// 	height: 100,
	// })

	// exampleBox := NewBox(BoxConfig{})
	// exampleBox.Copy(mainBox)

	boxA := New(Config{
		X:      10,
		Y:      0,
		Width:  100,
		Height: 50,
	})

	boxB := New(Config{})

	boxB.Computed(func(b *Box) {
		b.MoveBelow(boxA)
	})

	assert.Equal(t, 0, boxB.X(), "Box B X should be 0")
	assert.Equal(t, 50, boxB.Y(), "Box B Y should be 50")
	assert.Equal(t, 50, boxB.Y(), "Box B Y should still be 50")

	boxA.SetHeight(100)
	assert.Equal(t, 100, boxA.Height(), "Box A Height should be 100")

	assert.Equal(t, 0, boxB.X(), "Box B X should still be 0")
	assert.Equal(t, 100, boxB.Y(), "Box B Y should now be 100")

}

func TestLayeredComputed(t *testing.T) {
	boxA := New(Config{
		X:      0,
		Y:      0,
		Width:  100,
		Height: 50,
	})

	someValue := 5

	boxA.Computed(func(b *Box) {
		b.SetOrigin(someValue, someValue)
	})

	boxA.FlagNeedsRecalculation()

	assert.Equal(t, 5, boxA.X(), "Box A X should be 10")
	assert.Equal(t, 5, boxA.Y(), "Box A Y should be 10")

}
