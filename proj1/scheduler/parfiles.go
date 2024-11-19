package scheduler

import (
	"fmt"
	"proj1/locks"
	"proj1/png"
	"strings"
	"sync"
)

func RunParallelFiles(config Config) {
	dataDirs := strings.Split(config.DataDirs, "+")
	var taskQueue []*png.ImageTask

	// Parse effects.txt once
	tasks, err := parseEffectsFile()
	if err != nil {
		fmt.Println("Error parsing effects file:", err)
		return
	}

	// Create tasks for each data directory
	for _, dataDir := range dataDirs {
		for _, task := range tasks {
			newTask := *task // Copy the task
			newTask.DataDir = dataDir
			taskQueue = append(taskQueue, &newTask)
		}
	}

	numThreads := min(config.ThreadCount, len(taskQueue))
	var wg sync.WaitGroup
	wg.Add(numThreads)

	lock := &locks.TASLock{}

	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()
			for {
				var task *png.ImageTask

				// Safely access the queue
				lock.Lock()
				if len(taskQueue) > 0 {
					task = taskQueue[0]
					taskQueue = taskQueue[1:]
				}
				lock.Unlock()

				if task == nil {
					// No more tasks
					return
				}

				// Process the task
				processTask(task)
			}
		}()
	}

	wg.Wait()
}

func processTask(task *png.ImageTask) {
	// Construct the input file path
	inputFilePath := fmt.Sprintf("../data/in/%s/%s", task.DataDir, task.InPath)

	// Load the image
	img, err := png.Load(inputFilePath)
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	// Apply effects
	for _, effect := range task.Effects {
		switch effect {
		case "G":
			img.ApplyGrayscale(img.Bounds.Min.Y, img.Bounds.Max.Y)
		case "S":
			img.ApplyConvolution(png.SharpenKernel, img.Bounds.Min.Y, img.Bounds.Max.Y)
		case "E":
			img.ApplyConvolution(png.EdgeDetectionKernel, img.Bounds.Min.Y, img.Bounds.Max.Y)
		case "B":
			img.ApplyConvolution(png.BlurKernel, img.Bounds.Min.Y, img.Bounds.Max.Y)
		default:
			fmt.Println("Unknown effect:", effect)
		}
		// Added swap buffer to ensure next effect applies to the output of this effect
		img.SwapBuffers()
	}

	// Construct the output file path
	outputFilePath := fmt.Sprintf("../data/out/%s_%s", task.DataDir, task.OutPath)

	// Save the image
	err = img.Save(outputFilePath)
	if err != nil {
		fmt.Println("Error saving image:", err)
	}
}
