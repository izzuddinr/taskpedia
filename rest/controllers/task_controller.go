package controllers

import (
	"context"
	"net/http"
	"strconv"
	"taskpedia-rest/models"
	"time"

	"github.com/labstack/echo"
)

func GetTask(c echo.Context) error {

	log.Info("Get Task accessed")

	var result []models.Task
	db.Where("Name IS NOT NULL").Find(&result)

	return c.JSON(http.StatusOK, result)

}

func CreateTask(c echo.Context) error {
	log.Info("Create Task accessed")

	req := new(models.Task)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	t := time.Now()

	task := models.Task{
		Name:      req.Name,
		Desc:      req.Desc,
		UserID:    req.UserID,
		Username:  req.Username,
		Status:    req.Status,
		CreatedAt: t,
		UpdatedAt: t,
	}

	err := db.Create(&task).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {

		publishNatsMsg(&task)

		var result map[string]string = map[string]string{
			"message":   "Task insertion successful!",
			"timestamp": time.Now().Format(time.RFC850),
		}

		AddTransactionLogTask("CREATE TASK", &task)
		return c.JSON(http.StatusOK, result)
	}
}

func UpdateTask(c echo.Context) error {
	log.Info("Create Task accessed")

	req := new(models.Task)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	t := time.Now()

	task := models.Task{
		ID:         req.ID,
		Name:       req.Name,
		Desc:       req.Desc,
		UserID:     req.UserID,
		Username:   req.Username,
		Status:     req.Status,
		UpdatedAt:  t,
		UpdateFlag: true,
	}

	var result []models.Task
	var err error

	err = db.Where("ID = ?", task.ID).Find(&result).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if len(result) < 1 {
		var result map[string]string = map[string]string{
			"message":   "Task Not Found!",
			"timestamp": time.Now().Format(time.RFC850),
		}

		return c.JSON(http.StatusOK, result)
	}

	err = db.Model(&models.Task{}).Where("ID = ?", task.ID).
		Updates(map[string]interface{}{
			"Name":      task.Name,
			"Desc":      task.Desc,
			"UserID":    task.UserID,
			"Username":  task.Username,
			"Status":    task.Status,
			"UpdatedAt": task.UpdatedAt,
		}).
		Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {

		publishNatsMsg(&task)

		var result map[string]string = map[string]string{
			"message":   "Task update successful!",
			"timestamp": time.Now().Format(time.RFC850),
		}
		AddTransactionLogTask("UPDATE TASK", &task)
		return c.JSON(http.StatusOK, result)
	}
}

func publishNatsMsg(req *models.Task) {

	sendChan := make(chan *models.Task)
	ec.BindSendChan("request_subject", sendChan)

	log.Infof("Sending request %d | %v | %v", req.ID, req.Name, req.Desc)
	sendChan <- req
}

func AddTransactionLogTask(ttype string, task *models.Task) {

	t := models.Transaction{
		Type:      ttype,
		Timestamp: time.Now(),
		Details: []string{
			"ID: " + strconv.FormatUint(uint64(task.ID), 10),
			"Name: " + task.Name,
			"Desc: " + task.Desc,
			"UserID: " + strconv.FormatUint(uint64(task.UserID), 10),
			"Username: " + task.Username,
			"Status: " + task.Status,
			"CreatedAt: " + task.CreatedAt.String(),
			"UpdatedAt: " + task.UpdatedAt.String(),
		},
	}
	collection.InsertOne(context.TODO(), t)
}
