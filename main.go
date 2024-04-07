package main

import (
	"context"
	"net/http"
	"secaas_backend/db"
	"secaas_backend/svc"
	"secaas_backend/transport/controller"
	"secaas_backend/transport/router"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func main() {
	ctx := context.Background()

	logger := logrus.New()

	viper.SetEnvPrefix("SECAAS")

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	mongoCon := viper.GetString("SECAAS_MONGO_CONNECTION_STRING")
	mongoDb := viper.GetString("SECAAS_MONGO_DATABASE")

	db, err := db.New(db.MongoCfg{
		ConnectionString: mongoCon,
		Database:         mongoDb,
	}, logger)

	if err != nil {
		logger.WithError(err).Error("failed to connect to mongodb")
		return
	}

	svc := svc.New(logger, db)

	controller := controller.New(logger, svc)

	httpRouter, err := router.Init(logger, controller)

	if err != nil {
		logger.WithContext(ctx).WithError(err).Error("error while initialising the router.")
		return
	}

	err = viper.BindEnv("WEBSERVER_ADDRESS")
	if err != nil {
		logger.WithContext(ctx).WithError(err).Error("failed to bind web server address in env.")
		return
	}

	httpAddr := viper.GetString("SECAAS_WEBSERVER_ADDRESS")

	if !viper.IsSet("SECAAS_WEBSERVER_ADDRESS") {
		logger.WithContext(ctx).Error("cannot listen to the address.")
		return
	}

	httpRouter.Router.Run(httpAddr)
}
