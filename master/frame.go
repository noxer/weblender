package master

import ()

// =================================== STRUCT =================================

type frame struct {
	job       *job
	slaveName string
	frame     int
	data      []byte
	progress  byte
	completed bool
}

// =================================== CONSTRUCTOR ============================

func NewFrame(job *job, fr int) *frame {
	return &frame{
		job:       job,
		slaveName: "",
		frame:     fr,
		data:      nil,
		progress:  0,
		completed: false,
	}
}

// =================================== GETTER =================================

func (f *frame) Job() *job {
	return f.job
}

func (f *frame) SlaveName() string {
	return f.slaveName
}

func (f *frame) Frame() int {
	return f.frame
}

func (f *frame) AlignedFrame() int {
	return f.frame + f.Job().Start()
}

func (f *frame) Data() []byte {
	return f.data
}

func (f *frame) Progress() byte {
	return f.progress
}

func (f *frame) Completed() bool {
	return f.completed
}

func (f *frame) File() []byte {
	return f.Job().File()
}

// =================================== SETTER =================================

func (f *frame) SetSlaveName(slaveName string) {
	f.slaveName = slaveName
}

func (f *frame) SetData(data []byte) {
	f.data = data
}

func (f *frame) SetProgress(progress byte) {
	f.progress = progress
}

func (f *frame) SetCompleted(completed bool) {
	f.completed = completed
	if completed {
		f.SetProgress(100)
	}
}

// =================================== FUNCTIONS ==============================

func (f *frame) Reset() {
	f.SetSlaveName("")
	f.SetData(nil)
	f.SetProgress(0)
	f.SetCompleted(false)
}
