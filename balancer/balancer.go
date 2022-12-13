package balancer

import (
	"fmt"
	"math"
	"time"
)

var verbose bool = true

type stats struct {
	start_time time.Time     // set as soon as it comes in (both assign_request functions)
	handled    bool          //set to true once acknowledged here in balancer (don't go to server)
	duration   time.Duration //calculate once acknowledged here in balancer
}

type Balancer struct {
	task_queue    queue
	acks          []map[int]*Request
	servers       []*Server
	next_server   int
	request_stats map[*Request]*stats
}

func NewBalancer(servers []*Server) *Balancer {
	balancer := new(Balancer)
	balancer.task_queue = make(queue, 0)
	balancer.acks = make([]map[int]*Request, len(servers))
	for i := range servers {
		balancer.acks[i] = make(map[int]*Request)
	}
	balancer.request_stats = make(map[*Request]*stats, 0)

	balancer.servers = servers
	balancer.next_server = 0
	if verbose {
		fmt.Printf("Balancer created: %+v\n", *balancer)
	} // formatting dereferences and prints fields
	return balancer
}

func (balancer *Balancer) Assign_request(request *Request) {
	//first, create stats object
	balancer.request_stats[request] = new(stats)

	balancer.request_stats[request].handled = false
	balancer.request_stats[request].start_time = time.Now()

	servers := balancer.servers
	champ_server, champ_load := 0, math.Inf(1)
	//find next weighted least-loaded live server
	for i := range servers {
		if !servers[i].online {
			continue
		}
		curr_load := float64(servers[i].loadOnQueue / servers[i].cores)
		if curr_load < champ_load {
			champ_load = curr_load
			champ_server = i
		}
	}

	//give it the request, store in ack, and increment next server
	servers[champ_server].Add_request(request)
	balancer.acks[champ_server][request.id] = request
}

//just assigns to next live server, stores in acks
func (balancer *Balancer) Assign_request_round_robin(request *Request) {
	servers := balancer.servers

	balancer.request_stats[request] = new(stats)

	balancer.request_stats[request].handled = false
	balancer.request_stats[request].start_time = time.Now()

	//find next live server
	for !servers[balancer.next_server].online {
		balancer.next_server++
		println(balancer.next_server)
		if balancer.next_server >= len(servers) {
			balancer.next_server = 0
		}
	}
	//give it the request, store in ack, and increment next server
	servers[balancer.next_server].Add_request(request)
	balancer.acks[balancer.next_server][request.id] = request
	balancer.next_server++
	if balancer.next_server >= len(servers) {
		balancer.next_server = 0
	}
}

func (balancer *Balancer) Handle_death(*Server) {

}

func (balancer *Balancer) Ack_request(server_id int, request *Request) {
	balancer.request_stats[request].handled = true
	delete(balancer.acks[server_id], request.id)
	start_time := balancer.request_stats[request].start_time
	balancer.request_stats[request].duration = time.Since(start_time)
}

func (balancer *Balancer) Handle_wakeup(*Server) {

}
