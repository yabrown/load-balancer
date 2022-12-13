package main

import (
	b "cos316-load-balancer/balancer"
	"fmt"
	"sync"
	"time"
)

func system_info(servers []*b.Server, balancer *b.Balancer) {
	print("\n")
	for i := range servers {
		fmt.Printf("\tServer status: %+v\n", *servers[i]) // formatting dereferences and prints fields
	}
	fmt.Printf("Balancer status: %+v\n\n", balancer) // formatting dereferences and prints fields
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	// create 5 servers
	s1 := b.NewServer(0, 1)
	s2 := b.NewServer(1, 4)
	s3 := b.NewServer(2, 1)
	s4 := b.NewServer(3, 1)
	s5 := b.NewServer(4, 1)
	// create balancer
	servers := []*b.Server{s1, s2, s3, s4, s5}
	balancer := b.NewBalancer(servers)
	if balancer == nil {
	}
	// create requests
	requests := []*b.Request{}
	for i := 0; i < 10; i++ {
		requests = append(requests, b.NewRequest(i, 1))
	}
	//wake all servers up
	for i := range servers {
		servers[i].Wake_up()
	}
	// assign each request
	for i := range requests {
		// balancer.Assign_request(requests[i])
		balancer.Assign_request_round_robin(requests[i])
	}
	system_info(servers, balancer)
	//make each server active (start handling queued requests)
	for i := range servers {
		go servers[i].Work(balancer)
	}

	time.Sleep(5 * time.Second)
	wg.Done()
	// system_info(servers, balancer)
	fmt.Printf("Average Response Time: %.3fs\n", b.MeasureAverageResponseTime(balancer))
	fmt.Printf("Average Task Time: %.3fs\n", b.MeasureAverageTaskTime(balancer))

	wg.Wait()
}
