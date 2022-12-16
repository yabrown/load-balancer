package balancer

import (
	"fmt"
	s "github.com/montanaflynn/stats"
)

// average response time (latency) per task
func MeasureResponseTime(b *Balancer) (float64, float64, float64) {
	// the balancer has an instance variable tracking the total wall-clock time
	// for each task that has been completed
	// divide that by the balance instance variable tracking the total number of completed tasks

	var totalTime float64 = 0
	var totalTasksFinished int = 0
	times := make([]float64, len(b.request_stats))
	for _, stats := range b.request_stats {
		// fmt.Println(req.id, stats.duration.Seconds())
		totalTime += float64(stats.duration.Milliseconds())
		times = append(times, float64(stats.duration.Milliseconds()))
		if stats.handled {
			totalTasksFinished += 1
		}
	}

	fmt.Println("Number of tasks:", totalTasksFinished)

	averageTime := totalTime / float64(totalTasksFinished)

	stdDevTime, _ := s.StandardDeviation(times)
	medianTime, _ := s.Median(times)

	return averageTime, stdDevTime, medianTime

}

// average load of a task

//compute the average load per task (sum over all servers the timeCompleted * load_per_core)/num of completed tasks
// need to check if the units for this are equal or not
func MeasureAverageTaskLoad(b *Balancer) float64 {
	var totalLoad float64 = 0
	for _, s := range b.servers {
		totalLoad += s.timeCompleted * float64(s.cores)
	}

	return totalLoad / float64(len(b.request_stats))

}

//compute the average load per task (sum over all servers the timeCompleted * load_per_core)/num of completed tasks
// need to check if the units for this are equal or not
func MeasureAverageTaskTime(b *Balancer) (float64, float64, float64) {
	var totalTime float64 = 0
	times := make([]float64, len(b.servers))
	for _, s := range b.servers {
		totalTime += s.timeCompleted
		times = append(times, s.timeCompleted)
	}

	averageTime := totalTime / float64(len(b.request_stats))
	stdDevTime, _ := s.StandardDeviation(times)
	medianTime, _ := s.Median(times)

	return averageTime, stdDevTime, medianTime

}

// load distribution
// look at each server, check the total load on that server proportional
// to the number of cores, and compare to the overall load on the balancer
func MeasureLoadDistribution(b *Balancer) (float64, float64, float64) {
	load_per_core := make([]float64, len(b.servers))
	load_per_server := make([]float64, len(b.servers))
	tasks_per_server := make([]float64, len(b.servers))
	for i, s := range b.servers {
		isOnline, _ := s.State()
		if isOnline {
			load_per_core[i] = s.timeCompleted
			load_per_server[i] = s.timeCompleted * float64(s.cores)
			tasks_per_server[i] = float64(s.total_tasks)
		}

	}
	fmt.Println("Total runtime for each server:\t", load_per_core)

	fmt.Println("Total load for each server:\t", load_per_server)
	fmt.Println("Total tasks for each server:\t", tasks_per_server)

	averageRuntime, _ := s.Mean(load_per_core)
	stdDevRuntime, _ := s.StandardDeviation(load_per_core)
	medianRuntime, _ := s.Median(load_per_core)
	return averageRuntime, stdDevRuntime, medianRuntime

}
