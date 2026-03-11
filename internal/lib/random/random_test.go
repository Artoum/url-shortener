package random

import (
	"testing"
)

func TestNewRandomString(t *testing.T) {
	aliasLength := 6
	result := NewRandomString(aliasLength)

	if len(result) != aliasLength {
		t.Errorf("Expected length %d, got %d", aliasLength, len(result))
	}

	t.Logf("Generated string: %s", result)
}
