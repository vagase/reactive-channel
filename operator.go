package reactive_channel

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