package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"

	"github.com/go-redis/redis"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"taskpedia-rest/controllers"
	"taskpedia-rest/models"
)

var euri = os.Getenv("ELASTIC_URI")
var ruri = os.Getenv("REDIS_URI")
var nuri = os.Getenv("NATS_URI")
var muri = os.Getenv("MONGO_URI")

var ctx = context.Background()

var db = InitDB()
var mdb = InitMongo()
var e = echo.New()

var log = InitLogger()
var esc = InitSearchClient()
var rc = InitRedisClient()

func main() {
	log := InitLogger()

	controllers.InitControllers(log, db, mdb, e, esc, rc)

	var nc *nats.Conn

	for i := 0; i < 5; i++ {
		nConn, err := nats.Connect(nuri)
		if err == nil {
			nc = nConn
			break
		}

		fmt.Println("Waiting before connecting to NATS at:", nuri)
		time.Sleep(1 * time.Second)
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Panic(err)
	}
	defer ec.Close()

	controllers.InitNatsConn(nc, ec)

	log.Info("Connected to NATS and ready to publish messages")

	errWeb := e.Start(":18080")

	if errWeb != nil {
		log.Fatal(errWeb)
	} else {
		log.Info("Web Service started!")
	}
}

func InitDB() *gorm.DB {

	connString := "host=host.docker.internal user=postgres password=toor dbname=taskpedia port=5432 sslmode=disable"
	dbConn, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = dbConn.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatal(err)
	}

	return dbConn
}

func InitMongo() *mongo.Client {
	mco := options.Client().ApplyURI(muri)

	mc, err := mongo.Connect(context.TODO(), mco)
	if err != nil {
		log.Fatal(err)
	}

	err = mc.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	dbList, errm := mc.ListDatabaseNames(context.TODO(), bson.D{})
	if errm != nil {
		log.Fatal(errm)
	}
	log.Infof("%+v", dbList)

	return mc
}

func InitSearchClient() *elastic.Client {

	client, err := elastic.NewClient(elastic.SetURL(euri))
	if err != nil {
		log.Fatal(err)
	}

	info, code, err := client.Ping(euri).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Ping on ElasticSearch returned with code %d and version %s", code, info.Version.Number)

	exists, err := client.IndexExists("tasks").Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		createIndex, err := client.CreateIndex("tasks").Body(models.Mapping).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if !createIndex.Acknowledged {
			log.Info("Create [tasks] Index not acknowledged")
		}
	}

	return client
}

func InitRedisClient() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr: ruri,
		DB:   0,
	})

	pong, err := client.Ping().Result()

	if err != nil {
		log.Fatalf("Ping to REDIS Server Failed!", err)
	} else {
		log.Infof("Ping to REDIS Server return: ", pong)
	}
	return client
}

func InitLogger() *logrus.Logger {

	var log = logrus.New()
	return log

}
