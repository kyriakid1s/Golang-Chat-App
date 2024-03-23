package main

import (
	"net/http"
	"time"

	"chatapp.kyriakidis.net/internal/data"
	"chatapp.kyriakidis.net/internal/jwt"
)

func (app *application) registerHadler(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := app.readJSON(w, r, &user)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user.CreatedAt = time.Now()
	err = app.models.Users.Insert(&user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.Password = ""
	token, err := jwt.NewToken(user.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.createCookie(w, "jwt", token, "/")
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
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
	user, err := app.models.Users.Get(input.Username, input.Password)
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}
	//*Done this to omit password from response
	user.Password = ""
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