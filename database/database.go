package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sebastian-j-ibanez/flourish-backend/date"
)

const conStr = "postgres://postgres:postgres@localhost:5432/flourish"

type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Task struct {
	TaskId   int    `json:"taskId"`
	TaskCode string `json:"taskCode"`
	TaskName string `json:"taskName"`
	UserIds  []int  `json:"userIds"`
}

type TaskProgress struct {
	TaskProgressId int       `json:"taskProgressId"`
	UserId         int       `json:"userId"`
	TaskId         int       `json:"taskId"`
	Status         bool      `json:"status"`
	TaskDate       time.Time `json:"taskDate"`
}

// Abstracted structs
type ToDoTask struct {
	TaskId        int    `json:"taskId"`
	TaskName      string `json:"taskName"`
	TaskCompleted bool   `json:"taskCompleted"`
}

type TaskListing struct {
	TaskId        int    `json:"taskId"`
	TaskName      string `json:"taskName"`
	UserNum       int    `json:"userNum"`
	DaysCompleted int    `json:"daysCompleted"`
}

type TreeData struct {
	UserId        int    `json:"userId"`
	TaskId        int    `json:"taskId"`
	TaskName      string `json:"taskName"`
	TaskCode      string `json:"taskCode"`
	Days          []bool `json:"days"`
	DaysCompleted int    `json:"daysCompleted"`
}

