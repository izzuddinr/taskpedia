package main

import (
	"context"
	"fmt"
	"os"
	"taskpedia-search/models"
	"time"

	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	"github.com/olivere/elastic"
)

var ctx = context.Background()
var euri = os.Getenv("ELASTIC_URI")
var nuri = os.Getenv("NATS_URI")

func main() {

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

	recvChan := make(chan *models.Task)
	ec.BindRecvChan("request_subject", recvChan)

	esc, err := InitSearchClient()
	if err != nil {
		log.Fatalln("Failed to initialize ES Client! ", err)
	}

	for {
		// Wait for incoming messages
		req := <-recvChan

		log.Infof("Received request %+v", req)

		InsertSearchData(esc, req)
	}
}

func InitSearchClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL(euri))
	if err != nil {
		log.Panic(err)
	}

	info, code, err := client.Ping(euri).Do(ctx)
	if err != nil {
		log.Panic(err)
	}
	log.Infof("Ping on ElasticSearch returned with code %d and version %s", code, info.Version.Number)

	exists, err := client.IndexExists("tasks").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("tasks").Body(models.Mapping).Do(ctx)
		if err != nil {
			// Handle error
			log.Panic(err)
		}
		if !createIndex.Acknowledged {
			log.Info("Create [tasks] Index not acknowledged")
		}
	}

	return client, err

}

func InsertSearchData(esc *elastic.Client, req *models.Task) error {

	index, err := esc.Index().
		Index("tasks").
		Type("task").
		BodyJson(req).
		Do(ctx)

	if err != nil {
		log.Error(err)
		return err
	} else {
		log.Infof("SUCCESS: Insert task [%+v] with ID: %s", req, index.Id)
		return nil
	}
}
