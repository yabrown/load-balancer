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
func main() {
	var wg sync.WaitGroup
	alg := os.Args[1]
	numberOfTasks, _ := strconv.Atoi(os.Args[2])
	taskLoadRange, _ := strconv.Atoi(os.Args[3])
	//timeToSleep, _ := time.ParseDuration(os.Args[2] + "s")
	wg.Add(1)
	// create 5 servers
	s1 := b.NewServer(0, 1)
	s2 := b.NewServer(1, 4)
	s3 := b.NewServer(2, 1)
	s4 := b.NewServer(3, 8)
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
	fmt.Printf("Average Response Time: %.3fms\n", b.MeasureAverageResponseTime(balancer))
	fmt.Printf("Average Task Load: %.3fms\n", b.MeasureAverageTaskLoad(balancer))
	fmt.Printf("Average Task Time: %.3fms\n", b.MeasureAverageTaskTime(balancer))
	b.MeasureLoadDistribution(balancer)

	wg.Wait()
}
