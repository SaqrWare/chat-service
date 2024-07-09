package model

import (
	"github.com/gocql/gocql"
	"time"
)

type User struct {
	ID        gocql.UUID
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}
