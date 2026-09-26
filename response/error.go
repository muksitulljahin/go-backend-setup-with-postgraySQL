package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError is a custom error struct for Node.js style AppError handling
type AppError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func NewAppError(statusCode int, message string, err ...error) *AppError {
	var originalErr error
	if len(err) > 0 {
		originalErr = err[0]
	}
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Err:        originalErr,
	}
}

// Custom Error Response
func Error(c *gin.Context, statusCode int, message string, errDetail interface{}) {
	SendJSON(c, statusCode, false, message, nil, nil, errDetail)
}

// BadRequest (400)
func BadRequest(c *gin.Context, message string, errDetail interface{}) {
	Error(c, http.StatusBadRequest, message, errDetail)
}

// Unauthorized (401)
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden (403)
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

// NotFound (404)
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

// InternalServerError (500)
func InternalServerError(c *gin.Context, message string, errDetail interface{}) {
	Error(c, http.StatusInternalServerError, message, errDetail)
}

// HandleError automatically handles AppError or standard errors (Node.js errorHandler style)
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		var detail interface{}
		if appErr.Err != nil {
			detail = appErr.Err.Error()
		}
		Error(c, appErr.StatusCode, appErr.Message, detail)
		return
	}
	InternalServerError(c, "Something went wrong", err.Error())
}
