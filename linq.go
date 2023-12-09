package linq

type Index = int

func SelectMany[A, R any](xs Seq[A], f func(A) Seq[R]) Seq[R] {
	return Bind[A, R](xs, f)
}

func SelectManyWithIndex[A, R any](xs Seq[A], f func(A, Index) Seq[R]) Seq[R] {
	idx := 0
	return SelectMany(xs, func(a A) (r Seq[R]) {
		r = f(a, idx)
		idx++
		return
	})
}

// Select aka map
func Select[A, R any](xs Seq[A], f func(A) R) Seq[R] {
	return SelectWithIndex(xs, func(x A, _ Index) R {
		return f(x)
	})
}

func SelectWithIndex[A, R any](xs Seq[A], f func(A, Index) R) Seq[R] {
	return SelectManyWithIndex(xs, func(x A, i Index) Seq[R] {
		return Return(f(x, i))
	})
}

func First[A any](xs Seq[A]) (A, bool) {
	return Take(xs, 1).Next()
}

func FirstWhile[A any](xs Seq[A], f func(A) bool) (A, bool) {
	return TakeWhile(xs, f).Next()
}

func Last[A any](xs Seq[A]) (last A, ok bool) {
	return LastWhile(xs, Const[A](true))
}

func LastWhile[A any](xs Seq[A], f func(A) bool) (last A, ok bool) {
	Iterate(xs, func(x A) {
		if f(x) {
			ok = true
			last = x
		}
	})
	return
}

func Where[A any](xs Seq[A], f func(A) bool) Seq[A] {
	return WhereWithIndex(xs, func(x A, _ Index) bool {
		return f(x)
	})
}

func WhereWithIndex[A any](xs Seq[A], f func(A, Index) bool) Seq[A] {
	return SelectManyWithIndex(xs, func(x A, i Index) Seq[A] {
		if f(x, i) {
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

func TakeWhile[A any](xs Seq[A], f func(A) bool) Seq[A] {
	return TakeWhileWithIndex(xs, func(x A, _ Index) bool {
		return f(x)
	})
}

func TakeWhileWithIndex[A any](xs Seq[A], f func(A, Index) bool) Seq[A] {
	return WhereWithIndex(xs, f)
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

func SkipWhile[A any](xs Seq[A], f func(A) bool) Seq[A] {
	return SkipWhileWithIndex(xs, func(x A, _ Index) bool {
		return f(x)
	})
}

func SkipWhileWithIndex[A any](xs Seq[A], f func(A, Index) bool) Seq[A] {
	return SelectManyWithIndex(xs, func(x A, i Index) Seq[A] {
		if f(x, i) {
			return nil
		}
		return Return(x)
	})
}

func Aggregate[A, B, R any](
	xs Seq[A],
	init B,
	f func(acc B, cur A) B,
	selector func(B) R,
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
