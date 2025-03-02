package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sebastian-j-ibanez/flourish-backend/api"
	"github.com/sebastian-j-ibanez/flourish-backend/database"
)

func main() {
	p, err := database.ConnectToDatabase()
	if err != nil {
		msg := "error conecting to database: " + err.Error()
		panic(msg)
	}

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{
		"*",
	}

	r.Use(cors.New(config))

	r.POST("/login", api.LoginHandler(p))
	r.POST("/signup", api.SignupHandler(p))
	r.POST("/newtask", api.NewTaskHandler(p))
	r.POST("/updatetask", api.UpdateTaskHandler(p))
	r.POST("/todaytasks", api.TodoTaskHandler(p))
	r.POST("/usertasks", api.TaskListingHandler(p))
	r.POST("/treedata", api.TreeDataHandler(p))
	r.POST("/ping", api.Ping)

	r.Run(":8080")
}
