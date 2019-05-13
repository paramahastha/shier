package api

import (
	"net/http"
	"time"

	"shier/internal/models"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (s *Server) getAllRoles(c *gin.Context) {
	var roles []models.Role
	err := s.DB.Model(&roles).Select() // get all roles

	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
	}

	result := map[string]interface{}{
		"data": roles,
	}

	httpSuccessResponse(c, result["data"], http.StatusOK, "-")
}

func (s *Server) createRole(c *gin.Context) {
	form := &struct {
		Name string `form:"name" json:"name"`
	}{}
	c.BindJSON(form)

	// form validation
	err := validation.Errors{
		"name": validation.Validate(form.Name, validation.Required),
	}.Filter()

	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusBadRequest, "-")
		return
	}

	role := models.Role{
		Name:      form.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.DB.Insert(&role)
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	httpSuccessResponse(c, nil, http.StatusCreated, "Create role successfully")
}

func (s *Server) getRoleById(c *gin.Context) {
	var role models.Role

	id := c.Param("id")

	// get from database
	err := s.DB.Model(&role).Where("id = ?", id).Select()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	result := map[string]interface{}{
		"data": role,
	}

	httpSuccessResponse(c, result["data"], http.StatusOK, "-")
}

func (s *Server) updateRoleById(c *gin.Context) {
	var role models.Role

	form := &struct {
		Name string `form:"name" json:"name"`
	}{}
	id := c.Param("id")
	c.Bind(form)

	// form validation
	err := validation.Errors{
		"name": validation.Validate(form.Name, validation.Required, validation.In("user", "admin").Error("must be a 'user' or 'admin'")),
	}.Filter()

	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusBadRequest, "-")
		return
	}

	err = s.DB.Model(&role).Where("id = ?", id).Select()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	role = models.Role{
		ID:        role.ID,
		Name:      form.Name,
		UpdatedAt: time.Now(),
	}

	_, err = s.DB.Model(&role).
		Column("name").
		WherePK().Returning("*").Update()

	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	httpSuccessResponse(c, nil, http.StatusOK, "Update role successfully")
}

func (s *Server) deleteRoleById(c *gin.Context) {
	var role models.Role

	id := c.Param("id")
	err := validation.Errors{
		"id": validation.Validate(id, validation.Required),
	}.Filter()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusBadRequest, "-")
		return
	}

	err = s.DB.Model(&role).Where("id = ?", id).Select()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	err = s.DB.Delete(&role)
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	httpSuccessResponse(c, nil, http.StatusOK, "Delete role successfully")
}
