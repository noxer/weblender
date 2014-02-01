package master

import ()

type job struct {
	id         uint64
	file       []byte
	start, end int
	frames     [][]byte
}

func (i *instance) NewJob(file []byte, start, end int) *job {
	&job{
		id:     <-i.jobIds,
		file:   file,
		start:  start,
		end:    end,
		frames: make([][]byte, (end-start)+1),
	}
}

func (j *job) Id() uint64 {
	return j.id
}

func (j *job) File() []byte {
	return j.file
}

func (j *job) Start() int {
	return j.start
}

func (j *job) End() int {
	return j.End()
}

func (j *job) Frame(frame int) []byte {
	return j.frames(frame)
}

func (j *job) AlignedFrame(frame int) []byte {
	return j.Frame(frame + j.Start())
}

func (j *job) SetFrame(frame int, data []byte) {
	j.frames[frame] = data
}

func (j *job) SetAlignedFrame(frame int, data []byte) {
	j.SetFrame(frame+j.Start(), data)
}

func (j *job) Complete() bool {
	for _, f := range j.frames {
		if f == nil {
			return false
		}
	}
	return true
}

func (j *job) NextMissing() int {
	for i, f := range j.frames {
		if f == nil {
			return i
		}
	}
	return -1
}

func (j *job) NextMissingAligned() int {
	if m := j.NextMissing(); m < 0 {
		return m
	} else {
		return m + j.Start()
	}
}
