package controllers

import (
	"context"
	"net/http"
	"taskpedia-rest/models"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTransactionsLog(c echo.Context) error {

	log.Info("Get Transaction Log accessed")

	var result []models.Transaction

	findOptions := options.Find()
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		log.Error(err)
	}

	log.Infof("%+v", cursor)

	for cursor.Next(context.TODO()) {
		var cur models.Transaction
		err := cursor.Decode(&cur)
		if err != nil {
			log.Error(err)
		}
		result = append(result, cur)
	}

	return c.JSON(http.StatusOK, result)

}
