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

func TestRepeat(t *testing.T) {
	repeat1 := Repeat(From([]interface{}{1, 2, 3, 4}), 2, true)
	assertChanWithValues(t, repeat1, []interface{}{1,1,2,2,3,3,4,4})

	repeat2 := Repeat(From([]interface{}{1, 2, 3, 4}), 2, false)
	assertChanWithValues(t, repeat2, []interface{}{1,2,3,4,1,2,3,4})
}