package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (c *Config) Start() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", getAllUsers)
		v1.POST("/users", createUser)
		v1.GET("/users/:id", getUserById)
		v1.PUT("/users/:id", updateUserById)
		v1.DELETE("/users/:id", deleteUserById)
	}

	listenPort := fmt.Sprintf(":%s", c.ListenPort)
	router.Run(listenPort)
}
