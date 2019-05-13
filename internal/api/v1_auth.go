package api

import (
	"fmt"
	"net/http"
	"time"

	"shier/internal/models"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (s *Server) signup(c *gin.Context) {
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
		httpErrorResponse(c, err.Error(), http.StatusBadRequest, "-")
		return
	}

	// Check existing role in db
	for _, val := range uniqueNum(form.Roles) {

		err = s.DB.Model(&role).Where("id = ?", val).Select()
		if err != nil {
			httpErrorResponse(c, err.Error(), http.StatusBadRequest, fmt.Sprintf("pg: no role (%d) in result set", val))
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

	err = s.DB.Insert(&user)
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
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

		err = s.DB.Insert(&userRole)
		if err != nil {
			httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
			return
		}
	}

	httpSuccessResponse(c, nil, http.StatusCreated, "Register successfully")
}

func (s *Server) signin(c *gin.Context) {
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

	err = s.DB.Model(&user).Where("email = ?", form.Email).Select()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusForbidden, "Invalid email or password.")
		return
	}

	if form.Password != user.Password {
		httpErrorResponse(c, err.Error(), http.StatusForbidden, "Invalid email or password.")
		return
	}

	httpSuccessResponse(c, nil, http.StatusOK, "Login successfully")

}
