package non

import (
	"testing"
)

func TestEmpty_Nothing(t *testing.T) {
	str, err := Empty()

	if str != "" {
		t.Errorf("expected empty string, got %s", str)
	}
	if err != ErrEmpty {
		t.Errorf("expected ErrEmpty, got %s", err)
	}
}

func TestEmpty_EmptyStrings(t *testing.T) {
	str, err := Empty("", "", "")

	if str != "" {
		t.Errorf("expected an empty string, got %s", str)
	}
	if err != ErrEmpty {
		t.Errorf("expected ErrEmpty, got %s", err)
	}
}

func TestEmpty_Mixed(t *testing.T) {
	str, err := Empty("", "abc", "def")

	if str != "abc" {
		t.Errorf("expected 'abc', got %s", str)
	}
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}
