package invite

import (
	"math"
	"net/http"
	"secaas_backend/model"
	"secaas_backend/svc/errors"
	"secaas_backend/svc/invite"
	"secaas_backend/transport/controller/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type acceptInviteReq struct {
}

type InviteController struct {
	logger *logrus.Logger
	svc    *invite.InviteSVC
}

func New(svc *invite.InviteSVC, logger *logrus.Logger) *InviteController {
	uc := &InviteController{logger: logger, svc: svc}
	return uc
}

func (i *InviteController) SendInvite() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		var invite model.Invite

		err := gCtx.BindJSON(&invite)

		if err != nil {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "data/invalid-payload",
				Message: "Payload format is not valid",
			})
			i.logger.WithError(err).Error("error in decoding body")
			return
		}

		if invite.FromUserEmail == "" {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "user/invalid-email",
				Message: "Sender User Email is not valid",
			})
			i.logger.WithError(err).Error("invalid sender email when creating invite")
			return
		}

		if invite.OrganizationID == "" {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "organization/invalid-id",
				Message: "Organizaton ID is not valid",
			})
			i.logger.WithError(err).Error("invalid organization id when creating invite")
			return
		}

		if invite.ToUserEmail == "" {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "user/invalid-email",
				Message: "Receiver User Email is not valid",
			})
			i.logger.WithError(err).Error("invalid receiver email when creating invite")
			return
		}

		if invite.ExpiresAt.IsZero() {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "invite/invalid-expiry",
				Message: "Invite Expiry is not valid",
			})
			i.logger.WithError(err).Error("invalid invite expiry when creating invite")
			return
		}

		newInvite, err := i.svc.CreateInvite(gCtx.Request.Context(), invite)

		if err != nil {

			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		gCtx.JSON(http.StatusCreated, newInvite)

	}
}

func (i *InviteController) GetByOrganization() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		orgId := gCtx.Param("organizationId")

		if orgId == "" {
			err := response.ErrorResponse{
				Code:    "organization/invalid-id",
				Message: "Organization ID is not valid",
			}
			gCtx.JSON(http.StatusBadRequest, err)
			return
		}

		getOnlyActive := gCtx.Query("active") == "true"
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

		inviteList, err := i.svc.GetInvitesByOrganization(gCtx.Request.Context(), orgId, pageParams, getOnlyActive)

		if err != nil {
			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		pageResp := model.PaginationResponse{
			CurrentPage: page,
			Data:        inviteList,
			Limit:       limit,
			NextPage:    page + 1,
		}

		gCtx.JSON(http.StatusOK, pageResp)

	}
}

func (i *InviteController) GetForUser() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		email := gCtx.Query("email")

		if email == "" {
			err := response.ErrorResponse{
				Code:    "user/invalid-email",
				Message: "Email ID is not valid",
			}
			gCtx.JSON(http.StatusBadRequest, err)
			return
		}

		getOnlyActive := gCtx.Query("active") == "true"
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

		inviteList, err := i.svc.GetInvitesForUser(gCtx.Request.Context(), email, pageParams, getOnlyActive)

		if err != nil {
			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		pageResp := model.PaginationResponse{
			CurrentPage: page,
			Data:        inviteList,
			Limit:       limit,
			NextPage:    page + 1,
		}

		gCtx.JSON(http.StatusOK, pageResp)

	}
}

func (i *InviteController) DeleteInvite() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		inviteId := gCtx.Param("inviteId")

		if inviteId == "" {
			err := response.ErrorResponse{
				Code:    "invite/invalid-id",
				Message: "Invite ID is not valid",
			}
			gCtx.JSON(http.StatusBadRequest, err)
			return
		}

		deleteCount, err := i.svc.DeleteInvite(gCtx.Request.Context(), inviteId)

		if err != nil {
			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		if deleteCount > 0 {
			inviteId = ""
		}

		gCtx.JSON(http.StatusOK, gin.H{
			"deleted": deleteCount > 0,
			"id":      inviteId,
		})

	}
}

func (i *InviteController) AcceptInvite() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		inviteId := gCtx.Param("inviteId")

		if inviteId == "" {
			err := response.ErrorResponse{
				Code:    "invite/invalid-id",
				Message: "Invite ID is not valid",
			}
			gCtx.JSON(http.StatusBadRequest, err)
			return
		}

		err := i.svc.AcceptInvite(gCtx.Request.Context(), inviteId)

		if err != nil {
			if err == errors.ErrInviteNotFound {
				err := response.ErrorResponse{
					Code:    "invite/not-found",
					Message: "Invite not found.",
				}
				gCtx.JSON(http.StatusBadRequest, err)
				return
			}

			if err == errors.ErrInvalidID {
				err := response.ErrorResponse{
					Code:    "invite/invalid-id",
					Message: "Invite ID is not valid.",
				}
				gCtx.JSON(http.StatusBadRequest, err)
				return
			}

			err := response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "Internal server error has ocurred.",
			}
			gCtx.JSON(http.StatusInternalServerError, err)
			return

		}

		gCtx.JSON(http.StatusOK, gin.H{})

	}
}
