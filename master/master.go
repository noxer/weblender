package master

import (
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

type instance struct {
	slaveMutex sync.RWMutex
	slaves     map[string]*slave
	jobMutex   sync.RWMutex
	jobs       []*job
	jobIds     chan uint64
}

type slave struct {
	Name     string
	LastSeen time.Time
	JobId    uint64
}

func (i *instance) Start() http.Handler {
	// Start job id generator
	i.jobIds = make(chan int, 1)
	go i.idGen()

	// Create rpc server
	srv := rpc.NewServer()
	srv.RegisterName("master", &rpcInterface{master: i})

	return srv
}

func (i *instance) idGen() {
	nextId := uint64(1)

	for {
		i.jobIds <- nextId
		nextId++
	}
}
