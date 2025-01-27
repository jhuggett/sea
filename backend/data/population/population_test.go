package population

import (
	"testing"

	"github.com/jhuggett/sea/constructs/items"
	"github.com/jhuggett/sea/data"
	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {

	pop := Population{
		Persistent: data.Population{
			Size: 1000,

			Industries: []data.Industry{
				{
					Product: "grain",
					Workers: 1000,
				},
			},
		},
	}

	itemA := items.LookupItem(string(items.Grain))

	val, _ := pop.Value(itemA)

	assert.Equal(t, 0.0, val)
}
