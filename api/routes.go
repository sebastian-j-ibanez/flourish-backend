package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sebastian-j-ibanez/flourish-backend/database"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(p *pgxpool.Pool) gin.HandlerFunc {
	data := LoginData{}

	login := func(c *gin.Context) {
		if err := c.ShouldBindJSON(&data); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		success, err := database.AuthenticateUser(p, data.Username, data.Password)
		if !success || err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		c.Status(http.StatusOK)
	}

	return login
}

func Ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
