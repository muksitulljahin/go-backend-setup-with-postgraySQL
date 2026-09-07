package user

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/pkg/response"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	user, err := ctrl.service.CreateUser(&req)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to create user", err.Error())
		return
	}

	response.Created(c, "User created successfully", user)
}

func (ctrl *Controller) GetAllUsers(c *gin.Context) {
	users, err := ctrl.service.GetAllUsers()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve users", err.Error())
		return
	}

	response.Success(c, "Users retrieved successfully", users)
}

func (ctrl *Controller) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err.Error())
		return
	}

	user, err := ctrl.service.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to retrieve user", err.Error())
		return
	}

	response.Success(c, "User retrieved successfully", user)
}

func (ctrl *Controller) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err.Error())
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	user, err := ctrl.service.UpdateUser(uint(id), &req)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		if errors.Is(err, ErrUserAlreadyExists) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to update user", err.Error())
		return
	}

	response.Success(c, "User updated successfully", user)
}

func (ctrl *Controller) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err.Error())
		return
	}

	err = ctrl.service.DeleteUser(uint(id))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to delete user", err.Error())
		return
	}

	response.Success(c, "User deleted successfully", nil)
}
