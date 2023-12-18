package linq

type (
	Index                           = int
	Pred[T any]                     func(T) bool
	IdxPred[T any]                  func(T, Index) bool
	Selector[Source, Result any]    func(Source) Result
	IdxSelector[Source, Result any] func(Source, Index) Result
)

func SelectMany[A, R any](xs Seq[A], f func(A) Seq[R]) Seq[R] {
	return Bind[A, R](xs, f)
}

func SelectManyWithIndex[A, R any](xs Seq[A], f func(A, Index) Seq[R]) Seq[R] {
	idx := 0
	return SelectMany(xs, func(a A) (r Seq[R]) {
		r, idx = f(a, idx), idx+1
		return
	})
}

// Select aka map
func Select[A, R any](xs Seq[A], f Selector[A, R]) Seq[R] {
	return SelectWithIndex(xs, func(x A, _ Index) R {
		return f(x)
	})
}

func SelectWithIndex[A, R any](xs Seq[A], f IdxSelector[A, R]) Seq[R] {
	// // e.g.
	// i := 0
	// return SeqOf(func() (r R, ok bool) {
	// 	x, ok := xs.Next()
	// 	if !ok {
	// 		return
	// 	}
	// 	r, ok, i = f(x, i), true, i+1
	// 	return
	// })

	return SelectManyWithIndex(xs, func(x A, i Index) Seq[R] {
		return Return(f(x, i))
	})
}

func First[A any](xs Seq[A]) (A, bool) {
	return Take(xs, 1).Next()
}

func FirstWhile[A any](xs Seq[A], p Pred[A]) (A, bool) {
	return TakeWhile(xs, p).Next()
}

func Last[A any](xs Seq[A]) (last A, ok bool) {
	return LastWhile(xs, Const[A](true))
}

func LastWhile[A any](xs Seq[A], p Pred[A]) (last A, ok bool) {
	Iterate(xs, func(x A) {
		if p(x) {
			ok = true
			last = x
		}
	})
	return
}

func Where[A any](xs Seq[A], p Pred[A]) Seq[A] {
	return WhereWithIndex(xs, func(x A, _ Index) bool {
		return p(x)
	})
}

func WhereWithIndex[A any](xs Seq[A], p IdxPred[A]) Seq[A] {
	return SelectManyWithIndex(xs, func(x A, i Index) Seq[A] {
		if p(x, i) {
			return Return(x)
		}
		return nil
	})
}

func Take[A any](xs Seq[A], cnt int) Seq[A] {
	return SeqOf[A](func() (x A, ok bool) {
		if cnt <= 0 {
			return
		}
		cnt--
		return xs.Next()
	})
}

func TakeWhile[A any](xs Seq[A], p Pred[A]) Seq[A] {
	return TakeWhileWithIndex(xs, func(x A, _ Index) bool {
		return p(x)
	})
}

func TakeWhileWithIndex[A any](xs Seq[A], p IdxPred[A]) Seq[A] {
	return WhereWithIndex(xs, p)
}

func Skip[A any](xs Seq[A], cnt int) Seq[A] {
	return SelectMany(xs, func(x A) Seq[A] {
		if cnt <= 0 {
			return Return(x)
		}
		cnt--
		return nil
	})
}

func SkipWhile[A any](xs Seq[A], p Pred[A]) Seq[A] {
	return SkipWhileWithIndex(xs, func(x A, _ Index) bool {
		return p(x)
	})
}

func SkipWhileWithIndex[A any](xs Seq[A], p IdxPred[A]) Seq[A] {
	return SelectManyWithIndex(xs, func(x A, i Index) Seq[A] {
		if p(x, i) {
			return nil
		}
		return Return(x)
	})
}

func Aggregate[A, B, R any](
	xs Seq[A],
	init B,
	f func(acc B, cur A) B,
	selector Selector[B, R],
) (r R) {
	acc := init
	Iterate(xs, func(x A) {
		acc = f(acc, x)
	})
	return selector(acc)
}

// Fold Aggregate with init
func Fold[A, R any](xs Seq[A], init R, f func(acc R, cur A) R) (acc R) {
	acc = init
	Iterate(xs, func(x A) {
		acc = f(acc, x)
	})
	return
}

// Reduce Aggregate
func Reduce[A any](xs Seq[A], f func(acc A, cur A) A) (r A, ok bool) {
	r, ok = xs.Next()
	if !ok {
		return
	}
	Iterate(xs, func(x A) {
		r = f(r, x)
	})
	ok = true
	return
}

func All[A any](xs Seq[A], p Pred[A]) (r bool) {
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

func AnyElem[A any](xs Seq[A]) bool {
	_, ok := xs.Next()
	return ok
}

func Any[A any](xs Seq[A], p Pred[A]) (r bool) {
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

func Append[A any](xs Seq[A], a A) Seq[A] {
	end := false
	return SeqOf(func() (x A, ok bool) {
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

func Iterate[T any](xs Seq[T], f func(T)) {
	IterateWithIndex(xs, func(x T, _ Index) {
		f(x)
	})
}

func IterateWithIndex[T any](xs Seq[T], f func(T, Index)) {
	i := 0
	for {
		x, ok := xs.Next()
		if !ok {
			break
		}
		f(x, i)
		i++
	}
}

func ToSlice[T any](xs Seq[T]) (ys []T) {
	for {
		x, ok := xs.Next()
		if !ok {
			break
		}
		ys = append(ys, x)
	}
	return
}
