package non

import (
	"errors"
)

var (
	ErrEmpty = errors.New("empty string")
)

// Empty returns the first non-empty string. An empty string and ErrEmpty is
// returned if all strings are empty or nothing is given.
func Empty(strs ...string) (string, error) {
	for _, str := range strs {
		if len(str) > 0 {
			return str, nil
		}
	}
	return "", ErrEmpty
}
