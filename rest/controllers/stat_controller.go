package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

func GetDataStats(c echo.Context) error {

	log.Info("Get User accesed")

	ucount, err := rc.Get("UserCount").Result()
	if err != nil {
		log.Error(err)
	}
	tcount, err := rc.Get("TaskCount").Result()
	if err != nil {
		log.Error(err)
	}

	result := make(map[string]string)
	result["UserCount"] = ucount
	result["TaskCount"] = tcount

	return c.JSON(http.StatusOK, result)
}
