package balancer

import (
	"math/rand"
	"time"
)

type Request struct {
	load int
	id   int
}

type Server struct {
	q             queue
	cores         int
	online        bool
	loadOnQueue   int
	timeCompleted float32
	queueLength   int
}

func newRequest(id int) (request *Request) {
	request = new(Request)
	request.load = rand.Intn(100) + 1 // +1 since the random generation starts at 0
	request.id = id
	return request
}

func newServer() (server *Server) {
	server = new(Server)
	server.q = newQueue()
	server.cores = rand.Intn(4) + 1 // +1 since the random generation starts at 0
	server.online = false
	server.loadOnQueue = 0
	server.timeCompleted = 0
	return server
}

func (server *Server) add_request(request *Request) {
	if server.online {
		server.q.push(request)
		server.loadOnQueue += request.load
		server.queueLength += 1
	}
}

// do next task on queue, which just involves waiting the amount of time divided by cores
// return when the task completes
func (server *Server) do_request() {
	if server.online {
		nextReq, _ := server.q.pop()

		taskTime := float32(nextReq.load / server.cores)
		time.Sleep(time.Duration(taskTime) * time.Millisecond)

		server.timeCompleted += taskTime
		server.queueLength -= 1
		server.loadOnQueue -= nextReq.load
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

func (server *Server) die() {
	// need to send
	server.online = false
}

func (server *Server) wake_up() {
	server.online = true
}
