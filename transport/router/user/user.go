package user

import (
	"secaas_backend/transport/controller/user"

	"github.com/gin-gonic/gin"
)

func Add(router *gin.RouterGroup, controller user.UserController) {

	user := router.Group("/users")

	user.GET("/by/email", controller.GetUserByEmail())
	user.POST("/create", controller.CreateUser())
}
