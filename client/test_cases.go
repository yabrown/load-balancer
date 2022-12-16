package main

import (
	b "cos316-load-balancer/balancer"
	"fmt"
	"math/rand"
)

func analyze(balancer *b.Balancer) {
	RTAverage, RTStdDev, RTMedian := b.MeasureResponseTime(balancer)
	TTAverage, _, _ := b.MeasureAverageTaskTime(balancer)
	fmt.Printf("Average Task Load: %.3fms\n", b.MeasureAverageTaskLoad(balancer))
	fmt.Printf("Average Task Time: %.3f\n", TTAverage)
	fmt.Printf("Response Time: Average = %.3f, Variation = %.3f, Median = %.3fms\n", RTAverage, RTStdDev/RTAverage, RTMedian)

	b.MeasureLoadDistribution(balancer)
}

// Sophie, write whatever you want here
func data(balancer *b.Balancer) (float64, float64, float64, float64) {
	averageResponseTime, stdDevResponseTime, _ := b.MeasureResponseTime(balancer)
	averageLoadDist, stdDevLoadDist, _ := b.MeasureLoadDistribution(balancer)
	return averageResponseTime, stdDevResponseTime, averageLoadDist, stdDevLoadDist
}

// servers: 5, each with 1 core
// tasks: 1000, each random with load 1-100
func const_cores_random_requests(alg string) {
	//timeToSleep, _ := time.ParseDuration(os.Args[2] + "s")
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

	for i := 0; i < 1000; i++ {
		random := rand.Intn(5) + 1
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
	analyze(balancer)

}

// servers: 5, each with 1 core
// tasks: 1000, with worst distribution-- very uneven
func const_cores_linear_requests(alg string) {
	//timeToSleep, _ := time.ParseDuration(os.Args[2] + "s")
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

	for i := 0; i < 1000; i++ {
		load := ((i % 5) + 1)
		requests = append(requests, b.NewRequest(i, load))
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
	analyze(balancer)

}

// servers: 5, each with 1 core
// tasks: 1000, with best distribution-- even
func const_cores_const_requests(alg string) {
	//timeToSleep, _ := time.ParseDuration(os.Args[2] + "s")
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

	for i := 0; i < 1000; i++ {
		load := 3
		requests = append(requests, b.NewRequest(i, load))
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
	analyze(balancer)

}

// servers: 5, each with 1 core
// tasks: 1000, with best distribution-- even
func linear_cores_const_requests(alg string) {
	//timeToSleep, _ := time.ParseDuration(os.Args[2] + "s")
	// create 5 servers
	s1 := b.NewServer(0, 1)
	s2 := b.NewServer(1, 2)
	s3 := b.NewServer(2, 3)
	s4 := b.NewServer(3, 4)
	s5 := b.NewServer(4, 5)
	// create balancer
	servers := []*b.Server{s1, s2, s3, s4, s5}
	end_signal := make(chan bool)
	balancer := b.NewBalancer(servers, end_signal)
	if balancer == nil {
	}
	// create requests
	requests := []*b.Request{}

	for i := 0; i < 1000; i++ {
		load := 3
		requests = append(requests, b.NewRequest(i, load))
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
	print("got here")
	fmt.Println("Completing tasks with algorithm option", alg)
	//fmt.Println("Sleeping for", timeToSleep)
	<-end_signal
	// system_info(servers, balancer)
	analyze(balancer)

}

// servers: 5, each with 1 core
// tasks: 1000, with best distribution-- even
func linear_cores_reverse_requests(alg string) {
	//timeToSleep, _ := time.ParseDuration(os.Args[2] + "s")
	// create 5 servers
	s1 := b.NewServer(0, 1)
	s2 := b.NewServer(1, 2)
	s3 := b.NewServer(2, 3)
	s4 := b.NewServer(3, 4)
	s5 := b.NewServer(4, 5)
	// create balancer
	servers := []*b.Server{s1, s2, s3, s4, s5}
	end_signal := make(chan bool)
	balancer := b.NewBalancer(servers, end_signal)
	if balancer == nil {
	}
	// create requests
	requests := []*b.Request{}

	for i := 0; i < 1000; i++ {
		load := 5 - ((i % 5) + 1)
		requests = append(requests, b.NewRequest(i, load))
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
	print("got here")
	fmt.Println("Completing tasks with algorithm option", alg)
	//fmt.Println("Sleeping for", timeToSleep)
	<-end_signal
	// system_info(servers, balancer)
	analyze(balancer)

}

// servers: 5, each with 1 core
// tasks: 1000, with best distribution-- even
func linear_cores_random_requests(alg string) {
	//timeToSleep, _ := time.ParseDuration(os.Args[2] + "s")
	// create 5 servers
	s1 := b.NewServer(0, 1)
	s2 := b.NewServer(1, 2)
	s3 := b.NewServer(2, 3)
	s4 := b.NewServer(3, 4)
	s5 := b.NewServer(4, 5)
	// create balancer
	servers := []*b.Server{s1, s2, s3, s4, s5}
	end_signal := make(chan bool)
	balancer := b.NewBalancer(servers, end_signal)
	if balancer == nil {
	}
	// create requests
	requests := []*b.Request{}

	for i := 0; i < 1000; i++ {
		random := rand.Intn(5) + 1
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
	analyze(balancer)
}

// Case One:
// servers: 5, each with 1 core
// tasks: 1000, with best distribution-- even
func diff_cores_same_requests(alg string, num_tasks int) (float64, float64, float64, float64) {
	//timeToSleep, _ := time.ParseDuration(os.Args[2] + "s")
	// create 5 servers
	s1 := b.NewServer(0, 1)
	s2 := b.NewServer(1, 1)
	s3 := b.NewServer(2, 3)
	s4 := b.NewServer(3, 5)
	s5 := b.NewServer(4, 10)
	// create balancer
	servers := []*b.Server{s1, s2, s3, s4, s5}
	end_signal := make(chan bool)
	balancer := b.NewBalancer(servers, end_signal)
	if balancer == nil {
	}
	// create requests
	requests := []*b.Request{}

	for i := 0; i < num_tasks; i++ {
		load := 2
		requests = append(requests, b.NewRequest(i, load))
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
	avgResponse, stdResponse, avgLoadDist, stdLoadDist := data(balancer)

	//fmt.Printf("Average response time of %fms with a standard deviation of %fms.\n", avgResponse, stdResponse)
	//fmt.Printf("Average weighted load for each server of %f with a standard deviation of %f.\n", avgLoadDist, stdLoadDist)
	return avgResponse, stdResponse, avgLoadDist, stdLoadDist
}

// Case 2
func diff_cores_diff_requests(alg string, num_tasks int) (float64, float64, float64, float64) {
	//timeToSleep, _ := time.ParseDuration(os.Args[2] + "s")
	// create 5 servers
	s1 := b.NewServer(0, 1)
	s2 := b.NewServer(1, 1)
	s3 := b.NewServer(2, 3)
	s4 := b.NewServer(3, 5)
	s5 := b.NewServer(4, 10)
	// create balancer
	servers := []*b.Server{s1, s2, s3, s4, s5}
	end_signal := make(chan bool)
	balancer := b.NewBalancer(servers, end_signal)
	if balancer == nil {
	}
	// create requests
	requests := []*b.Request{}
	for i := 0; i < num_tasks; i++ {
		random := (rand.Intn(20) + 1)
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
	avgResponse, stdResponse, avgLoadDist, stdLoadDist := data(balancer)
	return avgResponse, stdResponse, avgLoadDist, stdLoadDist

}
