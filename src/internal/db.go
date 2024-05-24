package internal

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDatabase interface {
	Driver() *mongo.Database
}

type database struct {
	driver *mongo.Client
}

func (db *database) Driver() *mongo.Database {
	return db.driver.Database(Config.DbName)
}

func InitDb() {
	driver, err := mongo.Connect(context.Background(), &options.ClientOptions{})
	if err != nil {
		log.Errorf("error opening db")
		panic(err)
	}
	db := &database{driver: driver}
	Database = db
}

var Database IDatabase
