package linq

// Sequence Interface

type Seq[T any] interface {
	Next() (T, bool)
}

// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓
// Fun Sequence Implementations

type FSeq[T any] func() (T, bool)

func (f FSeq[T]) Next() (T, bool) { return f() }

// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓
// Slice Sequence Implementations

type SliceSeq[T any] struct {
	xs []T
	i  int
}

func (s *SliceSeq[T]) Next() (x T, ok bool) {
	if s.i >= len(s.xs) {
		return
	}
	x, s.i = s.xs[s.i], s.i+1
	return x, true
}
