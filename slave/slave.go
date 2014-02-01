package slave

import (
	"net/rpc"
	"weblender/blender"
	"weblender/master"
)

type slave struct {
	client *rpc.Client
	name   string
}

func New(host, path, name string) *slave {
	// Connect to master
	cl, err := rpc.DialHTTPPath("tcp", host, path)
	if err != nil {
		return nil
	}

	// Register at master
	var result bool
	cl.Call("master.RegisterSlave", name, &result)
	if !result {
		return nil
	}

	return &slave{
		client: cl,
		name:   name,
	}
}

func (s *slave) RequestJob() *master.RPCJob {
	var job master.RPCJob
	err := s.client.Call("master.RequestJob", s.name, &job)
	if err != nil {
		return nil
	}
	return &job
}

func (s *slave) RequestFrame(job *master.RPCJob) int {
	var frame int
	err := s.client.Call("master.RequestFrame", master.CreateFrameRequest(s.name, job), &frame)
	if err != nil {
		return -1
	}
	return frame
}

func (s *slave) UploadFrame(job *master.RPCJob, frame int, data []byte) error {
	var result bool
	err := 
}
