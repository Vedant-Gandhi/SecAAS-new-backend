package svc

import (
	"secaas_backend/db"
	"secaas_backend/svc/invite"
	"secaas_backend/svc/user"

	"github.com/sirupsen/logrus"
)

type SVC struct {
	logger *logrus.Logger
	db     *db.DB

	User   *user.UserSVC
	Invite *invite.InviteSVC
}

func New(logger *logrus.Logger, db *db.DB) *SVC {
	u := user.New(logger)
	i := invite.New(logger)

	s := &SVC{logger: logger, db: db, User: u, Invite: i}
	return s
}
