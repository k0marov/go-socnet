package core_err

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("requested entity was not found")

func Rethrow(description string, err error) error {
	return fmt.Errorf("%s: %w", description, err)
}
