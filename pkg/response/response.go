package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func JSON(c *gin.Context, statusCode int, success bool, message string, data interface{}, err interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: success,
		Message: message,
		Data:    data,
		Error:   err,
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusOK, true, message, data, nil)
}

func Created(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusCreated, true, message, data, nil)
}

func BadRequest(c *gin.Context, message string, err interface{}) {
	JSON(c, http.StatusBadRequest, false, message, nil, err)
}

func NotFound(c *gin.Context, message string) {
	JSON(c, http.StatusNotFound, false, message, nil, nil)
}

func InternalServerError(c *gin.Context, message string, err interface{}) {
	JSON(c, http.StatusInternalServerError, false, message, nil, err)
}
