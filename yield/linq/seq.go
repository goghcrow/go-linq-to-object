package linq

// Sequence

type (
	Iter[T any] <-chan T
	Next[T any] func() (T, bool)
)

func (i Iter[T]) Next() (T, bool) {
	x, ok := <-i
	return x, ok
}

func (i Iter[T]) ToSlice() (ys []T) {
	for x := range i {
		ys = append(ys, x)
	}
	return
}
