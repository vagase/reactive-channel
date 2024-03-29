package reactive_channel

import (
	"context"
	"github.com/bmizerany/assert"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

type SortInterfaceArray []interface{}

func (s SortInterfaceArray) Len() int {
	return len(s)
}

func (s SortInterfaceArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortInterfaceArray) Less(i, j int) bool {
	switch s[i].(type) {
	case int:
		return s[i].(int) < s[j].(int)
	case uint:
		return s[i].(uint) < s[j].(uint)
	case int8:
		return s[i].(int8) < s[j].(int8)
	case int16:
		return s[i].(int16) < s[j].(int16)
	case int32:
		return s[i].(int32) < s[j].(int32)
	case int64:
		return s[i].(int64) < s[j].(int64)
	case uint8:
		return s[i].(uint8) < s[j].(uint8)
	case uint16:
		return s[i].(uint16) < s[j].(uint16)
	case uint32:
		return s[i].(uint32) < s[j].(uint32)
	case uint64:
		return s[i].(uint64) < s[j].(uint64)
	case float32:
		return s[i].(float32) < s[j].(float32)
	case float64:
		return s[i].(float64) < s[j].(float64)
	case string:
		return s[i].(string) < s[j].(string)
	default:
		return true
	}
}

func isArrayEqual(v1 SortInterfaceArray, v2 SortInterfaceArray, strict bool) bool {
	if len(v1) != len(v2) {
		return false
	}

	if !strict {
		sort.Sort(v1)
		sort.Sort(v2)
	}

	for index, v := range v1 {
		switch v.(type) {
		case []interface{}:
			if !isArrayEqual(v.([]interface{}), v2[index].([]interface{}), strict) {
				return false
			}
		default:
			if v != v2[index] {
				return false
			}
		}
	}

	return true
}

func assertChanWithValues(t *testing.T, c chan interface{}, vals []interface{}) {
	inVals := Values(c)
	if !isArrayEqual(inVals, vals, true) {
		t.Errorf("assertChanWithValues fail, chan: %v, array: %v", inVals, vals)
	}
}

func assertEqual(t *testing.T, v1 interface{}, v2 interface{}) {
	if !reflect.DeepEqual(v1, v2) {
		t.Errorf("assertChanWithValues fail: %v, expected: %v", v1, v2)
	}
}

func timeoutContext(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}

func TestMap(t *testing.T) {
	in := From([]interface{}{1, 2, 3, 4})

	out := Map(in, func(i interface{}) interface{} {
		num := i.(int)
		return num * 2
	})

	assertChanWithValues(t, out, []interface{}{2, 4, 6, 8})
}

func TestReduce(t *testing.T) {
	in := From([]interface{}{1, 2, 3, 4})
	out := Reduce(in, func(i interface{}, i2 interface{}) interface{} {
		return i.(int) + i2.(int)
	}, 0)

	assertChanWithValues(t, out, []interface{}{10})
}

func TestFilter(t *testing.T) {
	in := From([]interface{}{1, 2, 3, 4})

	out := Filter(in, func(i interface{}) bool {
		num := i.(int)
		return num%2 == 0
	})

	assertChanWithValues(t, out, []interface{}{2, 4})
}

func TestMerge(t *testing.T) {
	c1 := From([]interface{}{1, 2, 9})
	c2 := From([]interface{}{3, 4, 10})
	c3 := From([]interface{}{5, 7, 11})
	c4 := From([]interface{}{6, 8, 12})

	ch := Merge(c1, c2, c3, c4)

	array := Values(ch)
	if !isArrayEqual(array, []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, false) {
		t.Error("merge values not equal to original ones")
	}
}

func TestSubscribe(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4}
	ch := From(arr)

	var wg sync.WaitGroup
	wg.Add(3)

	sub1 := Subscribe(ch)
	sub2 := Subscribe(ch)
	sub3 := Subscribe(ch)

	go func() {
		defer wg.Done()

		vals := Values(sub1)

		if !isArrayEqual(vals, arr, true) {
			t.Errorf("values not equal: %v, original: %v", vals, arr)
		}
	}()

	go func() {
		defer wg.Done()

		vals := Values(sub2)

		if !isArrayEqual(vals, arr, true) {
			t.Errorf("values not equal: %v, original: %v", vals, arr)
		}
	}()

	go func() {
		defer wg.Done()

		vals := Values(sub3)

		if !isArrayEqual(vals, arr, true) {
			t.Errorf("values not equal: %v, original: %v", vals, arr)
		}
	}()

	wg.Wait()
}

