package controller

import (
	"secaas_backend/svc"
	"secaas_backend/transport/controller/invite"
	"secaas_backend/transport/controller/user"

	"github.com/sirupsen/logrus"
)

type Controller struct {
	logger *logrus.Logger
	svc    *svc.SVC

	User   *user.UserController
	Invite *invite.InviteController
}

func New(logger *logrus.Logger, svc *svc.SVC) *Controller {
	u := user.New(svc.User, logger)
	i := invite.New(svc.Invite, logger)
	c := &Controller{logger: logger, svc: svc, User: u, Invite: i}
	return c
}