func ConnectToDatabase() (*pgxpool.Pool, error) {
	p, err := pgxpool.New(context.Background(), conStr)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Get User by Name and Login
func AuthenticateUser(p *pgxpool.Pool, name string, password string) (User, error) {
	query := "SELECT user_id, user_name, pass_word FROM users WHERE user_name = $1 AND pass_word = $2"

	var user User
	err := p.QueryRow(context.Background(), query, name, password).Scan(&user.UserId, &user.UserName, &user.Password)

	if err != nil {
		msg := "error validating user credentials: " + err.Error()
		return User{}, errors.New(msg)
	}

	return user, nil
}

// Insert new User with UserName and Password
func InsertUser(p *pgxpool.Pool, name string, password string) (User, error) {
	var user User

	query := "INSERT INTO users (user_name, pass_word) VALUES ($1, $2) RETURNING user_id, user_name, pass_word"

	err := p.QueryRow(context.Background(), query, name, password).Scan(&user.UserId, &user.UserName, &user.Password)
	if err != nil {
		msg := "error inserting user: " + err.Error()
		return User{}, errors.New(msg)
	}

	return user, nil
}

func NewTask(p *pgxpool.Pool, userId int, taskName string, taskCode string) error {
	tx, err := p.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := "INSERT INTO tasks (user_ids, task_name, task_code) VALUES ($1, $2, $3) RETURNING task_id"
	var taskId int
	err = tx.QueryRow(context.Background(), query, []int{userId}, taskName, taskCode).Scan(&taskId)
	if err != nil {
		return err
	}

	query = "INSERT INTO task_progress (user_id, task_id, status, task_date) VALUES ($1, $2, $3, $4)"
	_, err = tx.Exec(context.Background(), query, userId, taskId, false, date.GetToday())
	if err != nil {
		return err
	}

	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}

func UpdateTask(p *pgxpool.Pool, userId int, taskId int, status bool) error {
	query := "UPDATE task_progress SET status = $1 WHERE user_id = $2 AND task_id = $3"
	_, err := p.Exec(
		context.Background(),
		query,
		status,
		userId,
		taskId,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(p *pgxpool.Pool, userId int, taskId int) error {
	tx, err := p.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `
	UPDATE tasks 
	SET user_ids = array_remove(user_ids, $1) 
	WHERE task_id = $2
`
	_, err = tx.Exec(context.Background(), query, userId, taskId)
	if err != nil {
		return err
	}

	var userIds []int
	query = `
        SELECT user_ids 
        FROM tasks 
        WHERE task_id = $1
    `
	err = tx.QueryRow(context.Background(), query, taskId).Scan(&userIds)
	if err != nil {
		return err
	}

	query = `
	DELETE FROM task_progress
	WHERE user_id = $1 AND task_id = $2
	`
	_, err = tx.Exec(context.Background(), query, userId, taskId)
	if err != nil {
		return err
	}

	if len(userIds) <= 0 {
		query = `
			DELETE FROM tasks 
			WHERE task_id = $1
		`
		_, err = tx.Exec(context.Background(), query, taskId)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}

func JoinTask(p *pgxpool.Pool, userId int, taskCode string) error {
	tx, err := p.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `
		UPDATE tasks
		SET user_ids = array_append(user_ids, $1)
		WHERE task_code = $2
	`
	_, err = tx.Exec(
		context.Background(),
		query,
		userId,
		taskCode,
	)
	if err != nil {
		return err
	}

	var taskId int
	query = "SELECT task_id FROM tasks WHERE task_code = $1"
	err = tx.QueryRow(context.Background(), query, taskCode).Scan(&taskId)
	if err != nil {
		return err
	}

	query = "INSERT INTO task_progress (user_id, task_id, status, task_date) VALUES ($1, $2, $3, $4)"
	_, err = tx.Exec(
		context.Background(),
		query,
		userId,
		taskId,
		false,
		date.GetToday(),
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}

// Get Tasks by User Id
func GetTasksByUserId(p *pgxpool.Pool, user_id int) ([]Task, error) {
	query := "SELECT tasks.task_id, tasks.task_name, tasks.task_code, tasks.user_ids FROM tasks WHERE $1 = ANY(user_ids)"
	rows, err := p.Query(context.Background(), query, user_id)
	if err != nil {
		msg := "error getting tasks by user id: " + err.Error()
		return []Task{}, errors.New(msg)
	}

	var tasks []Task
	for rows.Next() {
		var task Task
		err = rows.Scan(
			&task.TaskId,
			&task.TaskName,
			&task.TaskCode,
			&task.UserIds,
		)
		if err != nil {
			msg := "error scanning task rows: " + err.Error()
			return []Task{}, errors.New(msg)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Get Task Progress by User Id and Task Id
func GetToDoTaskByUserIdTaskId(p *pgxpool.Pool, userId int, taskId int) (ToDoTask, error) {
	query := `SELECT tp.task_id, t.task_name, tp.status
	FROM task_progress tp
	JOIN users u ON tp.user_id = u.user_id
	JOIN tasks t ON tp.task_id = t.task_id
	WHERE tp.user_id = $1 AND tp.task_id = $2 AND tp.task_date = $3`

	var task ToDoTask
	todayDate := date.GetToday()
	err := p.QueryRow(
		context.Background(),
		query, userId,
		taskId,
		todayDate,
	).Scan(
		&task.TaskId,
		&task.TaskName,
		&task.TaskCompleted,
	)
	if err != nil {
		msg := "error getting task progress by user and task id: " + err.Error()
		return ToDoTask{}, errors.New(msg)
	}

	return task, nil
}

func GetTaskListingsByUserIdTaskId(p *pgxpool.Pool, userId int) ([]TaskListing, error) {
	query := `SELECT
		tp.task_id,
		t.task_name,
		array_length(t.user_ids, 1),
		SUM(CASE WHEN tp.status THEN 1 ELSE 0 END)
	FROM task_progress tp
	JOIN tasks t ON tp.task_id = t.task_id
	WHERE
		tp.user_id = $1
	GROUP BY
		tp.task_id, t.task_name, t.user_ids;`

	rows, err := p.Query(
		context.Background(),
		query, userId,
	)
	if err != nil {
		return []TaskListing{}, nil
	}

	var listings []TaskListing
	for rows.Next() {
		var listing TaskListing
		err = rows.Scan(
			&listing.TaskId,
			&listing.TaskName,
			&listing.UserNum,
			&listing.DaysCompleted,
		)
		if err != nil {
			return []TaskListing{}, nil
		}
		listings = append(listings, listing)
	}

	return listings, nil
}

func GetTreeDataByTaskId(p *pgxpool.Pool, taskId int) ([]TreeData, error) {
	query := `SELECT
		tp.user_id,
		tp.task_id,
		t.task_name,
		t.task_code,
		ARRAY_AGG(tp.status ORDER BY tp.task_date DESC) AS statuses,
	(SELECT COUNT(*) FROM UNNEST(ARRAY_AGG(tp.status)) AS status WHERE status = TRUE) AS true_count
	FROM
		task_progress tp
	JOIN
		tasks t ON tp.task_id = t.task_id
	WHERE
		t.task_id = $1
	GROUP BY
		tp.user_id, tp.task_id, t.task_name, t.task_code;`

	rows, err := p.Query(
		context.Background(),
		query,
		taskId,
	)
	if err != nil {
		return []TreeData{}, nil
	}

	var treeData []TreeData
	for rows.Next() {
		var datum TreeData
		err = rows.Scan(
			&datum.UserId,
			&datum.TaskId,
			&datum.TaskName,
			&datum.TaskCode,
			&datum.Days,
			&datum.DaysCompleted,
		)
		if err != nil {
			return []TreeData{}, err
		}
		treeData = append(treeData, datum)
	}

	return treeData, nil
}
