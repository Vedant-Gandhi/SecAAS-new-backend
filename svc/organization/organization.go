package organization

import (
	"github.com/sirupsen/logrus"
)

type OrganizationSVC struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *OrganizationSVC {
	u := &OrganizationSVC{logger: logger}
	return u
}
