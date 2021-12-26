package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *app) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(a.jwt.Verifier())
		r.Use(a.jwt.Authenticator)
		// r.Get("/users", a.getAllUsers)
		r.Put("/users", a.putUser)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/signup", a.userSignup)
		r.Post("/login", a.userLogin)
	})
	return r
}
