package controllers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

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
	if err != nil {
		return c.JSON(http.StatusOK, searchResult)
	}
	log.Info(searchResult)

	var rType models.Task
	for _, item := range searchResult.Each(reflect.TypeOf(rType)) {
		t := item.(models.Task)
		taskList = append(taskList, t)
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
