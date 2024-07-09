package repository

import (
	"github.com/gocql/gocql"
	"saqrware.com/chat/data/model"
	"time"
)

type MessageRepository struct {
	session *gocql.Session
}

func NewMessageRepository(session *gocql.Session) *MessageRepository {
	return &MessageRepository{
		session: session,
	}
}

func (repo *MessageRepository) CreateMessage(message model.Message) error {
	message.ID = gocql.TimeUUID()
	message.CreatedAt = time.Now()

	q := repo.session.Query(`INSERT INTO message (id, sender_id, receiver_id, content, created_at) VALUES ( ?, ?, ?, ?, ?)`,
		message.ID, message.Sender, message.Receiver, message.Content, message.CreatedAt)

	return q.Exec()
}

func (repo *MessageRepository) GetMessagesForUserWithPagination(userID gocql.UUID, lastID string, limit int) ([]model.Message, error) {
	var messages []model.Message
	var lastMessageID gocql.UUID
	var err error

	if lastID != "" {
		lastMessageID, err = gocql.ParseUUID(lastID)
		if err != nil {
			return nil, err
		}
	}

	query := `SELECT id, sender_id, receiver_id, content, created_at FROM message WHERE receiver_id = ? AND id > ? LIMIT ? ALLOW FILTERING`
	iter := repo.session.Query(query, userID, lastMessageID, limit).Iter()

	var message model.Message
	for iter.Scan(&message.ID, &message.Sender, &message.Receiver, &message.Content, &message.CreatedAt) {
		messages = append(messages, message)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return messages, nil
}
