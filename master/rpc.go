package master

import (
	gjob "weblender/job"
)

type rpcInterface struct {
	master *instance
}

func (r *rpcInterface) RegisterSlave(name string, answer *string) error {
	// Is the slave already present?
	if _, ok := r.master.slaves[name]; ok {
		answer = "name taken"
	} else {
		answer = "ok"
	}

	return nil
}

func (r *rpcInterface) UnregisterSlave(name string, answer *string) error {
	delete(r.master.slaves, name)
	answer = "ok"
	return nil
}

func (r *rpcInterface) RequestJob(name string, answer *gjob.Instance) {
	if sl, ok := r.master.slaves[name]; !ok {
		answer = nil
	}

}
