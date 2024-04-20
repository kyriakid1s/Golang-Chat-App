package main

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"chatapp.kyriakidis.net/internal/data"
	"github.com/google/uuid"
)

func (app *application) sendHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		From    string `json:"-"`
		To      string `json:"to"`
		Room    string `json:"room"`
		Message string `json:"message"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	input.From = app.getUserContext(r).Username
	message := &data.Message{
		From:    input.From,
		Room:    fmt.Sprintf("%s:%s", input.From, input.To),
		To:      input.To,
		Message: input.Message,
		ID:      uuid.New(),
		Date:    time.Now(),
	}

	err = app.models.Chats.SaveMessage(message)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"status": "success"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("name")
	//TODO: fix this shit line.
	room := fmt.Sprintf("%s:%s", r.Context().Value(contextUserKey).(*data.User).Username, username)
	data.SortUsernames(&room)
	//Check if user is authorized to get the room's message
	user := app.getUserContext(r)
	if !slices.Contains(user.Rooms, room) {
		app.notPermittedResponse(w, r)
		return
	}
	messages, err := app.models.Chats.GetMessages(room)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"messages": messages}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getChatsHandler(w http.ResponseWriter, r *http.Request) {
	chats, err := app.models.Chats.GetUserChats(app.getUserContext(r))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"chats": chats}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
