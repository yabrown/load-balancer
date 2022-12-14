package balancer

import (
	"fmt"
)

// average response time (latency) per task
func MeasureAverageResponseTime(b *Balancer) float32 {
	// the balancer has an instance variable tracking the total wall-clock time
	// for each task that has been completed
	// divide that by the balance instance variable tracking the total number of completed tasks

	var totalTime float32 = 0
	var totalTasks int = 0
	for _, stats := range b.request_stats {
		// fmt.Println(req.id, stats.duration.Seconds())

		totalTime += float32(stats.duration.Milliseconds())

		if stats.handled {
			totalTasks += 1
		}
	}

	fmt.Println("Number of tasks:", totalTasks)

	averageTime := totalTime / float32(totalTasks)

	return averageTime

}

// average load of a task

//compute the average load per task (sum over all servers the timeCompleted * load_per_core)/num of completed tasks
// need to check if the units for this are equal or not
func MeasureAverageTaskLoad(b *Balancer) float32 {
	var totalLoad float32 = 0
	for _, s := range b.servers {
		totalLoad += s.timeCompleted * float32(s.cores)
	}

	return totalLoad / float32(len(b.request_stats))

}

//compute the average load per task (sum over all servers the timeCompleted * load_per_core)/num of completed tasks
// need to check if the units for this are equal or not
func MeasureAverageTaskTime(b *Balancer) float32 {
	var totalTime float32 = 0
	for _, s := range b.servers {
		totalTime += s.timeCompleted
	}

	return totalTime / float32(len(b.request_stats))

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
	fmt.Println("Total runtime for each server:\t", load_per_core)

	fmt.Println("Total load for each server:\t", load_per_server)

}