func TestBuffer(t *testing.T) {
	out1 := Buffer(From([]interface{}{1, 2, 3, 4, 5, 6, 7}), 2, 3)
	assertChanWithValues(t, out1, []interface{}{[]interface{}{1, 2}, []interface{}{4, 5}})

	out2 := Buffer(From([]interface{}{1, 2, 3, 4}), 2, 0)
	assertChanWithValues(t, out2, []interface{}{[]interface{}{1, 2}, []interface{}{3, 4}})
}

func TestFlatMap(t *testing.T) {
	in := From([]interface{}{[]interface{}{1, 2}, []interface{}{3, 4}})
	out := FlatMap(in)
	assertChanWithValues(t, out, []interface{}{1, 2, 3, 4})
}

func TestGroupBy(t *testing.T) {
	out1 := Values(GroupBy(From([]interface{}{1, 2, 3, 4}), func(i interface{}) interface{} {
		num := i.(int)
		return num%2 == 0
	}))

	assertEqual(t, out1[0], map[interface{}][]interface{}{
		true:  {2, 4},
		false: {1, 3},
	})

	out2 := Values(GroupBy(From([]interface{}{1, 2, 3, 4}), func(i interface{}) interface{} {
		num := i.(int)
		if num%2 == 0 {
			return "even"
		} else {
			return "odd"
		}
	}))

	assertEqual(t, out2[0], map[interface{}][]interface{}{
		"odd":  {1, 3},
		"even": {2, 4},
	})
}

func TestDebounce(t *testing.T) {
	in := IntervalRange(0, 4, 40*time.Millisecond, 0)
	out := Debounce(in, time.Millisecond*100)

	assertChanWithValues(t, out, []interface{}{0, 3})
}

func TestDistinct(t *testing.T) {
	in := From([]interface{}{1, 2, 2, 3, 3, 5})
	out := Distinct(in)

	assertChanWithValues(t, out, []interface{}{1, 2, 3, 5})
}

func TestElementAt(t *testing.T) {
	out1 := ElementAt(From([]interface{}{0, 1, 2, 3}), 2)
	assertChanWithValues(t, out1, []interface{}{2})

	out2 := ElementAt(From([]interface{}{0, 1, 2, 3}), 6)
	assertChanWithValues(t, out2, []interface{}{})

	out3 := ElementAt(From([]interface{}{0, 1, 2, 3}), -1)
	assertChanWithValues(t, out3, []interface{}{})
}

func TestFirst(t *testing.T) {
	out1 := First(From([]interface{}{0, 1, 2, 3}))
	assertChanWithValues(t, out1, []interface{}{0})

	out2 := First(From([]interface{}{}))
	assertChanWithValues(t, out2, []interface{}{})
}

func TestLast(t *testing.T) {
	out1 := Last(From([]interface{}{1, 2, 3}))
	assertChanWithValues(t, out1, []interface{}{3})

	out2 := Last(From([]interface{}{}))
	assertChanWithValues(t, out2, []interface{}{})
}

func TestIgnoreElements(t *testing.T) {
	out1 := IgnoreElements(From([]interface{}{1, 2, 3}))
	assertChanWithValues(t, out1, []interface{}{})

	out2 := IgnoreElements(From([]interface{}{}))
	assertChanWithValues(t, out2, []interface{}{})
}

