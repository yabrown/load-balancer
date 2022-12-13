package balancer

import (
	"fmt"
	"time"
)

// average response time (latency) per task
func MeasureResponseTime(b *Balancer) (float32, float32) {
	// the balancer has an instance variable tracking the total wall-clock time
	// for each task that has been completed
	// divide that by the balance instance variable tracking the total number of completed tasks

	// also compute the average load per task (sum over all servers the timeCompleted * load_per_core)/num of completed tasks
	// need to check if the units for this are equal or not
	for len(b.acks) != 0 {

	}

	elapsedTime := b.totalTime / float32(b.totalTasks)

	return elapsedTime, 0

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
