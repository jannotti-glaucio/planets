package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	t.Run("Should be able to reverse string", func(t *testing.T) {
		input := "123456"
		expected := "654321"
		actual := Reverse(input)

		assert.EqualValues(t, expected, actual)
	})

	t.Run("Should be able to replace string in pattern", func(t *testing.T) {
		str := "123"
		pattern := "[0-9]"
		expected := "aaa"
		actual := ReplacePattern(str, pattern, "a")

		assert.EqualValues(t, expected, actual)
	})

	t.Run("Should be able to find string in string", func(t *testing.T) {
		str := "123456"
		substring := "456"
		ok := Contains(str, substring)

		assert.True(t, ok)
	})

	t.Run("Should be able to find string in pattern", func(t *testing.T) {
		str := "123456"
		substring := "[0-9]"
		ok := Matches(str, substring)

		assert.True(t, ok)
	})
}
