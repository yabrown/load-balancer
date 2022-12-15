package main

import (
	b "cos316-load-balancer/balancer"
	"fmt"
	"math/rand"
)

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
	fmt.Printf("Average Response Time: %.3fms\n", b.MeasureAverageResponseTime(balancer))
	fmt.Printf("Average Task Load: %.3fms\n", b.MeasureAverageTaskLoad(balancer))
	fmt.Printf("Average Task Time: %.3fms\n", b.MeasureAverageTaskTime(balancer))
	b.MeasureLoadDistribution(balancer)

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
	fmt.Printf("Average Response Time: %.3fms\n", b.MeasureAverageResponseTime(balancer))
	fmt.Printf("Average Task Load: %.3fms\n", b.MeasureAverageTaskLoad(balancer))
	fmt.Printf("Average Task Time: %.3fms\n", b.MeasureAverageTaskTime(balancer))
	b.MeasureLoadDistribution(balancer)

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
	fmt.Printf("Average Response Time: %.3fms\n", b.MeasureAverageResponseTime(balancer))
	fmt.Printf("Average Task Load: %.3fms\n", b.MeasureAverageTaskLoad(balancer))
	fmt.Printf("Average Task Time: %.3fms\n", b.MeasureAverageTaskTime(balancer))
	b.MeasureLoadDistribution(balancer)

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
	fmt.Printf("Average Response Time: %.3fms\n", b.MeasureAverageResponseTime(balancer))
	fmt.Printf("Average Task Load: %.3fms\n", b.MeasureAverageTaskLoad(balancer))
	fmt.Printf("Average Task Time: %.3fms\n", b.MeasureAverageTaskTime(balancer))
	b.MeasureLoadDistribution(balancer)

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
	fmt.Printf("Average Response Time: %.3fms\n", b.MeasureAverageResponseTime(balancer))
	fmt.Printf("Average Task Load: %.3fms\n", b.MeasureAverageTaskLoad(balancer))
	fmt.Printf("Average Task Time: %.3fms\n", b.MeasureAverageTaskTime(balancer))
	b.MeasureLoadDistribution(balancer)

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
	fmt.Printf("Average Response Time: %.3fms\n", b.MeasureAverageResponseTime(balancer))
	fmt.Printf("Average Task Load: %.3fms\n", b.MeasureAverageTaskLoad(balancer))
	fmt.Printf("Average Task Time: %.3fms\n", b.MeasureAverageTaskTime(balancer))
	b.MeasureLoadDistribution(balancer)

}
