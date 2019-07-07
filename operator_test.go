package reactive_channel

import (
	"reflect"
	"sort"
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
			if !isArrayEqual(v.([]interface{}), v2[index].([] interface{}), strict) {
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
	inVals := To(c)
	if !isArrayEqual(inVals, vals, true) {
		t.Errorf("assertChanWithValues fail, chan: %v, array: %v", inVals, vals)
	}
}

func assertEqual(t *testing.T, v1 interface{}, v2 interface{}) {
	if !reflect.DeepEqual(v1, v2) {
		t.Errorf("assertChanWithValues fail: %v, expected: %v", v1, v2)
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

func TestReduce(t *testing.T) {
	in := From([] interface{} {1,2,3,4})
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

	array := To(ch)
	if !isArrayEqual(array, []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, false) {
		t.Error("merge values not equal to original ones")
	}
}

func TestBroadcast(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4}
	ch := From(arr)

	var wg sync.WaitGroup
	wg.Add(3)

	sub1 := Broadcast(ch)
	sub2 := Broadcast(ch)
	sub3 := Broadcast(ch)

	go func() {
		defer wg.Done()

		vals := To(sub1)

		if !isArrayEqual(vals, arr, true) {
			t.Errorf("values not equal: %v, original: %v", vals, arr)
		}
	}()

	go func() {
		defer wg.Done()

		vals := To(sub2)

		if !isArrayEqual(vals, arr, true) {
			t.Errorf("values not equal: %v, original: %v", vals, arr)
		}
	}()

	go func() {
		defer wg.Done()

		vals := To(sub3)

		if !isArrayEqual(vals, arr, true) {
			t.Errorf("values not equal: %v, original: %v", vals, arr)
		}
	}()

	wg.Wait()
}

func TestBuffer(t *testing.T) {
	out1 := Buffer(From([]interface{}{1,2,3,4,5,6,7}), 2,3)
	assertChanWithValues(t, out1, []interface{}{[]interface{}{1,2}, []interface{}{4,5}})

	out2 := Buffer(From([]interface{}{1,2,3,4}), 2,0)
	assertChanWithValues(t, out2, []interface{}{[]interface{}{1,2}, []interface{}{3, 4}})
}

func TestFlatMap(t *testing.T) {
	in := From([]interface{}{[]interface{}{1,2}, []interface{}{3,4}})
	out := FlatMap(in)
	assertChanWithValues(t, out, []interface{}{1,2,3,4})
}

func TestGroupBy(t *testing.T) {
	out1 := To(GroupBy(From([]interface{}{1, 2, 3, 4}), func(i interface{}) interface{} {
		num := i.(int)
		return num % 2 == 0
	}))

	assertEqual(t, out1[0], map[interface{}][]interface{}{
		true:  {2, 4},
		false: {1, 3},
	})

	out2 := To(GroupBy(From([]interface{}{1, 2, 3, 4}), func(i interface{}) interface{} {
		num := i.(int)
		if num % 2 == 0 {
			return "even"
		} else {
			return "odd"
		}
	}))

	assertEqual(t, out2[0], map[interface{}][]interface{}{
		"odd": {1, 3},
		"even":  {2, 4},
	})
}

func TestDebounce(t *testing.T) {
	in := make(chan interface{})

	go func() {
		ticker := time.NewTicker(time.Millisecond *60)

		index := 0
		for {
			if index > 4 {
				ticker.Stop()
				close(in)
				return
			}

			<- ticker.C
			in <- index
			index++
		}
	}()

	out := Debounce(in, time.Millisecond * 100)

	assertChanWithValues(t, out, [] interface{} {0 ,2, 4})
}

func TestDistinct(t *testing.T) {
	in := From([]interface{} {1,2,2,3,3,5})
	out := Distinct(in)

	assertChanWithValues(t, out, []interface{}{1,2,3,5})
}

func TestElementAt(t *testing.T) {
	out1 := ElementAt(From([]interface{}{0,1,2,3}), 2)
	assertChanWithValues(t, out1, []interface{}{2})

	out2 := ElementAt(From([]interface{}{0,1,2,3}), 6)
	assertChanWithValues(t, out2, []interface{}{})


	out3 := ElementAt(From([]interface{}{0,1,2,3}), -1)
	assertChanWithValues(t, out3, []interface{}{})
}