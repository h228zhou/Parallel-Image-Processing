package scheduler

import (
	"bufio"
	"encoding/json"
	"os"
	"proj1/png"
)

func parseEffectsFile() ([]*png.ImageTask, error) {
	effectsPathFile := "../data/effects.txt"
	effectsFile, err := os.Open(effectsPathFile)
	if err != nil {
		return nil, err
	}
	defer effectsFile.Close()

	var tasks []*png.ImageTask
	scanner := bufio.NewScanner(effectsFile)
	for scanner.Scan() {
		line := scanner.Text()
		var task png.ImageTask
		err := json.Unmarshal([]byte(line), &task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// helper min function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
