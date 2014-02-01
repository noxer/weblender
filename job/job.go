package job

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"weblender/blender"
)

func init() {
	gob.Register(&instance{})
	gob.Register(&result{})
}

type Instance interface {
	Id() uint64
	Renderer() string
	Frame() (int, int)
	File() []byte
}

type instance struct {
	id           uint64 // Job ID
	renderer     string // Renderer
	start, frame int    // Start frame & frame offset
	file         []byte // Blend file
}

type result struct {
	id           uint64
	start, frame int
	file         []byte
	time         time.Duration
}

func New(id uint64, renderer string, start, frame int, file []byte) *instance {
	return &instance{
		id:       id,
		renderer: renderer,
		start:    start,
		frame:    frame,
		file:     file,
	}
}

func CreateResult(renderer blender.Instance) *result {
	start, frame := renderer.Frame()

	file, err := os.Open(renderer.OutputFile())
	if err != nil {
		return nil
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	return &result{
		id:    renderer.Id(),
		start: start,
		frame: frame,
		file:  data,
	}
}

func (i *instance) CreateRenderer() blender.Instance {
	// Create a temp file for the file
	tmpFile, err := ioutil.TempFile("", "weblender_")
	if err != nil {
		return nil
	}

	// Save the input and output name
	input := tmpFile.Name()
	output := fmt.Sprintf("%s_out_####", input)

	// Write the buffer
	tmpFile.Write(i.file)
	tmpFile.Close()

	return blender.New(i.id, input, output, i.renderer, i.start, i.frame)
}

func (i *instance) Id() uint64 {
	return i.id
}

func (i *instance) Renderer() string {
	return i.renderer
}

func (i *instance) Frame() (int, int) {
	return i.start, i.frame
}

func (i *instance) File() []byte {
	return i.file
}
