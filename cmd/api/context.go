package main

import (
	"context"
	"net/http"

	"chatapp.kyriakidis.net/internal/data"
)

type contextKey string

const contextUserKey = contextKey("user")

func (app *application) setUserContext(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), contextUserKey, user)
	return r.WithContext(ctx)
}

func (app *application) getUserContext(r *http.Request) *data.User {
	user, ok := r.Context().Value(contextUserKey).(*data.User)
	if !ok {
		panic("no user found in context")
	}
	return user
}
