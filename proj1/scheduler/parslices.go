package scheduler

import (
	"fmt"
	"proj1/png"
	"strings"
	"sync"
)

func RunParallelSlices(config Config) {
	tasks, err := parseEffectsFile()
	if err != nil {
		fmt.Println("Error parsing effects file:", err)
		return
	}

	dataDirs := strings.Split(config.DataDirs, "+")
	for _, dataDir := range dataDirs {
		for _, task := range tasks {
			task.DataDir = dataDir
			processImage(task, config.ThreadCount)
		}
	}
}

func processImage(task *png.ImageTask, numThreads int) {
	inputFilePath := fmt.Sprintf("../data/in/%s/%s", task.DataDir, task.InPath)
	outputFilePath := fmt.Sprintf("../data/out/%s_%s", task.DataDir, task.OutPath)

	// Load the image
	img, err := png.Load(inputFilePath)
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	// Get image dimensions
	bounds := img.Bounds
	height := bounds.Max.Y - bounds.Min.Y
	sliceHeight := (height + numThreads - 1) / numThreads // Ceiling division

	// Apply each effect
	for _, effect := range task.Effects {
		var wg sync.WaitGroup
		wg.Add(numThreads)

		for i := 0; i < numThreads; i++ {
			minY := bounds.Min.Y + i*sliceHeight
			maxY := min(minY+sliceHeight, bounds.Max.Y)

			go func(minY, maxY int) {
				defer wg.Done()
				switch effect {
				case "G":
					img.ApplyGrayscale(minY, maxY)
				case "S":
					img.ApplyConvolution(png.SharpenKernel, minY, maxY)
				case "E":
					img.ApplyConvolution(png.EdgeDetectionKernel, minY, maxY)
				case "B":
					img.ApplyConvolution(png.BlurKernel, minY, maxY)
				}
			}(minY, maxY)
		}

		wg.Wait()

		// Swap buffer after each effect
		img.SwapBuffers()
	}

	// Save the processed image
	err = img.Save(outputFilePath)
	if err != nil {
		fmt.Println("Error saving image:", err)
	}
}
