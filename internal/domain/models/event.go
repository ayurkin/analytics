package models

import "time"

type Event struct {
	TaskId          int32
	Time            time.Time
	Type            string
	User            string
	ApproversNumber int32
}
