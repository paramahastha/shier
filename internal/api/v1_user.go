package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"shier/internal/models"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-redis/redis"
)

func (s *Server) getAllUsers(c *gin.Context) {
	var users []models.User
	err := s.DB.Model(&users).Column("Roles").Select()

	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
	}

	result := map[string]interface{}{
		"data": users,
	}

	httpSuccessResponse(c, result["data"], http.StatusOK, "-")
}

func (s *Server) createUser(c *gin.Context) {
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

	httpSuccessResponse(c, nil, http.StatusCreated, "Create user successfully")
}

func (s *Server) getUserById(c *gin.Context) {
	var user models.User

	id := c.Param("id")

	val, err := s.Redis.Get(fmt.Sprintf("user_%s", id)).Result()
	if err == redis.Nil || val == "" {
		// get from database
		err := s.DB.Model(&user).Column("Roles").Where("id = ?", id).Select()
		if err != nil {
			httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
			return
		}

		json, err := json.Marshal(user)

		err = s.Redis.Set(fmt.Sprintf("user_%s", id), json, 0).Err()
		if err != nil {
			httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
			return
		}

		result := map[string]interface{}{
			"data": user,
		}

		httpSuccessResponse(c, result["data"], http.StatusOK, "-")
	} else if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	} else {
		byt := []byte(val)

		if err := json.Unmarshal(byt, &user); err != nil {
			httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
			return
		}

		result := map[string]interface{}{
			"data": user,
		}
		httpSuccessResponse(c, result["data"], http.StatusOK, "-")
	}

}

func (s *Server) updateUserById(c *gin.Context) {
	var user models.User
	var role models.Role
	var userRoles []models.UserRole

	form := &struct {
		FirstName string `form:"first_name" json:"first_name"`
		LastName  string `form:"last_name" json:"last_name"`
		Email     string `form:"email" json:"email"`
		Password  string `form:"password" json:"password"`
		Confirm   string `form:"confirm" json:"confirm"`
		Roles     []int  `form:"roles" json:"roles"`
	}{}
	id := c.Param("id")
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

	err = s.DB.Model(&user).Where("id = ?", id).Select()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	user = models.User{
		ID:        user.ID,
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Password:  form.Password,
		UpdatedAt: time.Now(),
	}

	_, err = s.DB.Model(&user).
		Column("first_name").
		Column("last_name").
		Column("email").
		Column("password").
		Column("updated_at").
		WherePK().Returning("*").Update()

	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	err = s.DB.Model(&userRoles).Where("user_id = ?", user.ID).Select()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusBadRequest, "-")
		return
	}

	for _, item := range userRoles {
		err = s.DB.Model(&item).Where("user_id = ?", user.ID).
			Where("role_id = ?", item.RoleID).Select()

		if err != nil {
			httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
			return
		}

		err = s.DB.Delete(&item)
		if err != nil {
			httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
			return
		}
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

	err = s.DB.Model(&user).Column("Roles").Where("id = ?", id).Select()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	// Update user cache
	json, err := json.Marshal(user)
	err = s.Redis.Set(fmt.Sprintf("user_%s", id), json, 0).Err()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	httpSuccessResponse(c, nil, http.StatusOK, "Update user successfully")
}

func (s *Server) deleteUserById(c *gin.Context) {
	var user models.User
	var userRoles []models.UserRole

	id := c.Param("id")
	err := validation.Errors{
		"id": validation.Validate(id, validation.Required),
	}.Filter()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusBadRequest, "-")
		return
	}

	err = s.DB.Model(&user).Column("Roles").Where("id = ?", id).Select()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	err = s.DB.Model(&userRoles).Where("user_id = ?", user.ID).Select()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	for _, item := range userRoles {
		err = s.DB.Model(&item).Where("user_id = ?", user.ID).
			Where("role_id = ?", item.RoleID).Select()

		if err != nil {
			httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
			return
		}

		err = s.DB.Delete(&item)
		if err != nil {
			httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
			return
		}
	}

	err = s.DB.Delete(&user)
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	err = s.Redis.Set(fmt.Sprintf("user_%s", id), "", 0).Err()
	if err != nil {
		httpErrorResponse(c, err.Error(), http.StatusInternalServerError, "-")
		return
	}

	httpSuccessResponse(c, nil, http.StatusOK, "Delete user successfully")
}
