package svc

import (
	"secaas_backend/db"
	"secaas_backend/svc/user"

	"github.com/sirupsen/logrus"
)

type SVC struct {
	logger *logrus.Logger
	db     *db.DB

	User *user.UserSVC
}

func New(logger *logrus.Logger, db *db.DB) *SVC {
	u := user.New(logger)

	s := &SVC{logger: logger, db: db, User: u}
	return s
}
