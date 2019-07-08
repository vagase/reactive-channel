package reactive_channel

import (
	"reflect"
	"sync"
	"time"
)

type MapFunc func(interface{}) interface{}

func Map(in <-chan interface{}, mapFunc MapFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for val := range in {
			out <- mapFunc(val)
		}
	}()

	return out
}

type ReduceFunc func(interface{}, interface{}) interface {}
func Reduce(in <- chan interface{}, reduceFunc ReduceFunc, initialVal interface{}) chan interface {} {
	out := make(chan interface{})

	go func() {
		result := initialVal

		for val := range in {
			result = reduceFunc(result, val)
		}

		out <- result
		close(out)
	}()

	return out
}

type FilterFunc func(interface{}) bool

func Filter(in <-chan interface{}, filterFunc FilterFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for val := range in{
			if filterFunc(val) {
				out <- val
			}
		}
	}()

	return out
}

/**
 * fan in
 */
func Merge(chans ...chan interface{}) chan interface{} {
	out := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(len(chans))

	for _, c := range chans {
		ch := c

		go func() {
			for val := range ch{
				out <- val
			}

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

var subscriptionMap = &sync.Map{}

/**
 * fan out
 */
func Broadcast(in chan interface{}) chan interface{} {
	var subs []chan interface{}

	val, ok := subscriptionMap.Load(in)
	if !ok {
		subs = make([]chan interface{}, 0)
	} else {
		subs = val.([]chan interface{})
	}

	out := make(chan interface{})

	subs = append(subs, out)
	subscriptionMap.Store(in, subs)

	// listen to the channel for the first time
	if !ok {
		go func() {
			for {
				val, ok := <-in

				mapVal, _ := subscriptionMap.Load(in)
				currentSubs := mapVal.([]chan interface{})

				if ok {
					// broadcast
					for _, sub := range currentSubs {
						sub <- val
					}
				} else {
					// unsubscribe all
					for _, sub := range currentSubs {
						close(sub)
					}

					subscriptionMap.Delete(in)

					return
				}
			}
		}()
	}

	return out
}

func Buffer(in chan interface{}, count int, skip int) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		var buffer [] interface{}
		counter := 0

		for val := range in {
			counter++
			if skip > 0 && counter % skip == 0 {
				continue
			}

			buffer = append(buffer, val)

			if len(buffer) == count {
				copyBuffer := make([]interface{}, count)
				copy(copyBuffer, buffer)
				out <- copyBuffer
				buffer = buffer[:0]
			}
		}
	}()

	return out
}

func FlatMap(in chan interface{}) chan interface {} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for val := range in {
			list := val.([] interface{})
			for _, item := range list {
				out <- item
			}
		}
	}()

	return out
}

type GroupByFunc func(interface{}) interface{}

func GroupBy(in chan interface{}, groupByFunc GroupByFunc ) chan interface{} {
	out := make(chan interface{})

	go func() {
		result := map[interface{}] []interface{}{}
		defer func() {
			out <- result
			close(out)
		}()

		for val := range in{
			key := groupByFunc(val)
			result[key] = append(result[key], val)
		}
	}()

	return out
}

func Debounce(in chan interface{}, duration time.Duration) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		var lastTime *time.Time

		for val := range in{
			now := time.Now()

			if lastTime != nil && now.Sub(*lastTime) < duration {
				continue
			}

			lastTime = &now

			out <- val
		}
	}()

	return out
}

func Distinct(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		var lastVal interface{}
		for val := range in {
			if val == lastVal {
				continue
			}

			lastVal = val
			out <- val
		}
	}()

	return out
}

func ElementAt(in chan interface{}, nth int) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		if nth < 0 {
			return
		}

		index := 0
		for val := range in {
			if index == nth {
				out <- val
				return
			}

			index++
		}
	}()

	return out
}

func First(in chan interface{}) chan interface{} {
	return ElementAt(in, 0)
}

func Last(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		var lastVal interface{}
		for val:= range in {
			lastVal = val
		}

		if lastVal != nil {
			out <- lastVal
		}
	}()

	return out
}

func IgnoreElements(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for range in {
			// do nothing
		}
	}()

	return out
}

func Sample(in chan interface{}, interval time.Duration) chan interface{} {
	out := make(chan interface{})

	go func() {
		ticker := time.NewTicker(interval)

		defer func() {
			close(out)
			ticker.Stop()
		}()

		var valToEmit interface{}

		for {
			select {
			case val, ok := <- in:
				if ok {
					valToEmit = val
				} else {
					return
				}
			case <- ticker.C:
				if valToEmit != nil {
					out <- valToEmit
					valToEmit = nil
				}
			}
		}
	}()

	return out
}

func Skip (in chan interface{}, count int) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		index := 0
		for val := range in {
			if index >= count {
				out <- val
			}
			index++
		}
	}()

	return out
}

func SkipLast(in chan interface{}, count int) chan interface{} {
	if count <= 0 {
		return Map(in, func(i interface{}) interface{} {
			return i
		})
	}

	out := make(chan interface{})

	go func() {
		defer close(out)

		var cache []interface {}

		for val := range in {
			if len(cache) == count {
				out <- cache[0]
				cache = append(cache[1:], val)
			} else {
				cache = append(cache, val)
			}
		}
	}()

	return out
}

func Take(in chan interface{}, count int) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		index := 0
		for val := range in {
			index++

			if index <= count {
				out <- val
			} else {
				return
			}
		}
	}()

	return out
}

func TakeLast(in chan interface{}, count int) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		var cache [] interface{}

		for val := range in {
			cache = append(cache, val)
			if len(cache) > count {
				cache = cache[1:]
			}
		}

		for _, val := range cache {
			out <- val
		}
	}()

	return out
}

func CombineLatest(chans... chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		cases := make([]reflect.SelectCase, len(chans))
		for i, ch := range chans {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}

		values := make([] interface{}, len(chans))

		for {
			chosen, value, ok := reflect.Select(cases)
			if ok {
				values[chosen] = value.Interface()

				nilFound := false
				for _, v := range values  {
					if v == nil {
						nilFound = true
						break
					}
				}

				if !nilFound {
					out <- values
				}
			} else {
				return
			}
		}
	}()

	return out
}

func StartWith(in chan interface{}, vals... interface{}) chan interface{} {
	out := make(chan interface{}, len(vals) + 1)

	for _, val := range vals {
		out <- val
	}

	go func() {
		defer close(out)

		for val := range in {
			out <- val
		}
	}()

	return out
}