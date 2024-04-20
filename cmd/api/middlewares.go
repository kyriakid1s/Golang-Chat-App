package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"chatapp.kyriakidis.net/internal/data"
	"chatapp.kyriakidis.net/internal/jwt"
	"golang.org/x/time/rate"
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
		user, err := app.models.Users.GetByUsernamePrivate(username)
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

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		// Add the "Vary: Access-Control-Request-Method" header.
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Check if the request has the HTTP method OPTIONS and contains the
			// "Access-Control-Request-Method" header. If it does, then we treat
			// it as a preflight request.
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				// Set the necessary preflight response headers, as discussed
				// previously.
				w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

				// Write the headers along with a 200 OK status and return from
				// the middleware with no further action.
				w.WriteHeader(http.StatusOK)
				return
			}

		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceedResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
