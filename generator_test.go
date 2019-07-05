package reactive_channel

import (
	"context"
	"testing"
	"time"
)

func TestFromTo(t *testing.T) {
	in := []interface{}{1, 2, 3, 4}
	c := From(in)
	out := To(c)

	if !isArrayEqual(in, out, true) {
		t.Error("From or To failed")
	}
}

func TestMap(t *testing.T) {
	in := From([]interface{}{1, 2, 3, 4})

	out := Map(in, func(i interface{}) interface{} {
		num := i.(int)
		return num * 2
	})

	assertChanWithValues(t, out, []interface{}{2, 4, 6, 8})
}

func TestInterval(t *testing.T) {
	ctx , _ := context.WithTimeout(context.Background(), time.Millisecond * 55)

	ch := Interval(ctx, time.Millisecond * 10)

	values := To(ch)
	if len(values) != 5 {
		t.Fail()
	}
}

func TestRange(t *testing.T) {
	ch := Range(3, 4)
	array := [] interface{} {3, 4, 5, 6}
	assertChanWithValues(t, ch, array)
}