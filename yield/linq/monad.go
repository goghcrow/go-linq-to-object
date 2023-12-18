package linq

// https://groups.google.com/g/elm-discuss/c/rAfKkv2w1GU
//
//                   (a -> b) -> a -> b                Names: apply, <|, $
// Functor f      => (a -> b) -> f a -> f b            Names: map, fmap, <$>, <~
// Applicative f  => f (a -> b) -> f a -> f b          Names: ap, <*>, ~
// Monad f        => (a -> f b) -> f a -> f b          Names: bind (flipped), flatMap, concatMap, =<<
//
// flips
//
//                    a -> (a -> b) -> b              Names: |>, #
// Functor f     => f a -> (a -> b) -> f b            (never used???)
// Applicative f => f a -> f (a -> b) -> f b          Names: <**> (rarely used???)
// Monad f       => f a -> (a -> f b) -> f b          Names: bind, flatMap?, >>= (more common then its flip)

// ↓↓↓↓↓↓ Sequence Monad ↓↓↓↓↓↓

// Unit aka return
func Unit[T any](x T) Iter[T] {
	return Of(x)
}

// Bind aka flatMap
func Bind[A, R any](xs Iter[A], f func(A) Iter[R]) Iter[R] {
	iter := make(chan R, internalChanCap)
	go func() {
		for x := range xs {
			if ys := f(x); ys != nil { // DON'T BLOCK
				for y := range ys {
					iter <- y
				}
			}
		}
		close(iter)
	}()
	return iter
}

// ↓↓↓↓↓↓ Alias ↓↓↓↓↓↓

func Return[T any](x T) Iter[T] { return Unit(x) }

func FlatMap[A, R any](xs Iter[A], f func(A) Iter[R]) Iter[R] { return Bind(xs, f) }
