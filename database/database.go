package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

const conStr = "postgres://postgres:postgres@localhost:5432/flourish"

func ConnectToDatabase() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), conStr)
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}

// Insert new User with UserName and Password
func InsertUser(p *pgxpool.Pool, name string, password string) error {
	query := "INSERT INTO users (user_name, pass_word) VALUES ($1 $2)"
	_, err := p.Exec(context.Background(), query, name, password)

	if err != nil {
		msg := "error inserting user: " + err.Error()
		return errors.New(msg)
	}

	return nil
}

// Get User by Name and Login
func AuthenticateUser(p *pgxpool.Pool, name string, password string) (bool, error) {
	query := "SELECT user_id FROM users WHERE user_name = $1 AND pass_word = $2"

	id := -1
	err := p.QueryRow(context.Background(), query, name, password).Scan(&id)

	if err != nil {
		msg := "error validating user credentials: " + err.Error()
		return false, errors.New(msg)
	}

	return id != -1, nil
}

// Get Tasks by User Id
func GetTasksByUserId(p *pgxpool.Pool, user_id int) ([]Task, error) {
	query := "SELECT * FROM tasks WHERE user_id = $1"
	rows, err := p.Query(context.Background(), query, user_id)
	if err != nil {
		msg := "error getting tasks by user id: " + err.Error()
		return []Task{}, errors.New(msg)
	}

	var tasks []Task
	for rows.Next() {
		var task Task
		err = rows.Scan(&task)
		tasks = append(tasks, task)
	}

	if err != nil {
		msg := "error validating user credentials: " + err.Error()
		return []Task{}, errors.New(msg)
	}

	return tasks, nil
}

// Get Task Progress by User Id and Task Id
func GetTaskProgressByUserTaskId(p *pgxpool.Pool, userId int, taskId int) (TaskProgress, error) {
	query := "SELECT * FROM task_progress WHERE user_id = $1 AND task_id = $2;"
	var tp TaskProgress
	err := p.QueryRow(context.Background(), query, userId, taskId).Scan(&tp)
	if err != nil {
		msg := "error getting task progress by user and task id: " + err.Error()
		return TaskProgress{}, errors.New(msg)
	}

	return tp, nil
}
