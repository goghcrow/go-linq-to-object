package example

import (
	"testing"

	"github.com/goghcrow/go-linq-object/yield/iter"
)

func TestIter(t *testing.T) {
	it := iter.Of(1, 2, 3, 4, 5)
	for v := range it {
		println(v)
	}

	for v := range it {
		println(v)
	}
}

func TestIterNext(t *testing.T) {
	it := iter.Of(1, 2, 3, 4, 5)
	for {
		nxt, has := it.Next()
		if !has {
			break
		}
		println(nxt)
	}

	for {
		nxt, has := it.Next()
		if !has {
			break
		}
		println(nxt)
	}
}

func TestInfXIter(t *testing.T) {
	it := iter.Infinite(42)
	for x := range it {
		println(x)
		// break
	}
}

func TestInfIter(t *testing.T) {
	it := iter.From(func() (int, bool) {
		return 42, true
	})
	for x := range it {
		println(x)
		// break
	}
}

func TestRange(t *testing.T) {
	it := iter.Range(10, 20)
	for i := range it {
		println(i)
	}
}
