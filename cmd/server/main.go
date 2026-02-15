package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/theo-mazars/postrun/internal/handlers"
)

func main() {
  s := handlers.Load()
  r := chi.NewRouter()

  r.Use(middleware.RequestID)
  r.Use(middleware.RealIP)
  r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)

  r.Use(s.AuthMiddleware)

  r.Post("/send", s.SendHandler)

  http.ListenAndServe(":3000", r)
}
