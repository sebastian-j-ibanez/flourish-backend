package database

import "time"

type TaskProgress struct {
	TaskProgressId int
	UserId         int
	TaskId         int
	Status         bool
	TaskDate       time.Time
}
