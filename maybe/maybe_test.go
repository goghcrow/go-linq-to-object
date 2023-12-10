package maybe

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, x, y any) {
	if !reflect.DeepEqual(x, y) {
		t.Fail()
	}
}

func TestMaybe(t *testing.T) {
	{
		x := Just(1)
		y := Just(2)
		z := Bind(x, func(a int) Maybe[int] {
			return Bind(y, func(b int) Maybe[int] {
				return Just[int](a + b)
			})
		})
		assertEqual(t, z, Just(3))
	}

	{
		x := Nothing[int]()
		y := Just(2)
		z := Bind(x, func(a int) Maybe[int] {
			return Bind(y, func(b int) Maybe[int] {
				return Just[int](a + b)
			})
		})
		assertEqual(t, z, Nothing[int]())
	}

	{
		x := Just(1)
		y := Nothing[int]()
		z := Bind(x, func(a int) Maybe[int] {
			return Bind(y, func(b int) Maybe[int] {
				return Just[int](a + b)
			})
		})
		assertEqual(t, z, Nothing[int]())
	}

	{
		x := Just(1)
		y := Map(x, func(a int) int {
			return a + 2
		})
		assertEqual(t, y, Just(3))
	}

	{
		x := Nothing[int]()
		y := Map(x, func(a int) int {
			return a + 2
		})
		assertEqual(t, y, Nothing[int]())
	}
}
