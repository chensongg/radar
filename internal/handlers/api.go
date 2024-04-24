package handlers

import (
	"radar/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)
	r.Route("/rooms", func(router chi.Router) {
		router.Use(middleware.Authorization)
		router.HandleFunc("/{roomID}/synclocation", SyncLocation)
	})

}
