package master

import (
	"net/http"
	"sync"
	"time"
)

// =================================== STRUCT =================================

type master struct {
	slaves map[string]*slave
	jobs   []*job

	sMutex, jMutex sync.RWMutex
}

// =================================== CONSTRUCTOR ============================

func NewMaster() *master {
	return &master{
		slaves: make(map[string]*slave),
		jobs:   make([]*job, 0),
	}
}

// =================================== GETTER =================================

func (m *master) SlaveCount() int {
	m.sMutex.RLock()
	defer m.sMutex.RUnlock()

	return len(m.slaves)
}

func (m *master) Slave(name string) *slave {
	m.sMutex.RLock()
	defer m.sMutex.RUnlock()

	return m.slaves[name]
}

func (m *master) JobCount() int {
	m.jMutex.RLock()
	defer m.jMutex.RUnlock()

	return len(m.jobs)
}

func (m *master) Job(index int) *job {
	m.jMutex.RLock()
	defer m.jMutex.RUnlock()

	return m.jobs[index]
}

// =================================== SETTER =================================

func (m *master) AddSlave(name string) bool {
	m.sMutex.Lock()
	defer m.sMutex.Unlock()

	if _, ok := m.slaves[name]; ok {
		return false
	}

	m.slaves[name] = NewSlave(name)
	return true
}

func (m *master) DelSlave(name string) {
	m.sMutex.Lock()
	defer m.sMutex.Unlock()
	m.jMutex.Lock()
	defer m.jMutex.Unlock()

	delete(m.slaves, name)

	for _, j := range m.jobs {
		j.RemoveSlave(name)
	}
}

func (m *master) AddJob(start, end int, renderer string, file []byte) uint64 {
	m.jMutex.Lock()
	defer m.jMutex.Unlock()

	jobId := uint64(len(m.jobs))
	m.jobs = append(m.jobs, NewJob(jobId, start, end, renderer, file))

	return jobId
}

// =================================== FUNCTIONS ==============================

func (m *master) RemoveOldSlaves() {
	m.sMutex.Lock()
	defer m.sMutex.Unlock()
	m.jMutex.Lock()
	defer m.jMutex.Unlock()

	now := time.Now()
	removed := make([]string, 0)

	// Remove the slaves from the slaves map
	for n, s := range m.slaves {
		if now.Sub(s.LastSeen()) > 5*time.Minute {
			delete(m.slaves, n)
			removed = append(removed, n)
		}
	}

	// Remove the slaves from all jobs they where working on
	for _, r := range removed {
		for _, j := range m.jobs {
			j.RemoveSlave(r)
		}
	}
}

func (m *master) NextJob() *job {
	m.jMutex.RLock()
	defer m.jMutex.RUnlock()

	for _, j := range m.jobs {
		if j.Open() {
			return j
		}
	}

	return nil
}

func (m *master) Handler() http.Handler {
	return NewRPCInterface(m)
}