func TestSample(t *testing.T) {
	in := IntervalRange(1, 7, time.Millisecond*30, 0)

	out := Sample(in, time.Millisecond*50)

	assertChanWithValues(t, out, []interface{}{1, 3, 4, 6})
}

func TestSkip(t *testing.T) {
	in := From([]interface{}{1, 2, 3, 4})
	out := Skip(in, 2)
	assertChanWithValues(t, out, []interface{}{3, 4})

	in2 := From([]interface{}{1, 2, 3, 4})
	out2 := Skip(in2, 5)
	assertChanWithValues(t, out2, []interface{}{})

	in3 := From([]interface{}{1, 2, 3, 4})
	out3 := Skip(in3, 0)
	assertChanWithValues(t, out3, []interface{}{1, 2, 3, 4})
}

func TestSkipLast(t *testing.T) {
	in1 := From([]interface{}{1, 2, 3, 4})
	out1 := SkipLast(in1, 2)
	assertChanWithValues(t, out1, []interface{}{1, 2})

	in2 := From([]interface{}{1, 2})
	out2 := SkipLast(in2, 3)
	assertChanWithValues(t, out2, []interface{}{})
}

func TestTake(t *testing.T) {
	in1 := From([]interface{}{1, 2, 3, 4})
	out1 := Take(in1, 0)
	assertChanWithValues(t, out1, []interface{}{})

	in2 := From([]interface{}{1, 2, 3, 4})
	out2 := Take(in2, 2)
	assertChanWithValues(t, out2, []interface{}{1, 2})

	in3 := From([]interface{}{1, 2})
	out3 := Take(in3, 3)
	assertChanWithValues(t, out3, []interface{}{1, 2})
}

func TestTakeLast(t *testing.T) {
	in1 := From([]interface{}{1, 2, 3, 4})
	out1 := TakeLast(in1, 0)
	assertChanWithValues(t, out1, []interface{}{})

	in2 := From([]interface{}{1, 2, 3, 4})
	out2 := TakeLast(in2, 2)
	assertChanWithValues(t, out2, []interface{}{3, 4})

	in3 := From([]interface{}{1, 2, 3, 4})
	out3 := TakeLast(in3, 5)
	assertChanWithValues(t, out3, []interface{}{1, 2, 3, 4})
}

func TestCombineLatest(t *testing.T) {
	in1 := IntervalRange(1, 4, time.Millisecond*30, 0)
	in2 := Map(IntervalRange(0, 2, time.Millisecond*50, 0), func(i interface{}) interface{} {
		return string(97 + i.(int))
	})

	out := Map(CombineLatest(in1, in2), func(i interface{}) interface{} {
		values := i.([]interface{})
		return strconv.Itoa(values[0].(int)) + values[1].(string)
	})

	assertChanWithValues(t, out, []interface{}{"1a", "2a", "3a", "3b", "4b"})
}

func TestStartWith(t *testing.T) {
	in1 := From([]interface{}{3, 4})
	out1 := StartWith(in1, 1, 2)
	assertChanWithValues(t, out1, []interface{}{1, 2, 3, 4})

	in2 := From([]interface{}{3, 4})
	out2 := StartWith(in2)
	assertChanWithValues(t, out2, []interface{}{3, 4})

	in3 := From([]interface{}{})
	out3 := StartWith(in3, 1, 2)
	assertChanWithValues(t, out3, []interface{}{1, 2})
}

func TestSwitch(t *testing.T) {
	var index int64 = 0
	in := Interval(timeoutContext(time.Millisecond*60), time.Millisecond*25, 0, func(i interface{}) interface{} {
		index++

		idx := index

		return Map(IntervalRange(1, 3, time.Millisecond*10, 0), func(i interface{}) interface{} {
			return strconv.FormatInt(idx, 10) + strconv.Itoa(i.(int))
		})
	})

	out := Switch(in)

	assertChanWithValues(t, out, []interface{}{"11", "12", "21", "22", "23"})
}

