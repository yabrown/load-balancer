package main

import (
	b "cos316-load-balancer/balancer"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

func system_info(servers []*b.Server, balancer *b.Balancer) {
	print("\n")
	for i := range servers {
		fmt.Printf("\tServer status: %+v\n", *servers[i]) // formatting dereferences and prints fields
	}
	fmt.Printf("Balancer status: %+v\n\n", balancer) // formatting dereferences and prints fields
}

// when running main.go, use the following form:
// go run main.go [algorithm name: {"robin", "state"}] [time for servers to complete (int)] [# of tasks to run (int)]
func scratch() {
	var wg sync.WaitGroup
	alg := os.Args[1]
	numberOfTasks, _ := strconv.Atoi(os.Args[2])
	taskLoadRange, _ := strconv.Atoi(os.Args[3])
	//timeToSleep, _ := time.ParseDuration(os.Args[2] + "s")
	wg.Add(1)
	// create 5 servers
	s1 := b.NewServer(0, 1)
	s2 := b.NewServer(1, 1)
	s3 := b.NewServer(2, 1)
	s4 := b.NewServer(3, 1)
	s5 := b.NewServer(4, 1)
	// create balancer
	servers := []*b.Server{s1, s2, s3, s4, s5}
	end_signal := make(chan bool)
	balancer := b.NewBalancer(servers, end_signal)
	if balancer == nil {
	}
	// create requests
	requests := []*b.Request{}

	for i := 0; i < numberOfTasks; i++ {
		random := rand.Intn(taskLoadRange)
		if random < 1 {
			random += 1
		}
		requests = append(requests, b.NewRequest(i, random))
	}
	//wake all servers up
	for i := range servers {
		servers[i].Wake_up()
	}
	// assign each request
	for i := range requests {
		if alg == "robin" {
			balancer.Assign_request_round_robin(requests[i])
		} else if alg == "state" {
			balancer.Assign_request(requests[i])
		}
	}
	// system_info(servers, balancer)
	//make each server active (start handling queued requests)
	for i := range servers {
		go servers[i].Work(balancer)
	}
	fmt.Println("Completing tasks with algorithm option", alg)
	//fmt.Println("Sleeping for", timeToSleep)
	<-end_signal
	// system_info(servers, balancer)
	//fmt.Printf("Average Response Time: %.3fms\n", b.MeasureAverageResponseTime(balancer))
	//fmt.Printf("Average Task Load: %.3fms\n", b.MeasureAverageTaskLoad(balancer))
	//fmt.Printf("Average Task Time: %.3fms\n", b.MeasureAverageTaskTime(balancer))
	b.MeasureLoadDistribution(balancer)

	wg.Wait()
}

func graph() {
	// println("***********************************************************")
	// println("Equal cores, even distribution:\n")
	// const_cores_const_requests("robin")
	// const_cores_const_requests("state")
	// println("***********************************************************")
	// println("Equal cores, uneven distribution:\n")
	// const_cores_linear_requests("robin")
	// const_cores_linear_requests("state")
	// println("***********************************************************")
	// println("Equal cores, random distribution:\n")
	// const_cores_random_requests("robin")
	// const_cores_random_requests("state")
	// println("***********************************************************")
	// println("Linear cores, even distribution:\n")
	// linear_cores_const_requests("robin")
	// linear_cores_const_requests("state")
	// println("***********************************************************")
	// println("Linear cores, reverse linear distribution:\n")
	// linear_cores_reverse_requests("robin")
	// linear_cores_reverse_requests("state")
	// println("***********************************************************")
	// println("Linear cores, random distribution:\n")
	// linear_cores_random_requests("robin")
	// linear_cores_random_requests("state")
	// taskData := make([]int, 0)
	// algData := make([]string, 0)
	// avgResponseData := make([]float64, 0)
	// stdResponseData := make([]float64, 0)
	// avgLoadData := make([]float64, 0)
	// stdLoadData := make([]float64, 0)

	same_reqs, _ := os.Create("diff_cores_same_requests_taskData.csv")
	diff_reqs, err := os.Create("diff_cores_diff_requests_taskData.csv")
	if err != nil {
		panic(err)
	}

	// write data to 2 CSV files as we run the algorithms, so that the files can be exported to a Jupyter notebook
	same_reqs.WriteString("Number of Tasks,Algorithm,Mean Response Time,Response Time Std. Dev.,Mean Weighted Load per Server,Weighted Load per Server Std. Dev.\n")
	diff_reqs.WriteString("Number of Tasks,Algorithm,Mean Response Time,Response Time Std. Dev.,Mean Weighted Load per Server,Weighted Load per Server Std. Dev.\n")

	defer same_reqs.Close()
	defer diff_reqs.Close()
	println("***********************************************************")
	var i int64
	for i = 1000; i < 4000; i += 250 {
		for _, alg := range [2]string{"robin", "state"} {

			// measure metrics for different cores, same load on each request

			avgResponse, stdResponse, avgLoadDist, stdLoadDist := diff_cores_same_requests(alg, int(i))
			// taskData = append(taskData, int(i))
			// algData = append(algData, alg)
			// avgResponseData = append(avgResponseData, avgResponse)
			// stdResponseData = append(stdResponseData, stdResponse)
			// avgLoadData = append(avgLoadData, avgLoadDist)
			// stdLoadData = append(stdLoadData, stdLoadDist)

			same_reqs.WriteString(strconv.FormatInt(i, 10) + ",")
			same_reqs.WriteString(alg + ",")
			same_reqs.WriteString(strconv.FormatFloat(avgResponse, 'f', 3, 64) + ",")
			same_reqs.WriteString(strconv.FormatFloat(stdResponse, 'f', 3, 64) + ",")
			same_reqs.WriteString(strconv.FormatFloat(avgLoadDist, 'f', 3, 64) + ",")
			same_reqs.WriteString(strconv.FormatFloat(stdLoadDist, 'f', 3, 64) + "\n")
			same_reqs.Sync()

			// measure metrics for different cores, randomized load on each request

			avgResponse, stdResponse, avgLoadDist, stdLoadDist = diff_cores_diff_requests(alg, int(i))
			diff_reqs.WriteString(strconv.FormatInt(i, 10) + ",")
			diff_reqs.WriteString(alg + ",")
			diff_reqs.WriteString(strconv.FormatFloat(avgResponse, 'f', 3, 64) + ",")
			diff_reqs.WriteString(strconv.FormatFloat(stdResponse, 'f', 3, 64) + ",")
			diff_reqs.WriteString(strconv.FormatFloat(avgLoadDist, 'f', 3, 64) + ",")
			diff_reqs.WriteString(strconv.FormatFloat(stdLoadDist, 'f', 3, 64) + "\n")
			diff_reqs.Sync()

		}
	}
	println("***********************************************************")

	same_reqs.Sync()
	diff_reqs.Sync()

	//scratch()
}

func main() {

	// println("***********************************************************")
	// const_cores_const_requests("robin")
	// println("***********************************************************")
	// const_cores_const_requests("state")
	// println("***********************************************************")

	// println("***********************************************************")
	// linear_cores_const_requests("robin")
	// println("***********************************************************")
	// linear_cores_const_requests("state")
	// println("***********************************************************")

	// println("***********************************************************")
	// const_cores_random_requests("robin")
	// println("***********************************************************")
	// const_cores_random_requests("state")
	// println("***********************************************************")

	// println("***********************************************************")
	// linear_cores_random_requests("robin")
	// println("***********************************************************")
	// linear_cores_random_requests("state")
	// println("***********************************************************")

	// println("***********************************************************")
	// const_cores_linear_requests("robin")
	// println("***********************************************************")
	// const_cores_linear_requests("state")
	// println("***********************************************************")

	// println("***********************************************************")
	// linear_cores_reverse_requests("robin")
	// println("***********************************************************")
	// linear_cores_reverse_requests("state")
	// println("***********************************************************")
}
