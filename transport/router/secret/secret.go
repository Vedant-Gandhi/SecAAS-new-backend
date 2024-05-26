package secret

import (
	"secaas_backend/transport/controller/secret"

	"github.com/gin-gonic/gin"
)

func Add(router *gin.RouterGroup, controller secret.SecretsController) {

	secret := router.Group("/secrets")

	secret.POST("", controller.Create())

	secret.GET("/organization/:organizationId/user/:userId", controller.GetForUserOrganization())
	secret.GET("/:secretId/organization/:organizationId/users", controller.GetUsersForSecret())
	secret.GET("/organization/:organizationId", controller.GetForOrganization())

	secret.POST("/:secretId/share", controller.ShareKey())

}
