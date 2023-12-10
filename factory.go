package linq

func SeqOf[T any](f FSeq[T]) Seq[T] {
	return f
}

func From[T any](xs ...T) Seq[T] {
	// return &SliceSeq[T]{xs: xs}
	var i int
	return SeqOf[T](func() (x T, ok bool) {
		ok = i < len(xs)
		if ok {
			x, i = xs[i], i+1
		}
		return
	})
}

func FromSlice[T any](xs []T) Seq[T] {
	// return From(xs...)
	return &SliceSeq[T]{xs: xs}
}

func FromMap[K comparable, V any](xs map[K]V) Seq[Cons[K, V]] {
	var ks []K
	for k := range xs {
		ks = append(ks, k)
	}
	return Select(FromSlice(ks), func(x K) Cons[K, V] {
		return Cons[K, V]{x, xs[x]}
	})
}

func Range(start, end int) Seq[int] {
	i := start
	return SeqOf[int](func() (x int, ok bool) {
		ok = i < end
		if ok {
			x, i = i, i+1
		}
		return
	})
}
