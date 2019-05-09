package api

import (
	"github.com/gin-gonic/gin"
)

func httpSuccessResponse(c *gin.Context, payload interface{}, code int, message string) {
	c.JSON(code, map[string]interface{}{
		"data":    payload,
		"code":    code,
		"message": message,
	})
}

func httpErrorResponse(c *gin.Context, errorPayload interface{}, code int, message string) {
	c.JSON(code, map[string]interface{}{
		"errors":  errorPayload,
		"code":    code,
		"message": message,
	})
}
