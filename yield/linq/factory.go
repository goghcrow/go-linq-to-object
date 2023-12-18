package linq

func Of[T any](xs ...T) Iter[T] {
	iter := make(chan T, len(xs))
	for _, x := range xs {
		iter <- x
	}
	close(iter)
	return iter
}

func From[T any](f Next[T]) Iter[T] {
	iter := make(chan T, internalChanCap)
	go func() {
		for {
			x, has := f()
			if !has {
				break
			}
			iter <- x
		}
		close(iter)
	}()
	return iter
}

func OfMap[K comparable, V any](xs map[K]V) Iter[Cons[K, V]] {
	var ks []K
	for k := range xs {
		ks = append(ks, k)
	}
	return Select(Of(ks...), func(x K) Cons[K, V] {
		return Cons[K, V]{x, xs[x]}
	})
}

func Range(minInclusive, maxExclusive int) Iter[int] {
	return From(func() (int, bool) {
		if minInclusive < maxExclusive {
			minInclusive++
			return minInclusive - 1, true
		}
		return 0, false
	})
}

func Infinite[T any](x T) Iter[T] {
	return From(func() (T, bool) {
		return x, true
	})
}
