package linq

import (
	"reflect"
	"strings"
	"testing"
)

type (
	Signed interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}
	Unsigned interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	}
	Integer interface{ Signed | Unsigned }
	Float   interface{ ~float32 | ~float64 }
	Number  interface{ Integer | Float }
)

func eq[A comparable](n A) func(A) bool { return func(x A) bool { return x == n } }
func lt[N Number](n N) func(N) bool     { return func(x N) bool { return x < n } }
func gt[N Number](n N) func(N) bool     { return func(x N) bool { return x > n } }
func le[N Number](n N) func(N) bool     { return func(x N) bool { return x <= n } }
func ge[N Number](n N) func(N) bool     { return func(x N) bool { return x >= n } }

func isEven(x int) bool { return x%2 == 0 }
func square(x int) int  { return x * x }
func double(x int) int  { return x * 2 }

func assertEqual(t *testing.T, x, y any) {
	if !reflect.DeepEqual(x, y) {
		t.Fail()
	}
}

func TestDeferred(t *testing.T) {
	xs := SeqOf[int](func() (int, bool) {
		panic("deferred")
	})
	_ = Select[int](xs, double)
}

func TestOnce(t *testing.T) {
	{
		xs := Of(1, 2, 3)
		ys := Select(xs, Id[int])
		assertEqual(t, ToSlice(ys), []int{1, 2, 3})

		zs := Select(xs, Id[int])
		assertEqual(t, len(ToSlice(zs)), 0)
	}

	{
		type T = Cons[string, int]
		xs := OfSlice([]string{"a", "b", "c"})
		ys := OfSlice([]int{1, 2, 3})
		zs := SelectMany(xs, func(x string) Seq[T] {
			return SelectMany(ys, func(y int) Seq[T] {
				return Return(T{x, y})
			})
		})
		assertEqual(t, ToSlice(zs), []T{
			{"a", 1}, T{"a", 2}, T{"a", 3},
		})
	}
}

func TestSelect(t *testing.T) {
	xs := Of(1, 2, 3)
	ys := Select(xs, square)
	assertEqual(t, ToSlice(ys), []int{1, 4, 9})
}

func TestSelectWithIndex(t *testing.T) {
	xs := Of(1, 2, 3)
	ys := SelectWithIndex(xs, func(a int, i Index) int {
		return a + i
	})
	slice := ToSlice(ys)
	assertEqual(t, slice, []int{1, 3, 5})
}

func TestWhere(t *testing.T) {
	xs := Range(1, 10)
	ys := Where(xs, isEven)
	assertEqual(t, ToSlice(ys), []int{2, 4, 6, 8})
}

func TestSkip(t *testing.T) {
	xs := Range(1, 10)
	ys := Skip(
		Where(xs, isEven),
		2,
	)
	assertEqual(t, ToSlice(ys), []int{6, 8})
}

func TestSkipWhile(t *testing.T) {
	xs := Range(1, 10)
	ys := SkipWhile(xs, lt(5))
	assertEqual(t, ToSlice(ys), []int{5, 6, 7, 8, 9})
}

func TestSkipWhileWithIndex(t *testing.T) {
	xs := Range(1, 10)
	ys := SkipWhileWithIndex(xs, func(x int, idx Index) bool {
		return x+idx < 10
	})
	assertEqual(t, ToSlice(ys), []int{6, 7, 8, 9})
}

func TestTake(t *testing.T) {
	xs := Range(1, 10)
	ys := Take(
		Where(xs, isEven),
		3,
	)
	assertEqual(t, ToSlice(ys), []int{2, 4, 6})
}

func TestTakeWhile(t *testing.T) {
	xs := Range(1, 10)
	ys := TakeWhile(xs, lt(5))
	assertEqual(t, ToSlice(ys), []int{1, 2, 3, 4})
}

func TestTakeWhileWithIndex(t *testing.T) {
	xs := Range(1, 10)
	ys := TakeWhileWithIndex(xs, func(x int, idx Index) bool {
		return x+idx < 10
	})
	assertEqual(t, ToSlice(ys), []int{1, 2, 3, 4, 5})
}

func TestFirst(t *testing.T) {
	xs := OfSlice([]int{9, 34, 65, 92, 87, 435, 3, 54,
		83, 23, 87, 435, 67, 12, 19})
	first, ok := First(xs)
	assertEqual(t, first, 9)
	assertEqual(t, ok, true)
}

func TestFirstWhile(t *testing.T) {
	xs := OfSlice([]int{9, 34, 65, 92, 87, 435, 3, 54,
		83, 23, 87, 435, 67, 12, 19})
	first, ok := FirstWhile(xs, gt(80))
	assertEqual(t, first, 92)
	assertEqual(t, ok, true)
}

func TestLast(t *testing.T) {
	xs := OfSlice([]int{9, 34, 65, 92, 87, 435, 3, 54,
		83, 23, 87, 67, 12, 19})
	first, ok := Last(xs)
	assertEqual(t, first, 19)
	assertEqual(t, ok, true)
}

func TestLastWhile(t *testing.T) {
	xs := OfSlice([]int{9, 34, 65, 92, 87, 435, 3, 54,
		83, 23, 87, 67, 12, 19})
	first, ok := LastWhile(xs, gt(80))
	assertEqual(t, first, 87)
	assertEqual(t, ok, true)
}

func TestAggregate(t *testing.T) {
	fruits := Of("apple", "mango", "orange", "passionfruit", "grape")
	s := Aggregate[string, string](fruits, "banana", func(longest string, next string) string {
		if len(next) > len(longest) {
			return next
		}
		return longest
	}, strings.ToUpper)

	assertEqual(t, s, "PASSIONFRUIT")
}

func TestFold(t *testing.T) {
	{
		xs := OfSlice([]int{1, 4, 5})
		r := Fold(xs, 5, func(acc int, cur int) int {
			return acc*2 + cur
		})
		assertEqual(t, r, 57)
	}

	{
		xs := OfSlice([]int{4, 8, 8, 3, 9, 0, 7, 8, 2})
		r := Fold(xs, 0, func(total int, cur int) int {
			if isEven(cur) {
				return total + 1
			}
			return total
		})
		assertEqual(t, r, 6)
	}
}

func TestReduce(t *testing.T) {
	sentence := "the quick brown fox jumps over the lazy dog"
	words := strings.Split(sentence, " ")
	reversed, ok := Reduce(OfSlice(words), func(acc string, cur string) string {
		return cur + " " + acc
	})
	assertEqual(t, ok, true)
	assertEqual(t, reversed, "dog lazy the over jumps fox brown quick the")
}

// left outer join
func TestCrossJoin(t *testing.T) {
	// from inner in items
	// from outer in function(items)
	// select projection(inner, outer)

	type T = Cons[string, int]

	xs := []string{"a", "b", "c"}
	ys := []int{1, 2, 3}

	zs := SelectMany(OfSlice(xs), func(x string) Seq[T] {
		return SelectMany(OfSlice(ys), func(y int) Seq[T] {
			return Return(T{x, y})
		})
	})

	assertEqual(t, ToSlice(zs), []T{
		{"a", 1}, T{"a", 2}, T{"a", 3},
		{"b", 1}, T{"b", 2}, T{"b", 3},
		{"c", 1}, T{"c", 2}, T{"c", 3},
	})
}
