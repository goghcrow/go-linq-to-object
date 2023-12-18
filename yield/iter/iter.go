package iter

// golang 没有迭代器,
// 其他语言 迭代器通常可以用 for 之类语法糖来遍历, golang 中
// 除了 slice/map, 还有 chan 可以用 for 来遍历
// so, 可以用 chan 来替代迭代器使用
// 且 chan 只能读取一次, 符合迭代器的语义
// 可以先不考虑 chan 的性能锁之类问题

type (
	Iter[T any] <-chan T
	Next[T any] func() (T, bool)
)

func (iter Iter[T]) Next() (T, bool) {
	x, ok := <-iter
	return x, ok
}

func Of[T any](xs ...T) Iter[T] {
	iter := make(chan T, len(xs))
	for _, x := range xs {
		iter <- x
	}
	close(iter)
	return iter
}

func From[T any](f Next[T]) Iter[T] {
	iter := make(chan T, 1)
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
