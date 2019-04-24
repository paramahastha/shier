package api

import (
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/paramahastha/shier/internal/models"
	"github.com/paramahastha/shier/pkg/db"
)

func getAllUsers(c *gin.Context) {
	var users []models.User
	err := db.GetConnection().Model(&users).Select() // get all users

	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
	}

	result := map[string]interface{}{
		"users": users,
	}

	httpOkResponse(c, result)
}

func createUser(c *gin.Context) {
	form := &struct {
		FirstName string `form:"first_name" json:"first_name"`
		LastName  string `form:"last_name" json:"last_name"`
		Email     string `form:"email" json:"email"`
		Password  string `form:"password" json:"password"`
		Confirm   string `form:"confirm" json:"confirm"`
		Role      string `form:"role" json:"role"`
	}{}
	c.Bind(form)

	// form validation
	err := validation.Errors{
		"first_name": validation.Validate(form.FirstName, validation.Required),
		"last_name":  validation.Validate(form.LastName, validation.Required),
		"email":      validation.Validate(form.Email, validation.Required, is.Email),
		"password":   validation.Validate(form.Password, validation.Required),
		"confirm":    validation.Validate(form.Confirm, validation.In(form.Password).Error("Your password and confirmation password do not match.")),
		"role":       validation.Validate(form.Role, validation.Required, validation.In("user", "admin").Error("must be a 'user' or 'admin'")),
	}.Filter()

	if err != nil {
		httpValidationErrorResponse(c, err.Error())
		return
	}

	user := models.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Password:  form.Password,
		Confirm:   form.Confirm,
		Role:      form.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = db.GetConnection().Insert(&user)
	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	result := map[string]interface{}{
		"users": user,
	}

	httpOkResponse(c, result)
}

func getUserById(c *gin.Context) {
	var user models.User

	id := c.Param("id")

	// get from database
	err := db.GetConnection().Model(&user).Where("id = ?", id).Select()
	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	result := map[string]interface{}{
		"users": user,
	}

	httpOkResponse(c, result)
}

func updateUserById(c *gin.Context) {
	var user models.User

	form := &struct {
		FirstName string `form:"first_name" json:"first_name"`
		LastName  string `form:"last_name" json:"last_name"`
		Email     string `form:"email" json:"email"`
		Password  string `form:"password" json:"password"`
		Confirm   string `form:"confirm" json:"confirm"`
		Role      string `form:"role" json:"role"`
	}{}
	id := c.Param("id")
	c.Bind(form)

	// form validation
	err := validation.Errors{
		"first_name": validation.Validate(form.FirstName, validation.Required),
		"last_name":  validation.Validate(form.LastName, validation.Required),
		"email":      validation.Validate(form.Email, validation.Required, is.Email),
		"password":   validation.Validate(form.Password, validation.Required),
		"confirm":    validation.Validate(form.Confirm, validation.In(form.Password).Error("Your password and confirmation password do not match.")),
		"role":       validation.Validate(form.Role, validation.Required, validation.In("user", "admin").Error("must be a 'user' or 'admin'")),
	}.Filter()

	if err != nil {
		httpValidationErrorResponse(c, err.Error())
		return
	}

	err = db.GetConnection().Model(&user).Where("id = ?", id).Select()
	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	user = models.User{
		ID:        user.ID,
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Password:  form.Password,
		Confirm:   form.Confirm,
		Role:      form.Role,
	}

	_, err = db.GetConnection().Model(&user).
		Column("first_name").
		Column("last_name").
		Column("email").
		Column("password").
		Column("confirm").
		Column("role").
		WherePK().Returning("*").Update()

	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	result := map[string]interface{}{
		"user": user,
	}

	httpOkResponse(c, result)
}

func deleteUserById(c *gin.Context) {
	var user models.User

	id := c.Param("id")
	err := validation.Errors{
		"id": validation.Validate(id, validation.Required),
	}.Filter()
	if err != nil {
		httpValidationErrorResponse(c, err.Error())
		return
	}

	err = db.GetConnection().Model(&user).Where("id = ?", id).Select()
	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	db.GetConnection().Delete(&user)

	result := map[string]interface{}{
		"user": "Delete user successfully",
	}

	httpOkResponse(c, result)
}
