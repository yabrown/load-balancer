package balancer

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestQueueBasic(t *testing.T) {
	req := NewRequest(1, 1)
	req2 := NewRequest(2, 1)
	/////// Test basic push, push, pop, pop ///////
	q := make(queue, 0)
	q = q.push(req)
	q = q.push(req2)
	popped, q := q.pop()
	if popped.id != 1 {
		t.Errorf("Failed to pop 'first', popped '%v' instead", popped.id)
	}
	popped, q = q.pop()
	if popped.id != 2 {
		t.Errorf("Failed to pop 'second', popped '%v' instead", popped.id)
	}

	//////// Test basic push pop //////
	q = make(queue, 0)
	q = q.push(NewRequest(3, 1))
	popped, q = q.pop()
	if popped.id != 3 {
		t.Errorf("Expected 'hi' to be popped, received '%v' instead", popped)
	}

	////// Popping empty queue ///////
	q = make(queue, 0)
	popped, q = q.pop()
	if popped != nil {
		t.Errorf("Expected 'failed', received '%v'", popped)
	}
	popped, q = q.pop()
	if popped != nil {
		t.Errorf("Expected 'failed', received '%v'", popped)
	}
	q = q.push(NewRequest(4, 1))
	popped, q = q.pop()
	popped, q = q.pop()
	if popped != nil {
		t.Errorf("Expected 'failed', received '%v'", popped)
	}

}

