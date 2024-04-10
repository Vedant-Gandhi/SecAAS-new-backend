package secret

import (
	"secaas_backend/transport/controller/secret"

	"github.com/gin-gonic/gin"
)

func Add(router *gin.RouterGroup, controller secret.SecretsController) {

	secret := router.Group("/secrets")

	secret.POST("", controller.CreateUser())

}
