package models

import (
	"time"
)

type DetailUser struct {
	User
	DeletedTime *time.Time `json:"delete_time"`
}

type DetailUsers []DetailUser
