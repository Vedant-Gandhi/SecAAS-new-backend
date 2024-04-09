package organization

import (
	"secaas_backend/transport/controller/organization"

	"github.com/gin-gonic/gin"
)

func Add(router *gin.RouterGroup, controller organization.OrganizationController) {

	organization := router.Group("/organizations")

	organization.POST("/", controller.CreateOrganization())
	organization.DELETE("/:organizationId", controller.DeleteOrganization())

}
