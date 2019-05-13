package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) SetupRouter(port string) {
	s.Router = gin.Default()

	s.Router.POST("/test", func(c *gin.Context) {
		form := &struct {
			Message string `form:"message" json:"message"`
		}{}
		c.BindJSON(form)

		c.JSON(http.StatusOK, gin.H{
			"message": form.Message,
		})
	})

	v1 := s.Router.Group("/api/v1")
	{
		// auth api
		v1.POST("/sign-up", s.signup)
		v1.POST("/sign-in", s.signin)

		// users api
		v1.GET("/users", s.getAllUsers)
		v1.POST("/users", s.createUser)
		v1.GET("/users/:id", s.getUserById)
		v1.PUT("/users/:id", s.updateUserById)
		v1.DELETE("/users/:id", s.deleteUserById)

		// role api
		v1.GET("/roles", s.getAllRoles)
		v1.POST("/roles", s.createRole)
		v1.GET("/roles/:id", s.getRoleById)
		v1.PUT("/roles/:id", s.updateRoleById)
		v1.DELETE("/roles/:id", s.deleteRoleById)
	}

	listenPort := fmt.Sprintf(":%s", port)
	s.Router.Run(listenPort)
}
