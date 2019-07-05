package reactive_channel

import (
	"sort"
	"sync"
	"testing"
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
		if v != v2[index] {
			return false
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
