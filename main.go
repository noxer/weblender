package main

import (
	"fmt"
	"time"
	"weblender/blender"
)

func main() {
	// Get the blender version
	//	verChan := blender.AsyncVersion()

	// Run a test
	// ToDo: Remove
	inst := blender.New(0, "test.blend", "test_out_", "BLENDER_RENDER", 1, 2)
	prog, errs := inst.StartRendering()

	for pr := range prog {
		fmt.Printf("Progress: %d%%\n", pr)
		time.Sleep(time.Second)
	}

	for err := range errs {
		fmt.Printf("Error: %s", err.Error())
	}
}
