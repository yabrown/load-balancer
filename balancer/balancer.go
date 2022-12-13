package balancer

import (
	"fmt"
	"time"
)

type Balancer struct {
	task_queue  queue
	acks        map[*Request]bool
	servers     []*Server
	next_server int
}

func newBalancer(servers []*Server) *Balancer {
	balancer := new(Balancer)
	balancer.task_queue = make(queue, 0)
	balancer.acks = make(map[*Request]bool)
	balancer.servers = servers
	balancer.next_server = 0
	return balancer
}

func (balancer *Balancer) assign_request(request Request) {
	servers := balancer.servers // so that don't have to recompute each time
	//champ := servers[0]
	// find least loaded server
	for i := range servers {
		if !(servers[i].online) { // if server is dead, ignore it
			continue
		}
	}

}

//just assigns to next live server
func (balancer *Balancer) assign_request_round_robin(request *Request) {
	servers := balancer.servers

	//find next live server
	for !servers[balancer.next_server].online {
		balancer.next_server++
		if balancer.next_server >= len(servers) { // if all servers are offline, wait a second, and check all again
			time.Sleep(1 * time.Second)
		}
	}

	servers[balancer.next_server].add_request(request) //give it the request
	balancer.next_server++                             //increment next server
}
