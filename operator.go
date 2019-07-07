package reactive_channel

import (
	"sync"
	"time"
)

type MapFunc func(interface{}) interface{}

func Map(in <-chan interface{}, mapFunc MapFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		for {
			val, ok := <-in
			if ok {
				out <- mapFunc(val)
			} else {
				close(out)
				return
			}
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
		for {
			val, ok := <-in
			if ok {
				if filterFunc(val) {
					out <- val
				}
			} else {
				close(out)
				return
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
			for {
				val, ok := <-ch
				if ok {
					out <- val
				} else {
					wg.Done()
					return
				}
			}
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