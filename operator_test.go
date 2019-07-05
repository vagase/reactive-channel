package reactive_channel

import "testing"

func isArrayEqual(v1[]interface{}, v2[]interface{}) bool {
	if len(v1) != len(v1) {
		return false
	}

	for index, v := range v1{
		if v != v2[index] {
			return false
		}
	}

	return true
}

func assertChanWithValues(t *testing.T, c chan interface{}, vals []interface{}) {
	inVals := To(c)
	if !isArrayEqual(inVals, vals) {
		t.Errorf("assertChanWithValues fail, chan: %v, array: %v", inVals, vals)
	}
}

func TestFromTo(t *testing.T) {
	in := []interface{}{1,2,3,4}
	c := From(in)
	out := To(c)

	if !isArrayEqual(in, out) {
		t.Error("From or To failed")
	}
}

func TestMap (t *testing.T) {
	in := From([]interface{}{1,2,3,4})

	out := Map(in, func(i interface{}) interface{} {
		num := i.(int)
		return num * 2
	})

	assertChanWithValues(t, out, []interface{}{2,4,6,8})
}