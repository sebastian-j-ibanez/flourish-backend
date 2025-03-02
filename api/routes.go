package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sebastian-j-ibanez/flourish-backend/code"
	"github.com/sebastian-j-ibanez/flourish-backend/database"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRequest struct {
	UserId int `json:"userId"`
}

type UserTaskRequest struct {
	UserId int `json:"userId"`
	TaskId int `json:"taskId"`
}

type UserTaskStatusRequest struct {
	UserId int  `json:"userId"`
	TaskId int  `json:"taskId"`
	Status bool `json:"status"`
}

type NewTaskRequest struct {
	UserId   int    `json:"userId"`
	TaskName string `json:"taskName"`
}

type TaskRequest struct {
	TaskId int `json:"taskId"`
}

func Ping(c *gin.Context) {
	c.Status(http.StatusOK)
}

func LoginHandler(p *pgxpool.Pool) gin.HandlerFunc {
	data := LoginData{}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&data); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		user, err := database.AuthenticateUser(p, data.Username, data.Password)
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func SignupHandler(p *pgxpool.Pool) gin.HandlerFunc {
	data := LoginData{}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&data); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		user, err := database.InsertUser(p, data.Username, data.Password)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func NewTaskHandler(p *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request NewTaskRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err := database.NewTask(
			p,
			request.UserId,
			request.TaskName,
			code.GenerateCode(),
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.Status(http.StatusOK)
	}
}

func UpdateTaskHandler(p *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request UserTaskStatusRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err := database.UpdateTask(
			p,
			request.UserId,
			request.TaskId,
			request.Status,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.Status(http.StatusOK)
	}
}

func DeleteTaskHandler(p *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request UserTaskRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		err := database.DeleteTask(
			p,
			request.UserId,
			request.TaskId,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.Status(http.StatusOK)
	}
}

func TodoTaskHandler(p *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request UserRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		userId := request.UserId

		tasks, err := database.GetTasksByUserId(p, userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var toDoTasks []database.ToDoTask
		for _, task := range tasks {
			toDoTask, err := database.GetToDoTaskByUserIdTaskId(p, userId, task.TaskId)
			if err == nil {
				toDoTasks = append(toDoTasks, toDoTask)
			}
		}

		if len(toDoTasks) <= 0 {
			msg := errors.New("no matching rows")
			c.JSON(http.StatusNoContent, msg)
		}

		c.JSON(http.StatusOK, toDoTasks)
	}
}

func TaskListingHandler(p *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request UserRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		listings, err := database.GetTaskListingsByUserIdTaskId(p, request.UserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, listings)
	}
}

func TreeDataHandler(p *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request TaskRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		treeData, err := database.GetTreeDataByTaskId(p, request.TaskId)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, treeData)
	}
}
