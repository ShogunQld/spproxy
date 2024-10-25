package assert

import (
	"strings"
	"testing"
)

func Equal[V comparable](t *testing.T, got, expected V) {
	t.Helper()

	if expected != got {
		t.Errorf(`assert.Equal(
     got: %v
expected: %v
)`, got, expected)
	}
}

func NotEmpty(t *testing.T, str *string) {
	t.Helper()

	if str == nil || strings.TrimSpace(*str) == "" {
		t.Errorf("assert.NotEmpty(%v)", str)
	}
}

func NotNull(t *testing.T, got interface{}) {
	t.Helper()

	if got != nil {
		t.Errorf(`assert.NotNull(nil)`)
	}
}
