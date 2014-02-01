package master

import (
	"encoding/gob"
	"errors"
	"net/rpc"
)

// =================================== STRUCT =================================

type rpcInterface struct {
	m *master
}

type RPCProgress struct {
	JobId    uint64
	Frame    int
	Progress byte
}

type RPCJob struct {
	Id       uint64
	Start    int
	Renderer string
	File     []byte
}

type RPCRequestFrame struct {
	JobId uint64
	Name  string
}

type RPCUploadFrame struct {
	JobId uint64
	Frame int
	Data  []byte
}

// =================================== INIT ===================================

func init() {
	gob.Register(&RPCProgress{})
	gob.Register(&RPCJob{})
	gob.Register(&RPCRequestFrame{})
	gob.Register(&RPCUploadFrame{})
}

// =================================== CONSTRUCTOR ============================

func NewRPCInterface(m *master) *rpc.Server {
	i := &rpcInterface{
		m: m,
	}

	r := rpc.NewServer()
	r.RegisterName("master", i)

	return r
}

func CreateFrameRequest(name string, job *RPCJob) *RPCRequestFrame {
	return &RPCRequestFrame{
		JobId: job.Id,
		Name:  name,
	}
}

func CreateRPCJob(job *job) *RPCJob {
	return &RPCJob{
		Id:       job.Id(),
		Start:    job.Start(),
		Renderer: job.Renderer(),
		File:     job.File(),
	}
}

// =================================== FUNCTIONS ==============================

func (r *rpcInterface) RegisterSlave(name string, answer *bool) error {
	*answer = r.m.AddSlave(name)
	return nil
}

func (r *rpcInterface) UnregisterSlave(name string, answer *bool) error {
	r.m.DelSlave(name)
	*answer = true
	return nil
}

func (r *rpcInterface) RequestJob(name string, answer *RPCJob) error {
	job := r.m.NextJob()
	if job == nil {
		return errors.New("No open jobs left.")
	}

	answer = CreateRPCJob(job)
	return nil
}

func (r *rpcInterface) RequestFrame(request RPCRequestFrame, answer *int) error {
	if int(request.JobId) >= r.m.JobCount() {
		return errors.New("Invalid job id.")
	}

	frame := r.m.Job(int(request.JobId)).NextFrame(request.Name)
	if frame == nil {
		return errors.New("No more frames in job.")
	}

	*answer = frame.Frame()
	return nil
}

func (r *rpcInterface) Progress(request RPCProgress, answer *string) error {
	if int(request.JobId) >= r.m.JobCount() {
		return errors.New("Invalid job id.")
	}

	r.m.Job(int(request.JobId)).Frame(request.Frame).SetProgress(request.Progress)
	*answer = "ok"
	return nil
}

func (r *rpcInterface) UploadFrame(request RPCUploadFrame, answer *bool) error {
	if int(request.JobId) >= r.m.JobCount() {
		return errors.New("Invalid job id.")
	}

	r.m.Job(int(request.JobId)).Frame(request.Frame).SetData(request.Data)
	r.m.Job(int(request.JobId)).Frame(request.Frame).SetCompleted(true)

	*answer = true
	return nil
}
