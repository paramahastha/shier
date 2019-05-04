package api

import (
	"fmt"
	"time"

	"shier/internal/models"
	"shier/pkg/db"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func signup(c *gin.Context) {
	var role models.Role

	form := &struct {
		FirstName string `form:"first_name" json:"first_name"`
		LastName  string `form:"last_name" json:"last_name"`
		Email     string `form:"email" json:"email"`
		Password  string `form:"password" json:"password"`
		Confirm   string `form:"confirm" json:"confirm"`
		Roles     []int  `form:"roles" json:"roles"`
	}{}
	c.BindJSON(form)

	// form validation
	err := validation.Errors{
		"first_name": validation.Validate(form.FirstName, validation.Required),
		"last_name":  validation.Validate(form.LastName, validation.Required),
		"email":      validation.Validate(form.Email, validation.Required, is.Email),
		"password":   validation.Validate(form.Password, validation.Required),
		"confirm":    validation.Validate(form.Confirm, validation.Required, validation.In(form.Password).Error("Your password and confirmation password do not match")),
		"roles":      validation.Validate(form.Roles, validation.Required, validation.Length(1, 2)),
	}.Filter()

	if err != nil {
		httpValidationErrorResponse(c, err.Error())
		return
	}

	// Check existing role in db
	for _, val := range uniqueNum(form.Roles) {

		err = db.GetConnection().Model(&role).Where("id = ?", val).Select()
		if err != nil {
			httpValidationErrorResponse(c, fmt.Sprintf("pg: no role (%d) in result set", val))
			return
		}
	}

	user := models.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Password:  form.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = db.GetConnection().Insert(&user)
	if err != nil {
		httpInternalServerErrorResponse(c, err.Error())
		return
	}

	for _, val := range uniqueNum(form.Roles) {
		userRole := models.UserRole{
			UserID:    user.ID,
			RoleID:    val,
			GrantDate: time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = db.GetConnection().Insert(&userRole)
		if err != nil {
			httpInternalServerErrorResponse(c, err.Error())
			return
		}
	}

	result := map[string]interface{}{
		"message": "Register successfully",
	}

	httpOkResponse(c, result)
}

func signin(c *gin.Context) {
	var user models.User

	form := &struct {
		Email    string `form:"email" json:"email"`
		Password string `form:"password" json:"password"`
	}{}
	c.BindJSON(form)

	// form validation
	err := validation.Errors{
		"email":    validation.Validate(form.Email, validation.Required, is.Email),
		"password": validation.Validate(form.Password, validation.Required),
	}.Filter()

	err = db.GetConnection().Model(&user).Where("email = ?", form.Email).Select()
	if err != nil {
		httpForbiddenErrorResponse(c, "Invalid email or password.")
		return
		// httpInternalServerErrorResponse(c, err.Error())
		// return
	}

	if form.Password != user.Password {
		httpForbiddenErrorResponse(c, "Invalid email or password.")
		return
	}

	result := map[string]interface{}{
		"message": "Login successfully",
	}

	httpOkResponse(c, result)

}
