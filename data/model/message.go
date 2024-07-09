package model

import (
	"github.com/gocql/gocql"
	"time"
)

type Message struct {
	ID        gocql.UUID
	Sender    gocql.UUID
	Receiver  gocql.UUID
	Content   string
	CreatedAt time.Time
}
