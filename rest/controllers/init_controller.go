package controllers

import (
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/nats-io/nats.go"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var ec *nats.EncodedConn
var esc *elastic.Client
var rc *redis.Client
var log *logrus.Logger
var db *gorm.DB
var mdb *mongo.Client
var collection *mongo.Collection

func InitControllers(
	Log *logrus.Logger,
	Db *gorm.DB,
	Mdb *mongo.Client,
	Echo *echo.Echo,
	Esc *elastic.Client,
	Rc *redis.Client,
) {

	log = Log
	esc = Esc
	rc = Rc
	db = Db
	mdb = Mdb

	e := Echo

	collection = mdb.Database("taskpedia").Collection("transactions")

	e.GET("/user/view", GetUser)
	e.POST("/user/create", InsertUser)
	e.POST("/user/update", UpdateUser)

	e.GET("/task/view", GetTask)
	e.POST("/task/create", CreateTask)
	e.POST("/task/update", UpdateTask)

	e.GET("/task/search/userid", SearchTaskByUserID)

	e.GET("/stat", GetDataStats)
	e.GET("/log", GetTransactionsLog)
}

func InitNatsConn(
	NatsConn *nats.Conn,
	EnConn *nats.EncodedConn) {
	ec = EnConn
}
