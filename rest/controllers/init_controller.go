package controllers

import (
	"github.com/labstack/echo"
	"github.com/nats-io/nats.go"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var ec *nats.EncodedConn
var esc *elastic.Client
var log *logrus.Logger
var db *gorm.DB

func InitControllers(
	Log *logrus.Logger,
	Database *gorm.DB,
	Echo *echo.Echo,
	Esc *elastic.Client,
) {

	log = Log
	esc = Esc
	db = Database

	e := Echo

	e.GET("/user/view", GetUser)
	e.POST("/user/create", InsertUser)
	e.POST("/user/update", UpdateUser)

	e.GET("/task/view", GetTask)
	e.POST("/task/create", CreateTask)
	e.POST("/task/update", UpdateTask)

	e.GET("/task/search", SearchTaskByUser)
}

func InitNatsConn(
	NatsConn *nats.Conn,
	EnConn *nats.EncodedConn) {
	ec = EnConn
}
