package blender

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	alreadyRunningError = errors.New("The process is already running.")

	regexProgressBlender = regexp.MustCompile("^Fra:([0-9]+) Mem:([0-9\\.]+[MG]?) \\(([0-9\\.]+[MG]?), peak ([0-9\\.]+[MG]?)\\) \\| [0-9a-zA-Z]+, Part ([0-9]+)\\-([0-9]+)$")
	regexProgressCycles  = regexp.MustCompile("^Fra:([0-9]+) Mem:([0-9\\.]+[MG]?) \\(([0-9\\.]+[MG]?), peak ([0-9\\.]+[MG]?)\\) \\| Mem: [0-9\\.]+[MG]?, Peak: [0-9\\.]+[MG]? \\| [0-9a-zA-Z]+, [0-9a-zA-Z]+ \\| Elapsed: [0-9:\\.]+ \\| Rendering \\| Path Tracing Tile ([0-9]+)/([0-9]+)$")
	regexSaved           = regexp.MustCompile("^Saved: ([0-9a-zA-Z\\./_\\-,]+) Time: ([0-9:\\.]+)$")
)

type Instance interface {
	StartRendering() (chan int, chan error)
	Progress() int
	Running() bool
	Cancel()
	InputFile() string
	OutputFile() string
	Time() (time.Duration, error)
	Id() uint64
	Frame() (int, int)
}

type instance struct {
	mutex sync.RWMutex

	process                               *exec.Cmd
	cur, max, start, frame                int
	jobId                                 uint64
	renderer, inputFile, outputFile, time string
}

func New(jobId uint64, inputFile, outputFile, renderer string, start, frame int) *instance {
	return &instance{
		cur:        -1,
		max:        -1,
		start:      start,
		frame:      frame,
		renderer:   renderer,
		inputFile:  inputFile,
		outputFile: outputFile,
	}
}

func (i *instance) StartRendering() (chan int, chan error) {
	// Create channels
	progress := make(chan int, 1)
	errors := make(chan error, 1)
	killProgress := make(chan bool, 1)

	// Wrap render()
	go func() {
		if err := i.render(); err != nil {
			errors <- err
		}
		close(killProgress)
		close(errors)
	}()

	// Create progess go routine
	go func() {
		defer close(progress)

		for {
			select {
			case <-killProgress:
				return
			case progress <- i.Progress():
			}
		}
	}()

	return progress, errors
}

func (i *instance) render() error {
	// Check if the process is already running
	if i.Running() {
		return alreadyRunningError
	}

	// Create the process
	i.process = exec.Command("blender", // Run blender
		"-b", i.inputFile, // The input file
		"-o", i.outputFile, // The output file
		"-s", fmt.Sprintf("%i", i.start), // The start frame
		"-f", fmt.Sprintf("+%i", i.frame), // The frame to render (relative to the start)
		"-E", i.renderer, // The renderer to use
		"-F", "PNG", // Output as PNG
		"-nojoystick", // No joystick support
		"-noaudio")    // No audio support

	// Get a reader
	output, err := i.process.StdoutPipe()
	if err != nil {
		return err
	}
	defer output.Close()

	// Wrap a scanner around the reader
	scanner := bufio.NewScanner(output)

	// Start blender
	if err = i.process.Start(); err != nil {
		return err
	}

	// Parse the output from the process
	for scanner.Scan() {
		// Check for read errors and read line
		if scanner.Err() != nil {
			return scanner.Err()
		}
		line := scanner.Text()

		// Scan for lines
		var parts [][]string

		// Scan for progress
		switch i.renderer {
		case "BLENDER_RENDER":
			parts = regexProgressBlender.FindAllStringSubmatch(line, 1)
		case "CYCLES":
			parts = regexProgressCycles.FindAllStringSubmatch(line, 1)
		}
		if len(parts) > 0 {
			i.mutex.Lock()
			if i.max < 0 {
				i.max, err = strconv.Atoi(parts[0][6])
				if err != nil {
					i.max = -2
				}
				i.cur = 0
			}
			i.cur++
			i.mutex.Unlock()

			// We found our match for this line
			continue
		}

		// Scan for saved
		parts = regexSaved.FindAllStringSubmatch(line, 1)
		if len(parts) > 0 {
			i.mutex.Lock()
			i.cur = i.max
			i.time = parts[0][2]
			i.mutex.Unlock()
		}

		//		fmt.Println("Progress:", i.Progress())
	}

	// Check if blender was successful
	if err = i.process.Wait(); err != nil {
		return err
	}
	return nil
}

// Progress calculates the progess of the blender instance in percent
func (i *instance) Progress() int {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	if i.cur < 0 || i.max < 0 {
		return 0
	}
	return (i.cur * 100) / i.max
}

func (i *instance) Cancel() {
	if !i.Running() {
		return
	}

	i.process.Process.Kill()
}

// Running checks if the renderer is running
func (i *instance) Running() bool {
	return i.process != nil && i.process.ProcessState != nil && !i.process.ProcessState.Exited()
}

func (i *instance) InputFile() string {
	return i.inputFile
}

func (i *instance) OutputFile() string {
	return i.outputFile
}

func (i *instance) Time() (time.Duration, error) {
	return time.ParseDuration(i.time)
}

func (i *instance) Id() uint64 {
	return i.jobId
}

func (i *instance) Frame() (int, int) {
	return i.start, i.frame
}

// Version tries to start blender and determin the version
func Version() string {
	// Prepare process
	proc := exec.Command("blender", "-v")

	// Prepare output reader
	out, err := proc.StdoutPipe()
	if err != nil {
		return ""
	}
	defer out.Close()
	scanner := bufio.NewScanner(out)

	// Start the process
	if err = proc.Start(); err != nil {
		return ""
	}

	for scanner.Scan() {
		line := scanner.Text()

		// Filter all lines not starting with "Blender "
		if !strings.HasPrefix(line, "Blender ") {
			continue
		}
		return strings.TrimSpace(strings.TrimPrefix(line, "Blender"))
	}

	return ""
}

// AsyncVersion is a wrapper that starts Version as a go routine
func AsyncVersion() chan string {
	// Create output channel
	out := make(chan string, 1)

	// Start the wrapper
	go func() {
		if ver := Version(); ver != "" {
			out <- ver
		}
		close(out)
	}()

	return out
}
