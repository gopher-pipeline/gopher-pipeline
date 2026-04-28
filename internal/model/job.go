package model

import (
	"github.com/google/uuid"
)

type Job struct {
	ID       uuid.UUID `json:"id"`
	Filename string    `json:"filename"`
	Value    int       `json:"value"`
	Status   string    `json:"status"`
}

type Result struct {
	JobID          uuid.UUID `json:"job_id"`
	ProcessedValue int       `json:"processed_value"`
	ErrorMessage   error     `json:"error_message"`
}