func TestZip(t *testing.T) {
	in11 := Range(0, 3)
	in12 := Range(5, 10)
	out1 := Map(Zip(in11, in12), func(i interface{}) interface{} {
		arr := i.([]interface{})
		return strconv.Itoa(arr[0].(int)) + strconv.Itoa(arr[1].(int))
	})

	assertChanWithValues(t, out1, []interface{}{"05", "16", "27"})

	in21 := Range(0, 0)
	in22 := Range(5, 10)
	out2 := Map(Zip(in21, in22), func(i interface{}) interface{} {
		arr := i.([]interface{})
		return strconv.Itoa(arr[0].(int)) + strconv.Itoa(arr[1].(int))
	})

	assertChanWithValues(t, out2, []interface{}{})

	in31 := Range(0, 3)
	in32 := Range(5, 10)
	in33 := Range(7, 2)
	out3 := Map(Zip(in31, in32, in33), func(i interface{}) interface{} {
		arr := i.([]interface{})
		return strconv.Itoa(arr[0].(int)) + strconv.Itoa(arr[1].(int)) + strconv.Itoa(arr[2].(int))
	})

	assertChanWithValues(t, out3, []interface{}{"057", "168"})
}

func TestDelay(t *testing.T) {
	in := IntervalRange(1, 3, time.Millisecond*10, 0)

	length := 0
	out := Delay(in, time.Millisecond*100)

	for val := range out {
		length++
		assertEqual(t, length, val)
	}

	assertEqual(t, length, 3)
}

func TestTimeInterval(t *testing.T) {
	in := Interval(timeoutContext(time.Millisecond*100), time.Millisecond*30, 0, nil)
	out := TimeInterval(in)
	vals := Values(out)

	assertEqual(t, len(vals), 3)
}

func TestTimeout(t *testing.T) {
	in := Interval(timeoutContext(time.Millisecond*100), time.Millisecond*30, 0, nil)

	out := Timeout(in, time.Millisecond*20)

	error, _ := <-out

	switch error.(type) {
	case TimeoutError:
	default:
		t.Errorf("got: %v, while expecting TimeoutError", error)
	}
}

func TestTimestamp(t *testing.T) {
	in := Interval(timeoutContext(time.Millisecond*50), time.Millisecond*20, 0, nil)
	out := Timestamp(in)
	vals := Values(out)

	assertEqual(t, len(vals), 2)

	v1 := vals[0].(TimestampItem)
	v2 := vals[1].(TimestampItem)
	diff := v2.timestamp.Sub(v1.timestamp) / time.Millisecond
	assert.T(t, diff <= 25 && diff >= 15)
}

func TestAll(t *testing.T) {
	in1 := From([]interface{}{1, 2, 3})
	out1 := All(in1, func(i interface{}) bool {
		return i.(int)%2 == 0
	})

	assertChanWithValues(t, out1, []interface{}{false})

	in2 := From([]interface{}{2, 4})
	out2 := All(in2, func(i interface{}) bool {
		return i.(int)%2 == 0
	})

	assertChanWithValues(t, out2, []interface{}{true})
}

func TestContains(t *testing.T) {
	in1 := From([]interface{}{1, 2, 3})
	out1 := Contains(in1, func(i interface{}) bool {
		return i.(int)%2 == 0
	})

	assertChanWithValues(t, out1, []interface{}{true})

	in2 := From([]interface{}{2, 4})
	out2 := Contains(in2, func(i interface{}) bool {
		return i.(int)%2 != 0
	})

	assertChanWithValues(t, out2, []interface{}{false})
}

func TestAmb(t *testing.T) {
	in1 := IntervalRange(0, 5, time.Millisecond*10, time.Millisecond*50)
	in2 := IntervalRange(10, 3, time.Millisecond*20, 0)

	out := Amb(in1, in2)

	assertChanWithValues(t, out, []interface{}{10, 11, 12})
}

