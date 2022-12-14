package balancer

import (
	"fmt"
	"time"
)

type Request struct {
	load     int
	id       int
	duration time.Duration
}

//defined load
func NewRequest(id int, load int) (request *Request) {
	request = new(Request)
	request.id = id
	request.load = load
	return request
}

type Server struct {
	id            int
	q             queue
	cores         int
	online        bool
	loadOnQueue   int
	timeCompleted float32
	queueLength   int
}

func NewServer(id int, cores int) *Server {
	server := new(Server)
	server.id = id
	server.q = newQueue()
	server.cores = cores
	server.online = false
	server.loadOnQueue = 0
	server.timeCompleted = 0
	if verbose {
		fmt.Printf("Server created: %+v\n", *server)
	} // formatting dereferences and prints fields
	return server
}

func (server *Server) Add_request(request *Request) {
	if verbose {
		fmt.Printf("Adding request %d to server %d\n", request.id, server.id)
	}
	if server.online {
		server.q = server.q.push(request)
		server.loadOnQueue += request.load
		server.queueLength += 1
	}
}

// do next task on queue, which just involves waiting the amount of time divided by cores
// return when the task completes
func (server *Server) Handle_request(balancer *Balancer) {
	if server.online {
		if verbose {
			fmt.Printf("Server %d handling next request.\n", server.id)
		}
		if verbose {
			fmt.Printf("\tStart handling: %+v\n", *server)
		} // formatting dereferences and prints fields
		var nextReq *Request
		nextReq, server.q = server.q.pop()

		taskTime := float32(nextReq.load) / float32(server.cores)
		fmt.Println(taskTime, "=", time.Duration(taskTime)*time.Second)
		if taskTime < 1 {
			taskTimeInMilli := taskTime * 1000
			time.Sleep(time.Duration(taskTimeInMilli) * time.Millisecond)
		} else {
			time.Sleep(time.Duration(taskTime) * time.Second)
		}

		server.timeCompleted += taskTime
		server.queueLength -= 1
		server.loadOnQueue -= nextReq.load
		balancer.Ack_request(server.id, nextReq)
		if verbose {
			fmt.Printf("\tFinished handling: %+v\n", *server)
		} // formatting dereferences and prints fields

	}

}

// return info about state-- whether online, total time (total load/cores), length of queue
func (server *Server) State() (bool, int) {
	if server.online {
		return server.online, server.loadOnQueue
	} else {
		return server.online, 0
	}
}

func (server *Server) Die() {
	// need to send
	server.online = false
}

func (server *Server) Wake_up() {
	server.online = true

}

// call as goroutine. keeps handling whatevers in front of the queue, returns when the server dies
func (server *Server) Work(balancer *Balancer) {
	if verbose {
		fmt.Printf("Server %d started listening.\n", server.id)
	}
	for {
		if len(server.q) != 0 {
			server.Handle_request(balancer)
		}
		if !server.online {
			if verbose {
				fmt.Printf("Server %d stopped listening.\n", server.id)
			}
			return
		}
	}
}
