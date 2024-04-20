package main

import (
	"net/http"
	"time"

	"chatapp.kyriakidis.net/internal/data"
	"chatapp.kyriakidis.net/internal/jwt"
	"chatapp.kyriakidis.net/internal/validator"
)

func (app *application) registerHadler(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := app.readJSON(w, r, &user)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateUser(v, &user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user.CreatedAt = time.Now()
	err = app.models.Users.Insert(&user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	token, err := jwt.NewToken(user.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.createCookie(w, "jwt", token, "/")
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user.PublicUser}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	v.Check(input.Username != "", "username", "can not be empty")
	v.Check(input.Password != "", "password", "can not be empty")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.models.Users.Get(input.Username, input.Password)
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}
	//Token creation
	token, err := jwt.NewToken(user.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.createCookie(w, "jwt", token, "/")

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
func (app *application) getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	username := app.getUserContext(r).Username
	err := app.writeJSON(w, http.StatusOK, envelope{"user": username}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
