package controller

import (
	"secaas_backend/svc"
	"secaas_backend/transport/controller/user"

	"github.com/sirupsen/logrus"
)

type Controller struct {
	logger *logrus.Logger
	svc    *svc.SVC

	User *user.UserController
}

func New(logger *logrus.Logger, svc *svc.SVC) *Controller {
	u := user.New(svc.User, logger)
	c := &Controller{logger: logger, svc: svc, User: u}
	return c
}
