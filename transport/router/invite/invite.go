package invite

import (
	"secaas_backend/transport/controller/invite"

	"github.com/gin-gonic/gin"
)

func Add(router *gin.RouterGroup, controller invite.InviteController) {

	invite := router.Group("/invites")

	invite.POST("/", controller.SendInvite())
	invite.DELETE("/:inviteId", controller.DeleteInvite())

	invite.GET("/organization/:organizationId", controller.GetByOrganization())
	invite.GET("/user/email", controller.GetForUser())

}
