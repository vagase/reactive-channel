package reactive_channel

import (
	"testing"
	"time"
)

func TestFromTo(t *testing.T) {
	in := []interface{}{1, 2, 3, 4}
	c := From(in)
	out := Values(c)

	if !isArrayEqual(in, out, true) {
		t.Error("From or Values failed")
	}
}

func TestInterval(t *testing.T) {
	out1 := Interval(timeoutContext(time.Millisecond*55), time.Millisecond*10, 0, nil)

	values := Values(out1)
	if len(values) != 5 {
		t.Fail()
	}

	index := 0
	out2 := Interval(timeoutContext(time.Millisecond*55), time.Millisecond*10, 0, func(i interface{}) interface{} {
		index++
		return index
	})

	assertChanWithValues(t, out2, []interface{}{1, 2, 3, 4, 5})
}

func TestRange(t *testing.T) {
	ch := Range(3, 4)
	array := []interface{}{3, 4, 5, 6}
	assertChanWithValues(t, ch, array)
}

func TestRepeat(t *testing.T) {
	repeat1 := Repeat(From([]interface{}{1, 2, 3, 4}), 2, true)
	assertChanWithValues(t, repeat1, []interface{}{1, 1, 2, 2, 3, 3, 4, 4})

	repeat2 := Repeat(From([]interface{}{1, 2, 3, 4}), 2, false)
	assertChanWithValues(t, repeat2, []interface{}{1, 2, 3, 4, 1, 2, 3, 4})
}

func TestJust(t *testing.T) {
	in := Just(1)
	assertChanWithValues(t, in, []interface{} {1})
}

func TestEmpty(t *testing.T) {
	in := Empty()
	assertChanWithValues(t, in, []interface{} {})
}