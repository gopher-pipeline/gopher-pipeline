package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/google/uuid"
	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
	"github.com/gopher-pipeline/gopher-pipeline/internal/parser"
	"github.com/gopher-pipeline/gopher-pipeline/internal/proccesor"
	"github.com/gopher-pipeline/gopher-pipeline/internal/writer"
)

// создать 3 пустых файла
// форик на 300 который будет создавать 300 джобов
// 300 джобов будут заполнять 3 массива
// через encode записать 3 массива в 3 файла
func generateFiles(numFiles int, numJobs int) {
	for i := 0; i < numFiles; i++ {
		var filename string = fmt.Sprintf("file_%d.json", i)
		file, err := os.Create("data/" + filename)
		if err != nil {
			log.Fatal("Error creating file: ", err)
		}

		jobs := generateJobs(numJobs, i)

		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")

		if err := enc.Encode(jobs); err != nil {
			log.Printf("encode error: %v", err)
		}

		defer file.Close()
	}
}

func generateJobs(numJobs int, fileNum int) []model.Job {
	jobs := make([]model.Job, numJobs)
	statuses := []string{"pending", "done", "error", "dolbaeb"}
	for i := 0; i < numJobs; i++ {
		id, _ := uuid.NewUUID()
		currentJob := model.Job{
			ID:       id,
			Filename: fmt.Sprintf("file_%d.json", fileNum),
			Value:    rand.Intn(201) - 100,
			Status:   statuses[rand.Intn(len(statuses))],
		}
		jobs[i] = currentJob
	}
	return jobs
}
func main() {
	// if err := os.RemoveAll("data/"); err != nil {
	// 	fmt.Errorf("Error: %v", err)
	// }
	generateFiles(3, 25)
	jobs, _ := parser.ParseFile("data/file_0.json")

	results := make([]model.Result, 0)

	for _, job := range jobs {
		current, _ := proccesor.Transform(job)
		results = append(results, current)
	}

	err := writer.WriteSummary(results, "data")
	if err != nil {
		fmt.Println(err)
	}

}
