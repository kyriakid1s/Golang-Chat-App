package data

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNoRecordFound = errors.New("no record found")
)

type Models struct {
	Users UserModel
	Chats ChatModel
}

func NewModels(db *redis.Client) Models {
	return Models{
		Users: UserModel{DB: db},
		Chats: ChatModel{DB: db},
	}
}
