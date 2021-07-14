package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"

	"taskpedia-rest/models"
)

func GetUser(c echo.Context) error {

	log.Info("Get User accesed")

	var result []models.User
	db.Where("Name IS NOT NULL").Find(&result)

	return c.JSON(http.StatusOK, result)
}

func InsertUser(c echo.Context) error {

	log.Info("Create User accesed")

	req := new(models.User)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user := models.User{
		Name: req.Name,
	}

	err := db.Create(&user).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		var result map[string]string = map[string]string{
			"message":   "User insertion successful!",
			"timestamp": time.Now().Format(time.RFC850),
		}
		AddTransactionLogUser("CREATE USER", &user)
		return c.JSON(http.StatusOK, result)
	}

}

func UpdateUser(c echo.Context) error {

	log.Info("Create User accesed")

	req := new(models.User)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user := models.User{
		ID:   req.ID,
		Name: req.Name,
	}

	err := db.Model(&models.User{}).Where("ID = ?", user.ID).Update("Name", user.Name).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		var result map[string]string = map[string]string{
			"message":   "User update successful!",
			"timestamp": time.Now().Format(time.RFC850),
		}
		AddTransactionLogUser("UPDATE USER", &user)
		return c.JSON(http.StatusOK, result)
	}
}

func AddTransactionLogUser(ttype string, user *models.User) {

	t := models.Transaction{
		Type:      ttype,
		Timestamp: time.Now(),
		Details: []string{
			"ID: " + strconv.FormatUint(uint64(user.ID), 10),
			"Name: " + user.Name},
	}
	collection.InsertOne(context.TODO(), t)
}
