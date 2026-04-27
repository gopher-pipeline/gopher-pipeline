package proccesor

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
)

func TestTransform(t *testing.T) {
	testUUID, _ := uuid.NewUUID()

	tests := []struct {
		name     string
		job      model.Job
		expected model.Result
		wantErr  bool
	}{
		{"valid job", model.Job{
			ID:       testUUID,
			Filename: "testfilename",
			Value:    500,
			Status:   "pending",
		}, model.Result{
			JobID:          testUUID,
			ProcessedValue: 1000,
			ErrorMessage:   nil,
		}, false,
		}, {
			"invalid value", model.Job{
				ID:       testUUID,
				Filename: "testfilename",
				Value:    -5,
				Status:   "done",
			}, model.Result{
				JobID:          testUUID,
				ProcessedValue: -5,
				ErrorMessage:   model.ErrInvalidValue,
			}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Transform(tt.job)
			if result != tt.expected {
				fmt.Errorf("error while transform method working")
			}
			if (err != nil) != tt.wantErr {
				fmt.Errorf("got err=%v, wantErr=%v", err, tt.wantErr)
			}
		})
	}
}
