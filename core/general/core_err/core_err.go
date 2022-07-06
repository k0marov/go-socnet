package core_err

import (
	"errors"
	"fmt"
	"github.com/k0marov/go-socnet/core/general/client_errors"
)

var ErrNotFound = errors.New("requested entity was not found")

func Rethrow(description string, err error) error {
	_, isClientError := err.(client_errors.ClientError)
	if isClientError {
		return err
	}
	return fmt.Errorf("%s: %w", description, err)
}
