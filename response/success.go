package response

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendJSON sends formatted API Response
func SendJSON(c *gin.Context, statusCode int, success bool, message string, data interface{}, meta *PaginationMeta, err interface{}) {
	c.JSON(statusCode, APIResponse{
		Success:    success,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta:       meta,
		Error:      err,
	})
}

// Success sends 200 OK with data payload
func Success(c *gin.Context, message string, data interface{}) {
	SendJSON(c, http.StatusOK, true, message, data, nil, nil)
}

// Created sends 201 Created with created object payload
func Created(c *gin.Context, message string, data interface{}) {
	SendJSON(c, http.StatusCreated, true, message, data, nil, nil)
}

// Paginated sends 200 OK with list data and pagination meta (page, limit, total, totalPage)
func Paginated(c *gin.Context, message string, data interface{}, page int, limit int, total int64) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	meta := &PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}

	SendJSON(c, http.StatusOK, true, message, data, meta, nil)
}