func TestDefaultIfEmpty(t *testing.T) {
	in1 := From([]interface{}{})
	out1 := DefaultIfEmpty(in1, 1)
	assertChanWithValues(t, out1, []interface{}{1})

	in2 := From([]interface{}{2})
	out2 := DefaultIfEmpty(in2, 1)
	assertChanWithValues(t, out2, []interface{}{2})
}

func TestSequenceEqual(t *testing.T) {
	in11 := From([]interface{}{1, 2, 3, 4})
	in12 := IntervalRange(1, 4, time.Millisecond*10, 0)
	out1 := SequenceEqual(in11, in12)
	assertChanWithValues(t, out1, []interface{}{true})

	in21 := From([]interface{}{})
	in22 := From([]interface{}{})
	out2 := SequenceEqual(in21, in22)
	assertChanWithValues(t, out2, []interface{}{true})

	in31 := From([]interface{}{1, 2, 3, 4})
	in32 := From([]interface{}{1, 2})
	out3 := SequenceEqual(in31, in32)
	assertChanWithValues(t, out3, []interface{}{false})

	in41 := From([]interface{}{1, 2, 3, 4})
	in42 := From([]interface{}{})
	out4 := SequenceEqual(in41, in42)
	assertChanWithValues(t, out4, []interface{}{false})
}

func TestSkipUntil(t *testing.T) {
	skipUntil1 := IntervalRange(0, 1, time.Millisecond*35, 0)
	in1 := IntervalRange(0, 5, time.Millisecond*10, 0)
	out1 := SkipUntil(in1, skipUntil1)

	assertChanWithValues(t, out1, []interface{}{3, 4})

	skipUntil2 := IntervalRange(0, 1, time.Millisecond*100, 0)
	in2 := IntervalRange(0, 5, time.Millisecond*10, 0)
	out2 := SkipUntil(in2, skipUntil2)
	assertChanWithValues(t, out2, []interface{}{})
}

func TestSkipWhile(t *testing.T) {
	in1 := From([]interface{}{1, 2, 3, 4})
	out1 := SkipWhile(in1, func(i interface{}) bool {
		return i.(int) > 2
	})
	assertChanWithValues(t, out1, []interface{}{3, 4})

	in2 := From([]interface{}{1, 2, 3, 4})
	out2 := SkipWhile(in2, func(i interface{}) bool {
		return i.(int) > 10
	})
	assertChanWithValues(t, out2, []interface{}{})

	in3 := From([]interface{}{})
	out3 := SkipWhile(in3, func(i interface{}) bool {
		return i.(int) > 10
	})
	assertChanWithValues(t, out3, []interface{}{})
}

func TestTakeUntil(t *testing.T) {
	takeUntil1 := IntervalRange(0, 1, time.Millisecond*35, 0)
	in1 := IntervalRange(0, 5, time.Millisecond*10, 0)
	out1 := TakeUntil(in1, takeUntil1)
	assertChanWithValues(t, out1, []interface{}{0, 1, 2})

	takeUntil2 := IntervalRange(0, 1, time.Millisecond*100, 0)
	in2 := IntervalRange(0, 5, time.Millisecond*10, 0)
	out2 := TakeUntil(in2, takeUntil2)
	assertChanWithValues(t, out2, []interface{}{0, 1, 2, 3, 4})
}

func TestTakeWhile(t *testing.T) {
	in1 := From([]interface{}{1, 2, 3, 4})
	out1 := TakeWhile(in1, func(i interface{}) bool {
		return i.(int) <= 2
	})
	assertChanWithValues(t, out1, []interface{}{1, 2})

	in2 := From([]interface{}{1, 2, 3, 4})
	out2 := TakeWhile(in2, func(i interface{}) bool {
		return i.(int) <= 10
	})
	assertChanWithValues(t, out2, []interface{}{1, 2, 3, 4})

	in3 := From([]interface{}{})
	out3 := TakeWhile(in3, func(i interface{}) bool {
		return i.(int) <= 10
	})
	assertChanWithValues(t, out3, []interface{}{})
}

