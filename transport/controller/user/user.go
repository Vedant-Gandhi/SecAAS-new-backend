package user

import (
	"math"
	"net/http"
	"secaas_backend/model"
	"secaas_backend/svc/errors"
	"secaas_backend/svc/user"
	"secaas_backend/transport/controller/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	logger *logrus.Logger
	svc    *user.UserSVC
}

func New(svc *user.UserSVC, logger *logrus.Logger) *UserController {
	uc := &UserController{logger: logger, svc: svc}
	return uc
}

func (u *UserController) GetUserByEmail() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		query := gCtx.Request.URL.Query()

		if !query.Has("email") {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "user/invalid-email",
				Message: "User Email is not valid",
			})
			return
		}

		email := query.Get("email")

		user, err := u.svc.GetByEmail(gCtx.Request.Context(), model.Email(email))

		if err != nil {
			if err == errors.ErrUserNotFound {
				gCtx.JSON(http.StatusNotFound, response.ErrorResponse{
					Code:    "user/not-found",
					Message: "User Not Found",
				})
				return
			}

			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		gCtx.JSON(http.StatusOK, user)

	}
}

func (u *UserController) CreateUser() gin.HandlerFunc {
	return func(gCtx *gin.Context) {

		var user model.User

		err := gCtx.BindJSON(&user)

		if err != nil {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "data/invalid-payload",
				Message: "Payload format is not valid",
			})
			u.logger.WithError(err).Error("error in decoding body")
			return
		}

		if user.Email == "" {
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "user/invalid-email",
				Message: "User Email is not valid",
			})
			u.logger.WithError(err).Error("invalid email when creating user")
			return
		}

		if user.PassHash.Hash == "" || user.PassHash.Alg == "" {
			u.logger.WithError(err).Error("invalid fields in pass hash found")
			gCtx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Code:    "user/invalid-password",
				Message: "User Password is not valid",
			})
			return
		}

		newUser, err := u.svc.CreateUser(gCtx.Request.Context(), user)

		if err != nil {
			if err == errors.ErrInvalidEmail {
				gCtx.JSON(http.StatusNotFound, response.ErrorResponse{
					Code:    "user/invalid-email",
					Message: "User Email is not Valid.",
				})
				return
			}

			gCtx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Code:    "server/internal-error",
				Message: "An Internal Server error has occurred",
			})
			return
		}

		gCtx.JSON(http.StatusCreated, newUser)

	}
}

func (s *UserController) GetUsersForOrganization() gin.HandlerFunc {
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

		data, err := s.svc.GetUsersByOrganization(gCtx.Request.Context(), model.OrganizationID(orgId), pageParams)

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

		if len(data) == 0 {
			resp.Data = []model.User{}
		}
		gCtx.JSON(http.StatusOK, resp)

	}
}
