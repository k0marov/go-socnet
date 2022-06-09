package option

type Option[T any] struct {
	value  T
	isNull bool
}

func (o *Option[T]) Fold(onValue func(value T), onNull func()) {
	if o.isNull {
		onNull()
	} else {
		onValue(o.value)
	}
}
