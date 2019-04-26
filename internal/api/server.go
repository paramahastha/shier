package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": map[string]string{
				"message": "Ok",
			},
		})
	})

	v1 := router.Group("/api/v1")
	{
		// users api
		v1.GET("/users", getAllUsers)
		v1.POST("/users", createUser)
		v1.GET("/users/:id", getUserById)
		v1.PUT("/users/:id", updateUserById)
		v1.DELETE("/users/:id", deleteUserById)

		// role api
		v1.GET("/roles", getAllRoles)
		v1.POST("/roles", createRole)
		v1.GET("/roles/:id", getRoleById)
		v1.PUT("/roles/:id", updateRoleById)
		v1.DELETE("/roles/:id", deleteRoleById)
	}

	return router
}

func (c *Config) Start() {
	router := SetupRouter()

	listenPort := fmt.Sprintf(":%s", c.ListenPort)
	router.Run(listenPort)
}
