package models

import "time"

type Transaction struct {
	Type      string
	Timestamp time.Time
	Details   []string
}
