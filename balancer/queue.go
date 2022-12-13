package balancer

import (
	"fmt"
)

var debug bool = false

////////////////////// QUEUE TYPE ////////////////////////////
type queue []*Request

func newQueue() queue {
	q := make(queue, 0)
	return q
}

func (queue queue) push(request *Request) queue {
	queue = append(queue, request)
	return queue
}

func (queue queue) pop() (*Request, queue) {
	// Testing
	if len(queue) == 0 {
		if debug {
			fmt.Print("tried to pop empty queue")
		}
		return nil, queue
	}
	var result = queue[0]
	queue = queue[1:]
	return result, queue
}

// Removes Request, returns queue and whether it was successful
func (queue queue) remove(request Request) (queue, bool) {
	var found bool = false
	var index int
	// Sets index to index of found request, otherwise 'found' stays false
	for i, element := range queue {
		if element.id == request.id {
			found = true
			index = i
		}
	}
	// If request was never found, return false now
	if !found {
		return nil, false
	}
	// Otherwise, request was found...
	// If it isn't the last item, push everything after it down one
	if index != len(queue) {
		// loop from request index to 1 less than last index
		for i := index; i < len(queue)-1; i++ {
			queue[i] = queue[i+1]
		}
	}
	// Then cut off last index (it's either extra or the target request)
	queue = queue[:len(queue)-1]
	return queue, true
}

/////////////////////////////////////////////////////////////
