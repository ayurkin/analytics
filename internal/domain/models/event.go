package models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	UUID            uuid.UUID `json:"uuid"`
	TaskId          int32     `json:"task_id"`
	Time            time.Time `json:"time"`
	Type            string    `json:"type"`
	User            string    `json:"user"`
	ApproversNumber int32     `json:"approvers_number"`
}
