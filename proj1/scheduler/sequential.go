package scheduler

import (
	"fmt"
	"proj1/png"
	"strings"
)

func RunSequential(config Config) {
	dataDirs := strings.Split(config.DataDirs, "+")
	tasks, err := parseEffectsFile()
	if err != nil {
		fmt.Println("Error parsing effects file:", err)
		return
	}

	for _, dataDir := range dataDirs {
		for _, task := range tasks {
			// Construct the input file path
			inputFilePath := fmt.Sprintf("../data/in/%s/%s", dataDir, task.InPath)

			// Load the image
			img, err := png.Load(inputFilePath)
			if err != nil {
				fmt.Println("Error loading image:", err)
				continue
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

			// Construct the output file path (prepend dataDir identifier)
			outputFilePath := fmt.Sprintf("../data/out/%s_%s", dataDir, task.OutPath)

			// Save the image
			err = img.Save(outputFilePath)
			if err != nil {
				fmt.Println("Error saving image:", err)
			}
		}
	}
}
