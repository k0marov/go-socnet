package ref

import "errors"

// Ref is a pointer which can't be nil
type Ref[T any] interface {
	Value() T
}

type ref[T any] struct {
	pointer *T
}

func (r ref[T]) Value() T {
	return *r.pointer
}

func NewRef[T any](pointer *T) (ref[T], error) {
	if pointer == nil {
		return ref[T]{}, errors.New("references can't be created from a nil pointer")
	}
	return ref[T]{pointer: pointer}, nil
}
