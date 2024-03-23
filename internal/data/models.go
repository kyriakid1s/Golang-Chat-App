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
}

func NewModels(db *redis.Client) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
