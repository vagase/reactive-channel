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
