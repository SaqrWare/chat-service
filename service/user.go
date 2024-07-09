package service

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"saqrware.com/chat/data"
	"saqrware.com/chat/data/model"
	"saqrware.com/chat/data/repository"
	"saqrware.com/chat/helper"
	dto2 "saqrware.com/chat/http/dto"
	"time"
)

func RegisterUser(userDto dto2.RegisterUserDto) error {
	var user model.User
	user.Username = userDto.Username
	user.FirstName = userDto.FirstName
	user.LastName = userDto.LastName
	user.Email = userDto.Email
	// encrypt password bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error", err)

		return err
	}
	user.Password = string(hashedPassword)

	// save user to cassandra
	userRepo := repository.NewUserRepository(data.CassandraSession)
	err = userRepo.CreateUser(user)
	if err != nil {
		fmt.Println("Error creating user", err)
		return err
	}
	return nil
}

func UserLogin(loginDto dto2.UserLoginDto) (string, error) {
	userRepo := repository.NewUserRepository(data.CassandraSession)
	user, err := userRepo.GetUserByUsernameOrEmail(loginDto.Identifier)
	if err != nil {
		return "", err
	}
	// check password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))
	if err != nil {
		return "", errors.New("WRONG_CREDENTIALS")
	}

	// [X] Generate random token
	token, err := helper.GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	// Save token and user data in redis
	_, err = data.RedisClient.HSet(data.RCTX, "session:"+token, "id", user.ID.String(), "username", user.Username, "email", user.Email).Result()
	if err != nil {
		return "", err
	}
	_, err = data.RedisClient.Expire(data.RCTX, "session:"+token, 30*24*time.Hour).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}
