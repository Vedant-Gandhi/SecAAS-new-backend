package router

import (
	"secaas_backend/transport/controller"
	"secaas_backend/transport/router/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type httpRouter struct {
	Router     *gin.Engine
	logger     *logrus.Logger
	controller *controller.Controller
}

func Init(logger *logrus.Logger, c *controller.Controller) (*httpRouter, error) {
	gr := gin.Default()

	apiV1 := gr.Group("/api/v1")

	user.Add(apiV1, *c.User)

	r := &httpRouter{logger: logger, Router: gr, controller: c}

	return r, nil
}
