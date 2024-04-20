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
	PublicUser
	Password  string    `json:"password,omitempty"`
	Rooms     []string  `json:"rooms,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type PublicUser struct {
	Username string `json:"username"`
	Online   bool   `json:"online"`
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

func (m UserModel) Get(username, password string) (*PublicUser, error) {
	userJSON := m.DB.Get(context.Background(), username).Val()
	var user User
	err := json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	user.Online = true
	newUserJSON, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = m.DB.Set(context.Background(), user.Username, newUserJSON, 0).Err()
	if err != nil {
		return nil, err
	}
	return &user.PublicUser, nil
}

func (m UserModel) GetByUsername(username string) (*PublicUser, error) {
	var user User
	userJSON, err := m.DB.Get(context.Background(), username).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNoRecordFound
	}
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, err
	}
	return &user.PublicUser, nil
}
func (m UserModel) GetByUsernamePrivate(username string) (*User, error) {
	var user User
	userJSON, err := m.DB.Get(context.Background(), username).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNoRecordFound
	}
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (m UserModel) GetUserForChat(usernames []string) ([]PublicUser, error) {
	var users []PublicUser
	for _, username := range usernames {
		user, err := m.GetByUsername(username)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}
