package balancer

import (
	"fmt"
	"sort"
)

// average response time (latency) per task
func MeasureAverageResponseTime(b *Balancer) float64 {
	// the balancer has an instance variable tracking the total wall-clock time
	// for each task that has been completed
	// divide that by the balance instance variable tracking the total number of completed tasks

	var totalTime float64 = 0
	var totalTasksFinished int = 0
	for _, stats := range b.request_stats {
		// fmt.Println(req.id, stats.duration.Seconds())
		totalTime += float64(stats.duration.Milliseconds())

		if stats.handled {
			totalTasksFinished += 1
		}
	}

	fmt.Println("Number of tasks:", totalTasksFinished)

	averageTime := totalTime / float64(totalTasksFinished)

	return averageTime

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
func MeasureAverageTaskTime(b *Balancer) float64 {
	var totalTime float64 = 0
	for _, s := range b.servers {
		totalTime += s.timeCompleted
	}

	return totalTime / float64(len(b.request_stats))

}

// load distribution
// look at each server, check the total load on that server proportional
// to the number of cores, and compare to the overall load on the balancer
func MeasureLoadDistribution(b *Balancer) {
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

}

func MeasureMedianResponseTime(b *Balancer) (float64, int) {

	times := make([]float64, len(b.request_stats))
	totalTasksFinished := 0
	for _, stats := range b.request_stats {
		// fmt.Println(req.id, stats.duration.Seconds())
		times = append(times, float64(stats.duration.Milliseconds()))

		if stats.handled {
			totalTasksFinished += 1
		}
	}

	sortedTimes := make([]float64, len(times))
	copy(sortedTimes, times)

	sort.Float64s(sortedTimes)

	l := len(sortedTimes)
	median := 0.0
	if l == 0 {
		return 0, 0
	} else if l%2 == 0 {
		median = (sortedTimes[l/2-1] + sortedTimes[l/2]) / 2
	} else {
		median = sortedTimes[l/2]
	}

	return median, totalTasksFinished

}
