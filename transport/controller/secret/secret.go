package secret

import (
	"math"
	"net/http"
	"secaas_backend/model"
	"secaas_backend/svc/errors"
	"secaas_backend/svc/secret"
	"secaas_backend/transport/controller/response"
	"strconv"

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

func (s *SecretsController) Create() gin.HandlerFunc {
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

		newSecret, err := s.svc.Create(gCtx.Request.Context(), secret)

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

func (s *SecretsController) GetForUserOrganization() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		orgId := gCtx.Param("organizationId")
		userId := gCtx.Param("userId")

		rawPage := gCtx.Query("page")
		rawLimit := gCtx.Query("limit")
		page, err := strconv.Atoi(rawPage)

		if err != nil || page < 0 {
			page = 1
		}

		limit, err := strconv.Atoi(rawLimit)

		if err != nil || limit < 0 || limit > 100 {
			limit = 10
		}

		pageParams := model.PaginationParams{
			Page:  page,
			Limit: limit,
			Skip:  int(math.Max(float64(page-1), 0)) * limit,
		}

		data, err := s.svc.GetAllSecretsforaUserInOrganization(gCtx.Request.Context(), model.UserID(userId), model.OrganizationID(orgId), pageParams)

		if err != nil {
			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		resp := model.PaginationResponse{
			CurrentPage: page,
			Data:        data,
			Limit:       limit,
			NextPage:    page + 1,
		}
		gCtx.JSON(http.StatusOK, resp)

	}
}

func (s *SecretsController) GetForOrganization() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		orgId := gCtx.Param("organizationId")

		rawPage := gCtx.Query("page")
		rawLimit := gCtx.Query("limit")
		page, err := strconv.Atoi(rawPage)

		if err != nil || page < 0 {
			page = 1
		}

		limit, err := strconv.Atoi(rawLimit)

		if err != nil || limit < 0 || limit > 100 {
			limit = 10
		}

		pageParams := model.PaginationParams{
			Page:  page,
			Limit: limit,
			Skip:  int(math.Max(float64(page-1), 0)) * limit,
		}

		data, err := s.svc.GetAllSecretsforOrganization(gCtx.Request.Context(), model.OrganizationID(orgId), pageParams)

		if err != nil {
			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		resp := model.PaginationResponse{
			CurrentPage: page,
			Data:        data,
			Limit:       limit,
			NextPage:    page + 1,
		}
		gCtx.JSON(http.StatusOK, resp)

	}
}
