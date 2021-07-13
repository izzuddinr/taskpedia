package models

import "time"

type Task struct {
	ID         uint      `json:"ID"`
	Name       string    `json:"Name"`
	Desc       string    `json:"Desc"`
	UserID     uint      `json:"UserID"`
	Username   string    `json:"Username"`
	Status     string    `json:"Status"`
	CreatedAt  time.Time `json:"CreatedAt"`
	UpdatedAt  time.Time `json:"UpdatedAt"`
	UpdateFlag bool      `json:"UpdateFlag"`
}
