package processor_helpers

import (
	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
)

func ValidateJob(job model.Job) error {
	if job.Value < 0 {
		return model.ErrInvalidValue
	}

	return nil
}

func JobToResult(job model.Job) model.Result {
	var result model.Result

	result.JobID = job.ID
	result.ProcessedValue = job.Value * 2

	err := ValidateJob(job)
	if err != nil {
		result.ErrorMessage = model.ErrInvalidValue
	} else {
		result.ErrorMessage = nil
	}

	return result
}
