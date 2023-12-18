package linq

type (
	Index                           = int
	Pred[T any]                     func(T) bool
	IdxPred[T any]                  func(T, Index) bool
	Selector[Source, Result any]    func(Source) Result
	IdxSelector[Source, Result any] func(Source, Index) Result
)

func SelectMany[A, R any](xs Iter[A], f func(A) Iter[R]) Iter[R] {
	return Bind[A, R](xs, f)
}

func SelectManyWithIndex[A, R any](xs Iter[A], f func(A, Index) Iter[R]) Iter[R] {
	idx := 0
	return SelectMany(xs, func(a A) (r Iter[R]) {
		r, idx = f(a, idx), idx+1
		return
	})
}

// Select aka map
func Select[A, R any](xs Iter[A], f Selector[A, R]) Iter[R] {
	return SelectWithIndex(xs, func(x A, _ Index) R {
		return f(x)
	})
}

func SelectWithIndex[A, R any](xs Iter[A], f IdxSelector[A, R]) Iter[R] {
	// // e.g.
	// i := 0
	// return From(func() (r R, ok bool) {
	// 	x, ok := xs.Next()
	// 	if !ok {
	// 		return
	// 	}
	// 	r, ok, i = f(x, i), true, i+1
	// 	return
	// })

	return SelectManyWithIndex(xs, func(x A, i Index) Iter[R] {
		return Return(f(x, i))
	})
}

func First[A any](xs Iter[A]) (A, bool) {
	return Take(xs, 1).Next()
}

func FirstWhile[A any](xs Iter[A], p Pred[A]) (A, bool) {
	return TakeWhile(xs, p).Next()
}

func Last[A any](xs Iter[A]) (last A, ok bool) {
	return LastWhile(xs, Const[A](true))
}

func LastWhile[A any](xs Iter[A], p Pred[A]) (last A, ok bool) {
	for x := range xs {
		if p(x) {
			ok = true
			last = x
		}
	}
	return
}

func Where[A any](xs Iter[A], p Pred[A]) Iter[A] {
	return WhereWithIndex(xs, func(x A, _ Index) bool {
		return p(x)
	})
}

func WhereWithIndex[A any](xs Iter[A], p IdxPred[A]) Iter[A] {
	return SelectManyWithIndex(xs, func(x A, i Index) Iter[A] {
		if p(x, i) {
			return Return(x)
		}
		return nil
	})
}

func Take[A any](xs Iter[A], cnt int) Iter[A] {
	return From[A](func() (x A, ok bool) {
		if cnt <= 0 {
			return
		}
		cnt--
		return xs.Next()
	})
}

func TakeWhile[A any](xs Iter[A], p Pred[A]) Iter[A] {
	return TakeWhileWithIndex(xs, func(x A, _ Index) bool {
		return p(x)
	})
}

func TakeWhileWithIndex[A any](xs Iter[A], p IdxPred[A]) Iter[A] {
	return WhereWithIndex(xs, p)
}

func Skip[A any](xs Iter[A], cnt int) Iter[A] {
	return SelectMany(xs, func(x A) Iter[A] {
		if cnt <= 0 {
			return Return(x)
		}
		cnt--
		return nil
	})
}

func SkipWhile[A any](xs Iter[A], p Pred[A]) Iter[A] {
	return SkipWhileWithIndex(xs, func(x A, _ Index) bool {
		return p(x)
	})
}

func SkipWhileWithIndex[A any](xs Iter[A], p IdxPred[A]) Iter[A] {
	return SelectManyWithIndex(xs, func(x A, i Index) Iter[A] {
		if p(x, i) {
			return nil
		}
		return Return(x)
	})
}

func Aggregate[A, B, R any](
	xs Iter[A],
	init B,
	f func(acc B, cur A) B,
	selector Selector[B, R],
) (r R) {
	acc := init
	for x := range xs {
		acc = f(acc, x)
	}
	return selector(acc)
}

// Fold Aggregate with init
func Fold[A, R any](xs Iter[A], init R, f func(acc R, cur A) R) (acc R) {
	acc = init
	for x := range xs {
		acc = f(acc, x)
	}
	return
}

// Reduce Aggregate
func Reduce[A any](xs Iter[A], f func(acc A, cur A) A) (r A, ok bool) {
	r, ok = xs.Next()
	if !ok {
		return
	}
	for x := range xs {
		r = f(r, x)
	}
	ok = true
	return
}

func All[A any](xs Iter[A], p Pred[A]) (r bool) {
	r = true
	for {
		x, ok := xs.Next()
		if !ok {
			break
		}
		r = r && p(x)
		if !r {
			break
		}
	}
	return
}

func AnyElem[A any](xs Iter[A]) bool {
	_, ok := xs.Next()
	return ok
}

func Any[A any](xs Iter[A], p Pred[A]) (r bool) {
	for {
		x, ok := xs.Next()
		if !ok {
			break
		}
		r = r || p(x)
		if r {
			break
		}
	}
	return
}

func Append[A any](xs Iter[A], a A) Iter[A] {
	end := false
	return From(func() (x A, ok bool) {
		if end {
			return
		}
		x, ok = xs.Next()
		if !ok {
			end = true
			return a, true
		}
		return
	})
}

func Count[A any](xs Iter[A]) (cnt int) {
	for _ = range xs {
		cnt++
	}
	return
}
