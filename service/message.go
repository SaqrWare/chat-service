package service

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"saqrware.com/chat/data"
	"saqrware.com/chat/data/model"
	"saqrware.com/chat/data/repository"
	"saqrware.com/chat/http/dto"
	"time"
)

func CreateMessage(dto dto.SendMessageDto, senderID string) error {
	log.Println("Starting message creation")
	var message model.Message
	var err error
	message.Sender, err = gocql.ParseUUID(senderID)
	if err != nil {
		log.Printf("Error parsing senderID to UUID: %v", err)
		return err
	}

	message.Receiver, err = gocql.ParseUUID(dto.Receiver)
	if err != nil {
		log.Printf("Error parsing receiverID to UUID: %v", err)
		return err
	}

	message.Content = dto.Message

	messageRepo := repository.NewMessageRepository(data.CassandraSession)
	err = messageRepo.CreateMessage(message)
	if err != nil {
		log.Printf("Error creating message: %v", err)
		return err
	}

	// Destroy last cached history
	cacheKey := generateCacheKey(message.Receiver.String(), "", 10)
	err = data.RedisClient.Del(data.RCTX, cacheKey).Err()
	if err != nil {
		log.Printf("Error deleting cache: %v", err)
		return err
	}

	log.Println("Message created successfully")
	return nil
}

func generateCacheKey(userID string, lastID string, limit int) string {
	return fmt.Sprintf("user:%s:messages:lastID:%s:limit:%d", userID, lastID, limit)
}

func GetMessageHistory(userID string, lastID string, limit int) ([]model.Message, error) {
	cacheKey := generateCacheKey(userID, lastID, limit)

	// Check cache
	cachedData, err := data.RedisClient.Get(data.RCTX, cacheKey).Result()
	if err == nil && cachedData != "" {
		var messages []model.Message
		if err := json.Unmarshal([]byte(cachedData), &messages); err == nil {
			return messages, nil
		}
	}

	// If no data in cache -> retrieve from database
	messageRepo := repository.NewMessageRepository(data.CassandraSession)
	userUUID, err := gocql.ParseUUID(userID)
	if err != nil {
		return nil, err
	}

	messages, err := messageRepo.GetMessagesForUserWithPagination(userUUID, lastID, limit)
	if err != nil {
		return nil, err
	}

	// Serialize and cache the data
	serializedData, err := json.Marshal(messages)
	if err == nil {
		// Save with 1 minute expiration
		err := data.RedisClient.Set(data.RCTX, cacheKey, serializedData, time.Minute).Err()
		if err != nil {
			return nil, err
		}
	}

	return messages, nil
}
