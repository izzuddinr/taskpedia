package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"taskpedia-rest/models"

	"github.com/labstack/echo"
	"github.com/olivere/elastic"
)

var ctx = context.Background()

func SearchTaskByUser(c echo.Context) error {

	req := new(models.Task)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var taskList []models.Task

	searchQuery := elastic.NewTermQuery("UserID", "6")
	fmt.Printf("searchQuery: %v\n", searchQuery)
	searchResult, err := esc.Search().
		Index("tasks").
		Query(searchQuery).
		Pretty(true).
		Do(ctx)
	log.Error(err)

	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d tasks\n", searchResult.Hits.TotalHits)

		for _, hit := range searchResult.Hits.Hits {
			var t models.Task
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				log.Error(err)
			}
			taskList = append(taskList, t)
		}
	} else {
		fmt.Print("Found no tasks\n")
	}
	return c.JSON(http.StatusOK, taskList)

	// searchSrc := elastic.NewSearchSource()
	// searchSrc.Query(elastic.NewMatchQuery("UserID", req.UserID))
	// searchSvc := esc.Search().Index("tasks").Index("task").SearchSource(searchSrc)

	// result, err := searchSvc.Do(ctx)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
	// log.Info(result)

	// for _, hit := range result.Hits.Hits {
	// 	var task models.Task
	// 	err := json.Unmarshal(*hit.Source, &task)
	// 	if err != nil {
	// 		log.Errorf("Error while doing json.Unmarshal: %v", err)
	// 	}
	// 	taskList = append(taskList, task)
	// }

	// return c.JSON(http.StatusOK, taskList)
}
