package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsEmptyString(t *testing.T) {
	test := func(str string, empty bool) {
		assert.Equal(t, IsEmptyString(str), empty)
	}
	test("", true)
	test(" ", true)
	test(" a ", false)
	test("	", true)
}
