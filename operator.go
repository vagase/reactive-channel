package reactive_channel

import (
	"context"
	"errors"
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

type ReduceFunc func(interface{}, interface{}) interface{}

func Reduce(in <-chan interface{}, reduceFunc ReduceFunc, initialVal interface{}) chan interface{} {
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
		for val := range in {
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
			for val := range ch {
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

		var buffer []interface{}
		counter := 0

		for val := range in {
			counter++
			if skip > 0 && counter%skip == 0 {
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

func FlatMap(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for val := range in {
			list := val.([]interface{})
			for _, item := range list {
				out <- item
			}
		}
	}()

	return out
}

type GroupByFunc func(interface{}) interface{}

func GroupBy(in chan interface{}, groupByFunc GroupByFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		result := map[interface{}][]interface{}{}
		defer func() {
			out <- result
			close(out)
		}()

		for val := range in {
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

		for val := range in {
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
		for val := range in {
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
			case val, ok := <-in:
				if ok {
					valToEmit = val
				} else {
					return
				}
			case <-ticker.C:
				if valToEmit != nil {
					out <- valToEmit
					valToEmit = nil
				}
			}
		}
	}()

	return out
}

func Skip(in chan interface{}, count int) chan interface{} {
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

		var cache []interface{}

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

		var cache []interface{}

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

func CombineLatest(chans ...chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		cases := make([]reflect.SelectCase, len(chans))
		for i, ch := range chans {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}

		values := make([]interface{}, len(chans))

		for {
			chosen, value, ok := reflect.Select(cases)
			if ok {
				values[chosen] = value.Interface()

				nilFound := false
				for _, v := range values {
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

func StartWith(in chan interface{}, vals ...interface{}) chan interface{} {
	out := make(chan interface{}, len(vals)+1)

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

func Switch(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		var cancelSubChan context.CancelFunc
		var subChanContext context.Context

		for newChan := range in {
			if cancelSubChan != nil {
				cancelSubChan()
				cancelSubChan = nil
			}

			var cancelSubChanContext context.Context
			cancelSubChanContext, cancelSubChan = context.WithCancel(context.Background())

			var subChanCancel context.CancelFunc
			subChanContext, subChanCancel = context.WithCancel(context.Background())

			go func() {
				for {
					select {
					case val, ok := <-newChan.(chan interface{}):
						if ok {
							out <- val
						} else {
							subChanCancel()
							return
						}
					case <-cancelSubChanContext.Done():
						return
					}
				}
			}()
		}

		<-subChanContext.Done()

		close(out)
	}()

	return out
}

func Zip(chans ...chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		valueCache := make([][]interface{}, len(chans))
		tryEmit := func() error {
			emit := true
			for _, arr := range valueCache {
				if len(arr) == 0 {
					emit = false
					break
				}
			}

			if !emit {
				return nil
			}

			emitVals := make([]interface{}, len(chans))
			for index, arr := range valueCache {
				if arr[0] == nil {
					return errors.New("end of chan")
				}

				emitVals[index] = arr[0]
				valueCache[index] = arr[1:]
			}

			out <- emitVals

			return nil
		}

		cases := make([]reflect.SelectCase, len(chans))
		for i, ch := range chans {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}

		for {
			index, value, ok := reflect.Select(cases)
			if ok {
				valueCache[index] = append(valueCache[index], value.Interface())
			} else {
				valueCache[index] = append(valueCache[index], nil)
			}

			if err := tryEmit(); err != nil {
				return
			}
		}
	}()

	return out
}

func Delay(in chan interface{}, delay time.Duration) chan interface{} {
	out := make(chan interface{})

	go func() {
		for {
			val, ok := <-in

			time.AfterFunc(delay, func() {
				if ok {
					out <- val
				} else {
					close(out)
				}
			})

			if !ok {
				return
			}
		}
	}()

	return out
}

type TimeIntervalItem struct {
	value    interface{}
	interval time.Duration
}

func TimeInterval(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		var lastTime *time.Time

		for val := range in {
			var interval time.Duration = 0

			now := time.Now()

			if lastTime != nil {
				interval = now.Sub(*lastTime)
			}

			lastTime = &now

			out <- TimeIntervalItem{val, interval}
		}
	}()

	return out
}

type TimeoutError struct {
}

func (timeoutError TimeoutError) Error() string {
	return "timeout error"
}

func Timeout(in chan interface{}, timeout time.Duration) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		select {
		case val, ok := <-in:
			if ok {
				out <- val
			} else {
				return
			}
		case <-time.After(timeout):
			out <- TimeoutError{}
			return
		}
	}()

	return out
}

type TimestampItem struct {
	value     interface{}
	timestamp time.Time
}

func Timestamp(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for val := range in {
			item := TimestampItem{val, time.Now()}
			out <- item
		}
	}()

	return out
}

type MatchFunc func(interface{}) bool

func All(in chan interface{}, match MatchFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for val := range in {
			if !match(val) {
				out <- false
				return
			}
		}

		out <- true
	}()

	return out
}

func Contains(in chan interface{}, match MatchFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for val := range in {
			if match(val) {
				out <- true
				return
			}
		}

		out <- false
	}()

	return out
}

func Amb(in ...chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		chans := in

		removeChanAtIndex := func(index int) {
			chans = append(chans[:index], chans[index+1:]...)
		}

		listen := func() (int, reflect.Value, bool) {
			cases := make([]reflect.SelectCase, len(chans))
			for i, ch := range chans {
				cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
			}

			return reflect.Select(cases)
		}

		var targetChan chan interface{}

		for {
			index, value, ok := listen()

			if ok {
				targetChan = chans[index]
				out <- value.Interface()
				break
			} else {
				removeChanAtIndex(index)

				if len(chans) == 0 {
					return
				}
			}
		}

		for val := range targetChan {
			out <- val
		}
	}()

	return out
}

func DefaultIfEmpty(in chan interface{}, defaultVal interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		empty := true

		for val := range in {
			out <- val
			empty = false
		}

		if empty {
			out <- defaultVal
		}
	}()

	return out
}

func SequenceEqual(ins ...chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		chans := ins

		valueCaches := make([][]interface{}, len(chans))
		isNotEqual := func() bool {
			var lastVal interface{}

			allEqual := true

			for _, arr := range valueCaches {
				if len(arr) > 0 {
					if lastVal == nil {
						lastVal = arr[0]
					} else if lastVal != arr[0] {
						return true
					}
				} else {
					allEqual = false
				}
			}

			if allEqual {
				for index, arr := range valueCaches {
					valueCaches[index] = arr[1:]
				}
			}

			return false
		}

		var wg sync.WaitGroup
		wg.Add(len(chans))

		for i, c := range chans {
			index := i
			ch := c

			go func() {
				defer wg.Done()

				for val := range ch {
					valueCaches[index] = append(valueCaches[index], val)
					if isNotEqual() {
						out <- false
						return
					}
				}
			}()
		}

		wg.Wait()

		for _, arr := range valueCaches {
			if len(arr) > 0 {
				out <- false
				return
			}
		}

		out <- true
	}()

	return out
}

