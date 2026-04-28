package proccesor

import (
	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
	processor_helpers "github.com/gopher-pipeline/gopher-pipeline/internal/proccesor/helpers"
)

func Transform(job model.Job) (model.Result, error) {
	res := processor_helpers.JobToResult(job)

	if res.ErrorMessage != nil {
		return res, res.ErrorMessage
	}

	return res, nil
}
