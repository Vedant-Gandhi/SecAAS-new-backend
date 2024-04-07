package user

import (
	"net/http"
	"secaas_backend/model"
	"secaas_backend/svc/errors"
	"secaas_backend/svc/user"
	"secaas_backend/transport/controller/response"

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
