package maybe

type Maybe[T any] struct {
	Just  bool
	Value T
}

func Just[T any](x T) Maybe[T] {
	return Maybe[T]{
		Just:  true,
		Value: x,
	}
}

func Nothing[T any]() Maybe[T] {
	return Maybe[T]{
		Just: false,
	}
}

func Unit[T any](x T) Maybe[T] {
	return Just[T](x)
}

func Bind[A, R any](m Maybe[A], f func(A) Maybe[R]) Maybe[R] {
	if m.Just {
		return f(m.Value)
	}
	return Nothing[R]()
}

// ----------------------------------------

func Return[A any](x A) Maybe[A] {
	return Unit[A](x)
}

func FlatMap[A, R any](m Maybe[A], f func(A) Maybe[R]) Maybe[R] {
	return Bind(m, f)
}

func Map[A, R any](m Maybe[A], f func(A) R) Maybe[R] {
	return Bind(m, func(a A) Maybe[R] {
		return Unit[R](f(a))
	})
}
