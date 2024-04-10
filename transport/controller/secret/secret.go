package secret

import (
	"secaas_backend/svc/secret"

	"github.com/sirupsen/logrus"
)

type SecretsController struct {
	logger *logrus.Logger
	svc    *secret.SecretsSVC
}

func New(svc *secret.SecretsSVC, logger *logrus.Logger) *SecretsController {
	uc := &SecretsController{logger: logger, svc: svc}
	return uc
}
