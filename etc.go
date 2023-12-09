package linq

type Cons[T1, T2 any] struct {
	Car T1
	Cdr T2
}

func Const[Any, R any](x R) func(Any) R {
	return func(Any) R {
		return x
	}
}

func Id[A any](a A) A {
	return a
}
