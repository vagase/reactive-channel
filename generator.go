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

func Interval(ctx context.Context, interval time.Duration) chan interface{} {
	out := make(chan interface{})

	go func() {
		ticker := time.NewTicker(interval).C

		defer close(out)

		for {
			select {
			case val, ok := <-ticker:
				if ok {
					out <- val
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
			for {
				val, ok := <- in
				if ok {
					for index := 0; index < repeat; index ++ {
						out <- val
					}
				} else {
					return
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