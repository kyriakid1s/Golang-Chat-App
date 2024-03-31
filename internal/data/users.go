package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"chatapp.kyriakidis.net/internal/validator"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type User struct {
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type UserModel struct {
	DB *redis.Client
}

// Validate User
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "username", "can not be empty")
	v.Check(len(user.Password) >= 8, "password", "must be at least 8 characters")
}

func (m UserModel) Insert(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = m.DB.Set(context.Background(), user.Username, userJSON, 0).Err()
	if err != nil {
		return err
	}
	//Add user to users
	err = m.DB.SAdd(context.Background(), "users", user.Username).Err()
	if err != nil {
		//if err, delete the user
		m.DB.Del(context.Background(), user.Username)
		return err
	}

	return nil
}

func (m UserModel) Get(username, password string) (*User, error) {
	userJSON := m.DB.Get(context.Background(), username).Val()
	var user User
	err := json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return &user, nil
}

func (m UserModel) GetByUsername(username string) (*User, error) {
	var user User
	userJSON := m.DB.Get(context.Background(), username)
	if errors.Is(userJSON.Err(), redis.Nil) {
		return nil, ErrNoRecordFound
	}
	err := json.Unmarshal([]byte(userJSON.Val()), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}