func SkipUntil(in chan interface{}, untilChan chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		accept := false

		for {
			select {
			case val, ok := <- in:
				if ok {
					if accept {
						out <- val
					}
				} else {
					return
				}
			case <- untilChan:
				accept = true
			}
		}
	}()

	return out
}

func SkipWhile(in chan interface{}, match MatchFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		accept := false

		for val := range in {
			if accept {
				out <- val
			} else if match(val) {
				accept = true
				out <- val
			}
		}
	}()

	return out
}

func TakeUntil(in chan interface{}, untilChan chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for {
			select {
			case val, ok := <- in:
				if ok {
					out <- val
				} else {
					return
				}
			case <- untilChan:
				return
			}
		}
	}()

	return out
}

func TakeWhile(in chan interface{}, match MatchFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for val := range in {
			if !match(val) {
				return
			}

			out <- val
		}
	}()

	return out
}

func Sum (in chan interface{}) chan interface{} {
	return Reduce(in, func(i interface{}, i2 interface{}) interface{} {
		return i.(int) + i2.(int)
	}, 0)
}

func Average(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		count := 0
		sum, _ := <- Reduce(in, func(i interface{}, i2 interface{}) interface{} {
			count++
			return i.(int) + i2.(int)
		}, 0)

		if count > 0 {
			avg := sum.(int) / count
			out <- avg
		}
	}()

	return out
}

func Count(in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		count := 0

		for _ = range in {
			count ++
		}

		out <- count
	} ()

	return out
}

func Min(in chan interface{}) chan interface {} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		var min interface {}

		for val := range in {
			if min == nil {
				min = val
			} else {
				less := false
				switch val.(type) {
				case string:
					less = val.(string) < min.(string)
				case int:
					less = val.(int) < min.(int)
				case int8:
					less = val.(int8) < min.(int8)
				case int16:
					less = val.(int16) < min.(int16)
				case int32:
					less = val.(int32) < min.(int32)
				case int64:
					less = val.(int64) < min.(int64)
				case uint:
					less = val.(uint) < min.(uint)
				case uint8:
					less = val.(uint8) < min.(uint8)
				case uint16:
					less = val.(uint16) < min.(uint16)
				case uint32:
					less = val.(uint32) < min.(uint32)
				case uint64:
					less = val.(uint64) < min.(uint64)
				case float32:
					less = val.(float32) < min.(float32)
				case float64:
					less = val.(float64) < min.(float64)
				}

				if less {
					min = val
				}
			}
		}

		if min != nil {
			out <- min
		}
	} ()

	return out
}