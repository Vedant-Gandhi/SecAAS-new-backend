package organization

import (
	"fmt"
	"net/http"
	"secaas_backend/model"
	"secaas_backend/svc/errors"
	"secaas_backend/svc/organization"
	"secaas_backend/svc/user"
	"secaas_backend/transport/controller/response"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type createOrgRequest struct {
	model.Organization
	AdminPvtKey string `json:"adminPvtKey"`
}

type OrganizationController struct {
	logger  *logrus.Logger
	svc     *organization.OrganizationSVC
	userSvc *user.UserSVC
}

func New(svc *organization.OrganizationSVC, userSvc *user.UserSVC, logger *logrus.Logger) *OrganizationController {
	uc := &OrganizationController{logger: logger, svc: svc, userSvc: userSvc}
	return uc
}

func (u *OrganizationController) CreateOrganization() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		fmt.Print("Received request in organization create")

		var orgReq model.Organization

		err := gCtx.BindJSON(&orgReq)

		if err != nil {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "data/invalid-payload",
				Message: "Payload format is not valid",
			})
			u.logger.WithError(err).Error("error in decoding body")
			return
		}

		if orgReq.AdminEmail == "" {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "organization/invalid-email",
				Message: "Organization Email is not valid",
			})
			u.logger.WithError(err).Error("invalid admin email when creating organization")
			return
		}

		organization := model.Organization{
			Name:         orgReq.Name,
			BillingEmail: orgReq.BillingEmail,
			AdminEmail:   orgReq.AdminEmail,
			SymmKey:      orgReq.SymmKey,
		}

		org, err := u.svc.CreateNew(gCtx.Request.Context(), organization)

		if err != nil {
			if err == errors.ErrInvalidAsymmetricKey {
				gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
					Code:    "security/invalid-asymm-key",
					Message: "Asymmetric key is not valid",
				})

				return
			}

			if err == errors.ErrInvalidEmail {
				gCtx.JSON(http.StatusNotFound, response.ErrorResponse{
					Code:    "organization/invalid-email",
					Message: "organization Admin Email is not Valid.",
				})
				return
			}

			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		fmt.Printf("%v", org)
		gCtx.JSON(http.StatusCreated, org)

	}
}

func (o *OrganizationController) DeleteOrganization() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		organizationId := gCtx.Param("organizationId")

		if organizationId == "" {
			err := response.ErrorResponse{
				Code:    "organization/invalid-id",
				Message: "Organization ID is not valid",
			}
			gCtx.JSON(http.StatusBadRequest, err)
			return
		}

		deleteCount, err := o.svc.DeleteOrganization(gCtx.Request.Context(), model.OrganizationID(organizationId))

		if err != nil {
			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		if deleteCount > 0 {
			organizationId = ""
		}

		gCtx.JSON(http.StatusOK, gin.H{
			"deleted": deleteCount > 0,
			"id":      organizationId,
		})

	}
}

func (u *OrganizationController) GetOrganizationsForUser() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		userId := gCtx.Param("userId")

		user, err := u.userSvc.GetByID(gCtx.Request.Context(), model.UserID(userId))

		if err != nil {
			if err == errors.ErrUserNotFound {
				gCtx.JSON(http.StatusNotFound, response.ErrorResponse{
					Code:    "user/not-found",
					Message: "User not found",
				})
				return
			}
			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		organizations := []model.OrganizationID{}

		for _, org := range user.Organization {
			if org.ID != "" {
				organizations = append(organizations, model.OrganizationID(org.ID))
			}
		}

		data, err := u.svc.GetOrganizationList(gCtx.Request.Context(), organizations)

		if err != nil {
			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		gCtx.JSON(http.StatusOK, data)

	}
}
