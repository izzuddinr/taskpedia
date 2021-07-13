package main

import (
	"fmt"
	"os"
	"taskpedia-worker/models"
	"time"

	"github.com/go-redis/redis"
	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	cron "github.com/robfig/cron/v3"
)

var nuri = os.Getenv("NATS_URI")
var ruri = os.Getenv("REDIS_URI")

func main() {

	rc, err := InitRedisClient()
	if err != nil {
		log.Fatal("Failed to initiate Redis client!")
	}

	db := InitDB()

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
		panic(err)
	}
	defer ec.Close()

	log.Info("Connected to NATS and ready to receive messages")

	sch := cron.New(cron.WithSeconds())
	sch.AddFunc("0 * * ? * *", func() {
		SetCountStat(db, rc)
	})
	sch.Start()

	recvChan := make(chan *models.Task)
	ec.BindRecvChan("request_subject", recvChan)

	sendChan := make(chan *models.Task)
	ec.BindSendChan("request_subject", sendChan)

	for {
		// Wait for incoming messages
		req := <-recvChan
		resp := req
		log.Infof("Received request %+v", req)

		sendChan <- resp
	}
}

func InitRedisClient() (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr: ruri,
		DB:   0,
	})

	pong, err := client.Ping().Result()

	if err != nil {
		log.Error("Ping to REDIS Server Failed!")
	} else {
		log.Infof("Ping to REDIS Server return: ", pong)
	}
	return client, err
}

func InitDB() *gorm.DB {

	connString := "host=host.docker.internal user=postgres password=toor dbname=test_01 port=5432 sslmode=disable"
	dbConn, err := gorm.Open(postgres.Open(connString), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to DB!")
	}

	return dbConn
}

func SetCountStat(db *gorm.DB, rc *redis.Client) {

	var ucount, tcount int64
	db.Model(&models.User{}).Count(&ucount)
	db.Model(&models.Task{}).Count(&tcount)

	err := rc.Set("UserCount", ucount, 0).Err()
	if err != nil {
		log.Error(err)
	}

	err = rc.Set("TaskCount", tcount, 0).Err()
	if err != nil {
		log.Error(err)
	}
}

func GetCountStat(rc *redis.Client) {

	ucount, err := rc.Get("UserCount").Result()
	if err != nil {
		log.Error(err)
	}
	tcount, err := rc.Get("TaskCount").Result()
	if err != nil {
		log.Error(err)
	}

	log.Infof("User Count: %v, Task Count: %v", ucount, tcount)

}
