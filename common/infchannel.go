package common

import "container/list"

func MakeInfinite() (chan <-interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})

	go func() {
		inQueue := list.List{}

		outCh := func() chan interface{} {
			if inQueue.Len() == 0 {
				return nil
			}
			return out
		}

		curVal := func() interface{} {
			if inQueue.Len() == 0 {
				return nil
			}
			return inQueue.Front().Value
		}

		for inQueue.Len() > 0 || in != nil {
			select {
			case v, ok := <-in:
				if !ok {
					in = nil
				} else {
					inQueue.PushBack(v)
				}
			case outCh() <- curVal():
				inQueue.Remove(inQueue.Front())
			}
		}
		close(out)
	}()

	return in, out
}

