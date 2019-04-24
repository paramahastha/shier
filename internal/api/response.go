package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func httpOkResponse(c *gin.Context, payload map[string]interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
		"data":   payload,
	})
}

func httpValidationErrorResponse(c *gin.Context, validationErrors string) {
	c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status": http.StatusBadRequest,
		"error": map[string]interface{}{
			"message": "Validation error",
			"errors":  validationErrors,
		},
	})
}

func httpInternalServerErrorResponse(c *gin.Context, errorMessage string) {
	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"status": http.StatusInternalServerError,
		"error": map[string]interface{}{
			"message": errorMessage,
		},
	})
}
