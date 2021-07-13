package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"taskpedia-rest/models"

	"github.com/labstack/echo"
	"github.com/olivere/elastic/v7"
)

var ctx = context.Background()

func SearchTaskByUserID(c echo.Context) error {

	req := new(models.Task)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var taskList []models.Task

	searchQuery := elastic.NewTermQuery("UserID", req.UserID)
	searchResult, err := esc.Search().
		Index("tasks").
		Query(searchQuery).
		Pretty(true).
		Do(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if searchResult.Hits.TotalHits.Value > 0 {

		for _, hit := range searchResult.Hits.Hits {
			var t models.Task
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				log.Error(err)
			}
			taskList = append(taskList, t)
		}
		return c.JSON(http.StatusOK, taskList)
	} else {
		return c.JSON(http.StatusOK, "Found no Task!")
	}
}
