package name

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {

	firstName := Generate(2)
	lastName := Generate(4)

	nickname := GenerateNickName()

	n := firstName + " " + lastName + " the " + nickname

	assert.Equal(t, n, "")
}
