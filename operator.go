package reactive_channel

import "sync"

func From(in []interface{}) chan interface{} {
	out := make(chan interface{}, len(in))

	for _, val := range in {
		out <- val
	}

	close(out)

	return out
}

func To(in chan interface{}) [] interface{} {
	out := make([]interface{}, 0)
	for val := range in {
		out = append(out, val)
	}

	return out
}

type MapFunc func(interface{} ) interface{}

func Map(in <- chan interface{}, mapFunc MapFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		for {
			val, ok := <- in
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

type FilterFunc func(interface{}) bool

func Filter(in <- chan interface{}, filterFunc FilterFunc) chan interface{} {
	out := make(chan interface{})

	go func() {
		for {
			val, ok := <- in
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

func Merge(chans ... chan interface{}) chan interface{} {
	out := make(chan interface{})

	mutex := sync.Mutex{}
	closedCount := 0

	for _, c := range chans {
		ch := c

		go func() {
			for {
				val, ok := <- ch
				if ok {
					out <- val
				} else {
					mutex.Lock()
					closedCount++

					if closedCount == len(chans) {
						close(out)
					}

					mutex.Unlock()
					return
				}
			}
		}()
	}

	return out
}