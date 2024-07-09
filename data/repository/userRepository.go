package repository

import (
	"github.com/gocql/gocql"
	"saqrware.com/chat/data/model"
	"time"
)

type UserRepository struct {
	session *gocql.Session
}

func NewUserRepository(session *gocql.Session) *UserRepository {
	return &UserRepository{
		session: session,
	}
}

func (repo *UserRepository) CreateUser(user model.User) error {
	user.ID = gocql.TimeUUID()
	user.CreatedAt = time.Now()

	q := repo.session.Query(`INSERT INTO user (id, username, first_name, last_name, email, password, created_at) VALUES ( ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.CreatedAt)

	return q.Exec()
}

func (repo *UserRepository) GetUserByUsernameOrEmail(input string) (model.User, error) {
	var user model.User
	// Query by username
	q := repo.session.Query(`SELECT id, username, first_name, last_name, email, password, created_at FROM user WHERE username = ? LIMIT 1`, input)
	err := q.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err == nil {
		return user, nil
	}

	// Query by email if not found by username
	q = repo.session.Query(`SELECT id, username, first_name, last_name, email, password, created_at FROM user WHERE email = ? LIMIT 1`, input)
	err = q.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err == nil {
		return user, nil
	}

	return model.User{}, err // return error if not found
}
