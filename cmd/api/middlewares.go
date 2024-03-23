package main

import (
	"errors"
	"net/http"

	"chatapp.kyriakidis.net/internal/data"
	"chatapp.kyriakidis.net/internal/jwt"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("jwt")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				r = app.setUserContext(r, data.AnonymousUser)
				next.ServeHTTP(w, r)
				return
			default:
				app.serverErrorResponse(w, r, err)
				return
			}
		}

		username, err := jwt.VerifyToken(token.Value)
		//TODO: Fix errors (switch)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		user, err := app.models.Users.GetByUsername(username)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNoRecordFound):
				app.invalidAuthenticationTokenResponse(w, r)
				return
			default:
				app.serverErrorResponse(w, r, err)
				return
			}
		}
		r = app.setUserContext(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.getUserContext(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
