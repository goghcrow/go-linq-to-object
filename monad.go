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
func Unit[T any](x T) Seq[T] {
	hasNxt := true
	return SeqOf[T](func() (z T, ok bool) {
		if hasNxt {
			hasNxt = false
			return x, true
		}
		return
	})
}

// Bind aka flatMap
func Bind[A, R any](xs Seq[A], f func(A) Seq[R]) Seq[R] {
	var y Seq[R]
	return SeqOf[R](func() (z R, ok bool) {
	begin:
		if y != nil {
			z, ok = y.Next()
			if ok {
				return
			}
			y = nil
		}
		for {
			var x A
			x, ok = xs.Next()
			if !ok {
				return
			}
			y = f(x)
			goto begin
		}
	})
}

// ↓↓↓↓↓↓ Alias ↓↓↓↓↓↓

func Return[T any](x T) Seq[T] { return Unit(x) }

func FlatMap[A, R any](xs Seq[A], f func(A) Seq[R]) Seq[R] { return Bind(xs, f) }
