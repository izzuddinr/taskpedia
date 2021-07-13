package models

import "time"

type Task struct {
	ID         uint      `gorm:"primaryKey" json:"ID"`
	Name       string    `json:"Name"`
	Desc       string    `json:"Desc"`
	UserID     uint      `json:"UserID"`
	Username   string    `json:"Username"`
	Status     string    `json:"Status"`
	CreatedAt  time.Time `json:"CreatedAt"`
	UpdatedAt  time.Time `json:"UpdatedAt"`
	UpdateFlag bool      `gorm:"-" json:"UpdateFlag"`
}

const Mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	}
}`
