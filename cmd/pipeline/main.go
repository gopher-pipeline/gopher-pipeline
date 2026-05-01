package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
	"github.com/gopher-pipeline/gopher-pipeline/internal/parser"
	"github.com/gopher-pipeline/gopher-pipeline/internal/pipeline"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generateFiles(3, 5)

	bufJobsCh := make(chan model.Job, 20)
	bufErrCh := make(chan error, 20)
	bufResCh := make(chan model.Result, 20)

	paths, err := filepath.Glob("data/*.json")
	if err != nil {
		fmt.Printf("Error: ur files are shit")
	}
	go func() {
		for _, f := range paths {
			parsedFile, err := parser.ParseFile(f)

			if err != nil {
				bufErrCh <- err
				continue
			}

			for _, job := range parsedFile {
				bufJobsCh <- job
			}
		}

		close(bufJobsCh)
	}()

	p := pipeline.NewPipeline(bufJobsCh, bufResCh, bufErrCh, 5)
	go p.Run(ctx)

	results := make([]model.Result, 0)
	for res := range bufResCh {
		results = append(results, res)
	}

	err = writer.WriteSummary(results, "out")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

}
