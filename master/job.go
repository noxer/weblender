package master

import ()

// =================================== STRUCT =================================

type job struct {
	id         uint64
	start, end int
	frames     []*frame
	renderer   string
	file       []byte
}

// =================================== CONSTRUCTOR ============================

func NewJob(id uint64, start, end int, renderer string, file []byte) *job {
	// Create job object
	j := &job{
		id:       id,
		start:    start,
		end:      end,
		frames:   make([]*frame, 1+end-start),
		renderer: renderer,
		file:     file,
	}
	// Create frames list
	for i := 0; i <= end-start; i++ {
		j.frames[i] = NewFrame(j, i)
	}

	return j
}

// =================================== GETTER =================================

func (j *job) Id() uint64 {
	return j.id
}

func (j *job) Start() int {
	return j.start
}

func (j *job) End() int {
	return j.end
}

func (j *job) FrameCount() int {
	return len(j.frames)
}

func (j *job) CompleteFrameCount() int {
	count := 0
	for _, f := range j.frames {
		if f.Completed() {
			count++
		}
	}
	return count
}

func (j *job) IncompleteFrameCount() int {
	count := 0
	for _, f := range j.frames {
		if !f.Completed() {
			count++
		}
	}
	return count
}

func (j *job) Frame(index int) *frame {
	return j.frames[index]
}

func (j *job) Renderer() string {
	return j.renderer
}

func (j *job) File() []byte {
	return j.file
}

// =================================== FUNCTIONS ==============================

func (j *job) RemoveSlave(name string) {
	for _, f := range j.frames {
		if !f.Completed() && f.SlaveName() == name {
			f.Reset()
		}
	}
}

func (j *job) Completed() bool {
	return j.IncompleteFrameCount() == 0
}

func (j *job) Open() bool {
	if j.Completed() {
		return false
	}

	for _, f := range j.frames {
		if !f.Completed() && f.SlaveName() == "" {
			return true
		}
	}
	return false
}

func (j *job) NextFrame(slave string) *frame {
	for _, f := range j.frames {
		if !f.Completed() && f.SlaveName() == "" {
			f.SetSlaveName(slave)
			return f
		}
	}
	return nil
}
