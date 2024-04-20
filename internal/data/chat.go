package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Message struct {
	ID      uuid.UUID `json:"id"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	Room    string    `json:"room"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

type Chat struct {
	Members  []User
	Messages []Message
}

type ChatModel struct {
	DB *redis.Client
}

// sortUsernames sorts room's users usernames.
func SortUsernames(roomName *string) {
	names := strings.Split(*roomName, ":")
	sort.Strings(names)
	*roomName = strings.Join(names, ":")

}

func (m ChatModel) getByUsername(username string) (*User, error) {
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

// checkRoomIfExists checks on both sender and recipient if a room exists between them. If not, creates one.
func (m ChatModel) checkRoomIfExists(sender, recipient *User, messageRoom string) error {
	var foundSenderRoom bool
	var foundRecipientRoom bool

	for _, room := range sender.Rooms {
		if room == messageRoom {
			foundSenderRoom = true
			fmt.Println("sender found")
			break
		}
	}
	for _, room := range recipient.Rooms {
		if room == messageRoom {
			foundRecipientRoom = true
			fmt.Println("recipient found")
			break
		}
	}
	if !foundSenderRoom {
		sender.Rooms = append(sender.Rooms, messageRoom)
		data, err := json.Marshal(sender)
		fmt.Println("sender not found")
		if err != nil {
			return err
		}

		return m.DB.Set(context.Background(), sender.Username, data, 0).Err()
	}
	if !foundRecipientRoom {
		recipient.Rooms = append(recipient.Rooms, messageRoom)
		data, err := json.Marshal(recipient)
		fmt.Println("recipient not found")
		if err != nil {
			return err
		}
		return m.DB.Set(context.Background(), recipient.Username, data, 0).Err()
	}
	return nil
}

func (m ChatModel) SaveMessage(message *Message) error {
	sender, err := m.getByUsername(message.From)
	if err != nil {
		return err
	}
	recipient, err := m.getByUsername(message.To)
	if err != nil {
		return err
	}

	//Sort room name so we have one room for both users
	SortUsernames(&message.Room)
	//Check if room exists on both sender and recipient. If not, adds it.
	err = m.checkRoomIfExists(sender, recipient, message.Room)
	if err != nil {
		return err
	}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = m.DB.ZAdd(context.Background(), message.Room, redis.Z{Score: float64(message.Date.Unix()), Member: data}).Err()
	if err != nil {
		return err
	}
	return nil
}

func (m ChatModel) GetMessages(room string) ([]Message, error) {
	messages, err := m.DB.ZRange(context.Background(), room, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var chat []Message
	for _, message := range messages {
		var m Message
		if err := json.Unmarshal([]byte(message), &m); err != nil {
			return nil, err
		}
		chat = append(chat, m)
	}
	return chat, nil
}

func roomNameToUsernames(room string) []string {
	names := strings.Split(room, ":")
	return names
}

// TODO: !! returns user's password, fix it
func (m ChatModel) GetUserChats(user *User) ([]Chat, error) {
	var chats []Chat
	for _, room := range user.Rooms {
		var chat Chat
		names := roomNameToUsernames(room)
		for _, name := range names {
			if name != user.Username {
				chatUser, err := m.getByUsername(name)
				if err != nil {
					return nil, err
				}
				//*Do this to omit these fields.
				chatUser.Rooms = []string{}
				chatUser.Password = ""
				fmt.Println(chatUser)
				chat.Members = append(chat.Members, *chatUser)
				messages, err := m.GetMessages(room)
				if err != nil {
					return nil, err
				}
				chat.Messages = append(chat.Messages, messages...)
				chats = append(chats, chat)
			}
		}
	}
	return chats, nil
}
