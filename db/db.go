package db

import (
	"github.com/kamva/mgm/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
}

type MongoCfg struct {
	ConnectionString string
	Database         string
}

func New(cfg MongoCfg, logger *logrus.Logger) (*DB, error) {
	db := &DB{}
	err := mgm.SetDefaultConfig(nil, cfg.Database, options.Client().ApplyURI(cfg.ConnectionString))

	logger.Debug("Connected to MongoDB")
	return db, err
}
