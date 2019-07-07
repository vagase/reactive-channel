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

func To(in chan interface{}) []interface{} {
	out := make([]interface{}, 0)
	for val := range in {
		out = append(out, val)
	}

	return out
}

func Interval(ctx context.Context, interval time.Duration, mapFunc MapFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		ticker := time.NewTicker(interval)

		defer func() {
			close(out)
			ticker.Stop()
		}()

		for {
			select {
			case val, ok := <- ticker.C:
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

func Range(start int, size int) chan interface{} {
	out := make(chan interface{}, size)

	for index := 0 ; index < size; index++ {
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
				for index := 0; index < repeat; index ++ {
					out <- val
				}
			}
		} else {
			vals := To(in)

			for index := 0; index < repeat; index++ {
				for _, val := range vals{
					out <- val
				}
			}
		}
	}()

	return out
}