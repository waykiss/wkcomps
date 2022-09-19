package list

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayStringReverse(t *testing.T) {
	values := []struct {
		input    []string
		expected []string
	}{
		{[]string{"a", "e", "i", "o", "u"}, []string{"u", "o", "i", "e", "a"}},
		{[]string{"1", "2", "3", "4", "5"}, []string{"5", "4", "3", "2", "1"}},
		{[]string{"a1", "b2", "c3", "d4", "e5", "f6"}, []string{"f6", "e5", "d4", "c3", "b2", "a1"}},
	}
	for _, v := range values {
		got := ArrayStringReverse(v.input)
		assert.Equal(t, v.expected, got, "Conversion of the value '%v' is expected '%v' but got '%v'",
			v.input, v.expected, got)
	}
}
