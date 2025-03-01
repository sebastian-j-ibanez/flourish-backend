package database

type Task struct {
	TaskId   int
	TaskCode string
	TaskName string
	UserIds  []int
}
