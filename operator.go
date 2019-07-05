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

/**
 * fan in
 */
func Merge(chans ... chan interface{}) chan interface{} {
	out := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(len(chans))

	for _, c := range chans {
		ch := c

		go func() {
			for {
				val, ok := <- ch
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
	var subs [] chan interface{}

	val, ok := subscriptionMap.Load(in)
	if !ok {
		subs = make([] chan interface{}, 0)
	} else {
		subs = val.([] chan interface{})
	}

	out := make(chan interface{})

	subs = append(subs, out)
	subscriptionMap.Store(in, subs)

	// listen to the channel for the first time
	if !ok {
		go func() {
			for {
				val, ok := <- in

				mapVal, _ := subscriptionMap.Load(in)
				currentSubs := mapVal.([] chan interface{})

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