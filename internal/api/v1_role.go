package api

import (
	"time"

	"shier/internal/models"
	"shier/pkg/db"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

func getAllRoles(c *gin.Context) {
	var roles []models.Role
	err := db.GetConnection().Model(&roles).Select() // get all roles

	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
	}

	result := map[string]interface{}{
		"roles": roles,
	}

	httpOkResponse(c, result)
}

func createRole(c *gin.Context) {
	form := &struct {
		Name string `form:"name" json:"name"`
	}{}
	c.BindJSON(form)

	// form validation
	err := validation.Errors{
		"name": validation.Validate(form.Name, validation.Required),
	}.Filter()

	if err != nil {
		httpValidationErrorResponse(c, err.Error())
		return
	}

	role := models.Role{
		Name:      form.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = db.GetConnection().Insert(&role)
	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	result := map[string]interface{}{
		"role": "Create role successfully",
	}

	httpOkResponse(c, result)
}

func getRoleById(c *gin.Context) {
	var role models.Role

	id := c.Param("id")

	// get from database
	err := db.GetConnection().Model(&role).Where("id = ?", id).Select()
	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	result := map[string]interface{}{
		"roles": role,
	}

	httpOkResponse(c, result)
}

func updateRoleById(c *gin.Context) {
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
		httpValidationErrorResponse(c, err.Error())
		return
	}

	err = db.GetConnection().Model(&role).Where("id = ?", id).Select()
	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	role = models.Role{
		ID:        role.ID,
		Name:      form.Name,
		UpdatedAt: time.Now(),
	}

	_, err = db.GetConnection().Model(&role).
		Column("name").
		WherePK().Returning("*").Update()

	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	result := map[string]interface{}{
		"role": "Update role successfully",
	}

	httpOkResponse(c, result)
}

func deleteRoleById(c *gin.Context) {
	var role models.Role

	id := c.Param("id")
	err := validation.Errors{
		"id": validation.Validate(id, validation.Required),
	}.Filter()
	if err != nil {
		httpValidationErrorResponse(c, err.Error())
		return
	}

	err = db.GetConnection().Model(&role).Where("id = ?", id).Select()
	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	err = db.GetConnection().Delete(&role)
	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	result := map[string]interface{}{
		"role": "Delete role successfully",
	}

	httpOkResponse(c, result)
}
