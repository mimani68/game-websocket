package generic

// Optional type for generic optional values

type Optional[T any] struct {
	value *T
}

func NewOptional[T any](val T) Optional[T] {
	return Optional[T]{value: &val}
}

func EmptyOptional[T any]() Optional[T] {
	return Optional[T]{value: nil}
}

func (o Optional[T]) IsPresent() bool {
	return o.value != nil
}

func (o Optional[T]) Get() (T, bool) {
	if o.value == nil {
		var zero T
		return zero, false
	}
	return *o.value, true
}

func (o Optional[T]) OrElse(other T) T {
	if o.value == nil {
		return other
	}
	return *o.value
}
