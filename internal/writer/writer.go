package writer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
)

func WriteSummary(results []model.Result, path string) error {
	file, err := os.Create(path + "/summary.json")
	if err != nil {
		return fmt.Errorf("error while creating file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(results); err != nil {
		return err
	}

	return nil
}
