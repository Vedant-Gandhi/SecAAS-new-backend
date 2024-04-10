package secret

import (
	"net/http"
	"secaas_backend/model"
	"secaas_backend/svc/errors"
	"secaas_backend/svc/secret"
	"secaas_backend/transport/controller/response"

	"github.com/gin-gonic/gin"
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

func (s *SecretsController) CreateUser() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		var secret model.Secret

		err := gCtx.BindJSON(&secret)

		if err != nil {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "data/invalid-payload",
				Message: "Payload format is not valid",
			})
			s.logger.WithError(err).Error("error in decoding body")
			return
		}

		if secret.EncryptedData == "" {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "secret/invalid-payload",
				Message: "Secret Payload is not valid",
			})
			s.logger.WithError(err).Error("secret payload is not valid")
			return
		}

		if secret.CreatorEmail == "" {
			s.logger.WithError(err).Error("invalid secret creator key")
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "secret/invalid-email",
				Message: "User Email is not valid",
			})
			return
		}

		newSecret, err := s.svc.CreateForUser(gCtx.Request.Context(), secret)

		if err != nil {
			if err == errors.ErrInvalidEmail {
				gCtx.JSON(http.StatusNotFound, response.ErrorResponse{
					Code:    "secret/invalid-email",
					Message: "Secret Email is not Valid.",
				})
				return
			}

			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		gCtx.JSON(http.StatusCreated, newSecret)

	}
}
