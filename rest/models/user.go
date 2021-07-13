package models

type User struct {
	ID   uint   `gorm:"primaryKey"  json:"ID"`
	Name string `json:"Name"`
}
