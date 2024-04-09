package organization

import (
	"secaas_backend/svc/organization"

	"github.com/sirupsen/logrus"
)

type OrganizationController struct {
	logger *logrus.Logger
	svc    *organization.OrganizationSVC
}

func New(svc *organization.OrganizationSVC, logger *logrus.Logger) *OrganizationController {
	uc := &OrganizationController{logger: logger, svc: svc}
	return uc
}
