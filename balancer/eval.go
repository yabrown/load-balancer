package balancer

import (
	"fmt"
	"time"
)

// response time per request
// send a request, set a timer, see how long it takes to return?
func MeasureResponseTime(b *Balancer) time.Duration {
	startTime := time.Now()
	task := NewRequest(-1, 100000)
	b.Assign_request_round_robin(task)
	// use channels to get a message back when the task is finished?
	elapsed := time.Now().Sub(startTime)
	return time.Duration(elapsed.Milliseconds())

}

// load distribution
// look at each server, check the total load on that server proportional
// to the number of cores, and compare to the overall load on the balancer
func MeasureLoadDistribution(b *Balancer) {
	load_per_core := make([]float32, len(b.servers))
	load_per_server := make([]float32, len(b.servers))
	for i, s := range b.servers {
		isOnline, _ := s.State()
		if isOnline {
			load_per_core[i] = s.timeCompleted
			load_per_server[i] = s.timeCompleted * float32(s.cores)
		}

	}
	fmt.Println("Load per core for each server: ", load_per_core)
	fmt.Println("-------------------------------------------")

	fmt.Println("Total load for each server", load_per_server)

}
