package controller

import (
	"secaas_backend/svc"
	"secaas_backend/transport/controller/invite"
	"secaas_backend/transport/controller/organization"
	"secaas_backend/transport/controller/secret"
	"secaas_backend/transport/controller/user"

	"github.com/sirupsen/logrus"
)

type Controller struct {
	logger *logrus.Logger
	svc    *svc.SVC

	User         *user.UserController
	Invite       *invite.InviteController
	Organization *organization.OrganizationController
	Secrets      *secret.SecretsController
}

func New(logger *logrus.Logger, svc *svc.SVC) *Controller {
	u := user.New(svc.User, logger)
	i := invite.New(svc.Invite, logger)
	o := organization.New(svc.Organization, svc.User, logger)
	sec := secret.New(svc.Secrets, logger)

	c := &Controller{logger: logger, svc: svc, User: u, Invite: i, Organization: o, Secrets: sec}
	return c
}