func TestSum(t *testing.T) {
	in1 := From([]interface{}{})
	out1 := Sum(in1)
	assertChanWithValues(t, out1, []interface{}{0})

	in2 := From([]interface{}{1})
	out2 := Sum(in2)
	assertChanWithValues(t, out2, []interface{}{1})

	in3 := From([]interface{}{1, 2, 3})
	out3 := Sum(in3)
	assertChanWithValues(t, out3, []interface{}{6})
}

func TestAverage(t *testing.T) {
	in1 := From([]interface{}{})
	out1 := Average(in1)
	assertChanWithValues(t, out1, []interface{}{})

	in2 := From([]interface{}{1})
	out2 := Average(in2)
	assertChanWithValues(t, out2, []interface{}{1})

	in3 := From([]interface{}{1, 2, 3})
	out3 := Average(in3)
	assertChanWithValues(t, out3, []interface{}{2})
}

func TestCount(t *testing.T) {
	in1 := From([]interface{}{})
	out1 := Count(in1)
	assertChanWithValues(t, out1, []interface{}{0})

	in2 := From([]interface{}{1, 2})
	out2 := Count(in2)
	assertChanWithValues(t, out2, []interface{}{2})
}

func TestMin(t *testing.T) {
	in1 := From([]interface{}{})
	out1 := Min(in1)
	assertChanWithValues(t, out1, []interface{}{})

	in2 := From([]interface{}{1, -2, 3})
	out2 := Min(in2)
	assertChanWithValues(t, out2, []interface{}{-2})

	in3 := From([]interface{}{1.0, -2.0, 3.0})
	out3 := Min(in3)
	assertChanWithValues(t, out3, []interface{}{-2.0})

	in4 := From([]interface{}{"a", "ab", "b"})
	out4 := Min(in4)
	assertChanWithValues(t, out4, []interface{}{"a"})
}

func TestMax(t *testing.T) {
	in1 := From([]interface{}{})
	out1 := Max(in1)
	assertChanWithValues(t, out1, []interface{}{})

	in2 := From([]interface{}{1, -2, 3})
	out2 := Max(in2)
	assertChanWithValues(t, out2, []interface{}{3})

	in3 := From([]interface{}{1.0, -2.0, 3.0})
	out3 := Max(in3)
	assertChanWithValues(t, out3, []interface{}{3.0})

	in4 := From([]interface{}{"a", "ab", "b"})
	out4 := Max(in4)
	assertChanWithValues(t, out4, []interface{}{"b"})
}

func TestConcat(t *testing.T) {
	in11 := From([]interface{}{1, 2})
	in12 := From([]interface{}{"a", "b"})
	out1 := Concat(in11, in12)
	assertChanWithValues(t, out1, []interface{}{1, 2, "a", "b"})

	in21 := From([]interface{}{})
	in22 := From([]interface{}{"b"})
	out2 := Concat(in21, in22)
	assertChanWithValues(t, out2, []interface{}{"b"})
}

func TestIsEmpty(t *testing.T) {
	in1 := From([]interface{}{})
	out1 := IsEmpty(in1)
	assertChanWithValues(t, out1, []interface{}{true})

	in2 := From([]interface{}{1})
	out2 := IsEmpty(in2)
	assertChanWithValues(t, out2, []interface{}{false})
}

func TestWindowWithCount(t *testing.T) {
	in1 := From([]interface{}{1, 2, 3})
	out1 := WindowWithCount(in1, 2)
	values1 := Values(out1)
	assertChanWithValues(t, values1[0].(chan interface{}), []interface{}{1, 2})
	assertChanWithValues(t, values1[1].(chan interface{}), []interface{}{3})

	in2 := From([]interface{}{1})
	out2 := WindowWithCount(in2, 2)
	values2 := Values(out2)
	assertChanWithValues(t, values2[0].(chan interface{}), []interface{}{1})
}
