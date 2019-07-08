package reactive_channel

import (
	"context"
	"time"
)

func From(in []interface{}) chan interface{} {
	out := make(chan interface{}, len(in))

	for _, val := range in {
		out <- val
	}

	close(out)

	return out
}

func Values(in chan interface{}) []interface{} {
	out := make([]interface{}, 0)
	for val := range in {
		out = append(out, val)
	}

	return out
}

func Interval(ctx context.Context, interval time.Duration, delay time.Duration, mapFunc MapFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		if delay > 0 {
			time.Sleep(delay)
		}

		ticker := time.NewTicker(interval)

		defer func() {
			close(out)
			ticker.Stop()
		}()

		for {
			select {
			case val, ok := <-ticker.C:
				if ok {
					var outVal interface{} = val
					if mapFunc != nil {
						outVal = mapFunc(val)
					}
					out <- outVal
				} else {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func IntervalRange(start int, count int, interval time.Duration, delay time.Duration) chan interface{} {
	timeout := interval*time.Duration(count) + time.Duration(float64(interval)*0.5)
	timeoutContext, _ := context.WithTimeout(context.Background(), timeout)

	index := 0
	return Interval(timeoutContext, interval, delay, func(i interface{}) interface{} {
		res := start + index

		index++

		return res
	})
}

func Range(start int, count int) chan interface{} {
	out := make(chan interface{}, count)

	for index := 0; index < count; index++ {
		out <- start + index
	}

	close(out)

	return out
}

/**
 * @immediate:
	- true: in(1->2->3), out (1->1->2->2->3->3)
	- false: in(1->2->3), out (1->2->3->1->2->3)
*/
func Repeat(in chan interface{}, repeat int, immediate bool) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		if immediate {
			for val := range in {
				for index := 0; index < repeat; index++ {
					out <- val
				}
			}
		} else {
			vals := Values(in)

			for index := 0; index < repeat; index++ {
				for _, val := range vals {
					out <- val
				}
			}
		}
	}()

	return out
}

func Just(val interface{}) chan interface{} {
	return From([] interface{} {val})
}

func Empty() chan interface {} {
	out := make(chan interface{})
	close(out)
	return out
}
