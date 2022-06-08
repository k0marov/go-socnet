package ref

import "errors"

// Ref is a pointer which can't be nil
type Ref[T any] struct {
	pointer *T
}

func (r Ref[T]) Value() T {
	return *r.pointer
}

func NewRef[T any](pointer *T) (Ref[T], error) {
	if pointer == nil {
		return Ref[T]{}, errors.New("references can't be created from a nil pointer")
	}
	return Ref[T]{pointer: pointer}, nil
}
