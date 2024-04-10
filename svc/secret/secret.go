package secret

import (
	"context"
	"secaas_backend/model"
	"secaas_backend/svc/errors"

	"github.com/sirupsen/logrus"
)

type SecretsSVC struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *SecretsSVC {
	u := &SecretsSVC{logger: logger}
	return u
}
func (s *SecretsSVC) GetForUser(ctx context.Context, userId model.UserID, organizationId string) (sec model.Secret, err error) {

	if userId == "" {
		s.logger.WithContext(ctx).Error("invalid user id to get secret list for user")
		err = errors.ErrInvalidID
		return
	}

	if organizationId == "" {
		s.logger.WithContext(ctx).Error("invalid organization id to get secret list for user")
		err = errors.ErrInvalidOrganizationID
		return
	}

	return
}
