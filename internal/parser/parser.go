package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
)

func ParseFile(filename string) ([]model.Job, error) {
	jobs := make([]model.Job, 3)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
	}(f)

	if err := json.NewDecoder(f).Decode(&jobs); err != nil {
		return nil, fmt.Errorf("loadConfig: %w", err)
	}

	fmt.Println(jobs)

	return jobs, nil
}
