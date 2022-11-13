package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	t.Run("Must be able to create token", func(t *testing.T) {
		hash := New()
		hashed := hash.Create("12345")
		check := hash.Compare(hashed, "12345")

		assert.True(t, check)
	})
}
