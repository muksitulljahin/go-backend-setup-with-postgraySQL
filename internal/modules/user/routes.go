package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	controller := NewController(service)

	users := rg.Group("/users")
	{
		users.POST("", controller.CreateUser)
		users.GET("", controller.GetAllUsers)
		users.GET("/:id", controller.GetUserByID)
		users.PUT("/:id", controller.UpdateUser)
		users.DELETE("/:id", controller.DeleteUser)
	}
}
