package balancer

import (
	"fmt"
)

var verbose bool = true

type Balancer struct {
	task_queue  queue
	acks        []map[*Request]bool
	servers     []*Server
	next_server int
}

func NewBalancer(servers []*Server) *Balancer {
	balancer := new(Balancer)
	balancer.task_queue = make(queue, 0)
	balancer.acks = make([]map[*Request]bool, 0)
	balancer.servers = servers
	balancer.next_server = 0
	if verbose {
		fmt.Printf("Balancer created: %+v\n", *balancer)
	} // formatting dereferences and prints fields
	return balancer
}

func (balancer *Balancer) Assign_request(request *Request) {
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
func (balancer *Balancer) Assign_request_round_robin(request *Request) {
	servers := balancer.servers

	//find next live server
	for !servers[balancer.next_server].online {
		balancer.next_server++
		println(balancer.next_server)
		if balancer.next_server >= len(servers) {
			balancer.next_server = 0
		}
	}
	//give it the request and increment next server
	servers[balancer.next_server].Add_request(request)
	balancer.next_server++
	if balancer.next_server >= len(servers) {
		balancer.next_server = 0
	}
}

func (balancer *Balancer) Handle_death(*Server) {

}

func (balancer *Balancer) Handle_wakeup(*Server) {

}