// test 5 servers, each with one core, on 15 tasks that have a load of 1
func TestSameCoreSameTask(t *testing.T) {
	algs := [2]string{"robin", "state"}
	for _, alg := range algs {
		servers := make([]*Server, 5)
		for i := 0; i < 5; i++ {
			servers[i] = NewServer(i, 1)
		}
		balancer := NewBalancer(servers)

		requests := []*Request{}

		for i := 0; i < 15; i++ {

			requests = append(requests, NewRequest(i, 1))
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

		for i := range servers {
			go servers[i].Work(balancer)
		}

		time.Sleep(time.Duration(3100) * time.Millisecond)

		avgResponseTime := MeasureAverageResponseTime(balancer)

		if alg == "robin" {
			if math.Abs(float64(avgResponseTime))-2 > .02 {
				t.Errorf("Expected average response time of 2, received average response time of %v", avgResponseTime)
			} else {
				fmt.Printf("Success for algorithm \"%s\". Average response time is %f seconds.\n", alg, avgResponseTime)
			}
		} else if alg == "state" {
			if math.Abs(float64(avgResponseTime))-float64(2) > .02 {
				t.Errorf("Expected average response time of 2, received average response time of %v", avgResponseTime)
			} else {
				fmt.Printf("Success for algorithm \"%s\". Average response time is %f seconds.\n", alg, avgResponseTime)
			}
		}
	}

}

// test 3 servers, two with 1 core and one with 2 cores, on 15 tasks that have a load of 2
func TestDiffCoreSameTask(t *testing.T) {
	algs := [2]string{"robin", "state"}
	for _, alg := range algs {
		servers := make([]*Server, 3)
		for i := 0; i < 2; i++ {
			servers[i] = NewServer(i, 1)
		}

		servers[2] = NewServer(2, 2)
		balancer := NewBalancer(servers)

		requests := []*Request{}

		for i := 0; i < 15; i++ {

			requests = append(requests, NewRequest(i, 2))
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

		for i := range servers {
			go servers[i].Work(balancer)
		}

		time.Sleep(time.Duration(15000) * time.Millisecond)

		avgResponseTime := MeasureAverageResponseTime(balancer)

		if alg == "robin" {
			correctAvg := 5.0
			if math.Abs(avgResponseTime-correctAvg) > .1 {
				t.Errorf("Failure for algorithm \"%s\". Expected average response time of %f, received average response time of %v", alg, correctAvg, avgResponseTime)
			} else {
				fmt.Printf("Success for algorithm \"%s\". Average response time is %f seconds.\n", alg, avgResponseTime)
			}
		} else if alg == "state" {
			correctAvg := 4.66
			if math.Abs(avgResponseTime-correctAvg) > .1 {
				t.Errorf("Failure for algorithm \"%s\". Expected average response time of %f, received average response time of %v", alg, correctAvg, avgResponseTime)
			} else {
				fmt.Printf("Success for algorithm \"%s\". Average response time is %f seconds.\n", alg, avgResponseTime)
			}
		}
	}

}

// test 3 servers, each with 1 core, on 15 tasks that have a load of 1 or 2
// the task listing is an alternating list of 1, 2, 1, 2....
func TestSameCoreDiffTask(t *testing.T) {
	algs := [2]string{"robin", "state"}
	for _, alg := range algs {
		servers := make([]*Server, 3)
		for i := 0; i < 3; i++ {
			servers[i] = NewServer(i, 1)
		}

		balancer := NewBalancer(servers)

		requests := []*Request{}

		for i := 0; i < 15; i++ {
			// 8 tasks of size 1, 7 tasks of size 2, alternating starting with 1
			taskSize := (i % 2) + 1
			requests = append(requests, NewRequest(i, taskSize))
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

		for i := range servers {
			go servers[i].Work(balancer)
		}

		time.Sleep(time.Duration(10000) * time.Millisecond)

		avgResponseTime := MeasureAverageResponseTime(balancer)

		if alg == "robin" {
			correctAvg := 4.4
			if math.Abs(avgResponseTime-correctAvg) > .1 {
				t.Errorf("Failure for algorithm \"%s\". Expected average response time of %f, received average response time of %v", alg, correctAvg, avgResponseTime)
			} else {
				fmt.Printf("Success for algorithm \"%s\". Average response time is %f seconds.\n", alg, avgResponseTime)
			}
		} else if alg == "state" {
			correctAvg := 4.3
			if math.Abs(avgResponseTime-correctAvg) > .1 {
				t.Errorf("Failure for algorithm \"%s\". Expected average response time of %f, received average response time of %v", alg, correctAvg, avgResponseTime)
			} else {
				fmt.Printf("Success for algorithm \"%s\". Average response time is %f seconds.\n", alg, avgResponseTime)
			}
		}
	}

}

// test 3 servers, two with 1 core and one with 2 cores, on 15 tasks that have a load of 1 or 2
// the task listing is an alternating list of 1, 2, 1, 2....
func TestDiffCoreDiffTask(t *testing.T) {
	algs := [2]string{"robin", "state"}
	for _, alg := range algs {
		servers := make([]*Server, 3)
		for i := 0; i < 2; i++ {
			servers[i] = NewServer(i, 1)
		}

		servers[2] = NewServer(2, 2)

		balancer := NewBalancer(servers)

		requests := []*Request{}

		for i := 0; i < 15; i++ {
			// 8 tasks of size 1, 7 tasks of size 2, alternating starting with 1
			taskSize := (i % 2) + 1
			requests = append(requests, NewRequest(i, taskSize))
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

		for i := range servers {
			go servers[i].Work(balancer)
		}

		time.Sleep(time.Duration(10000) * time.Millisecond)

		avgResponseTime := MeasureAverageResponseTime(balancer)

		if alg == "robin" {
			correctAvg := 3.7
			if math.Abs(avgResponseTime-correctAvg) > .1 {
				t.Errorf("Failure for algorithm \"%s\". Expected average response time of %f, received average response time of %v", alg, correctAvg, avgResponseTime)
			} else {
				fmt.Printf("Success for algorithm \"%s\". Average response time is %f seconds.\n", alg, avgResponseTime)
			}
		} else if alg == "state" {
			correctAvg := 3.4
			if math.Abs(avgResponseTime-correctAvg) > .1 {
				t.Errorf("Failure for algorithm \"%s\". Expected average response time of %f, received average response time of %v", alg, correctAvg, avgResponseTime)
			} else {
				fmt.Printf("Success for algorithm \"%s\". Average response time is %f seconds.\n", alg, avgResponseTime)
			}
		}
	}

}
