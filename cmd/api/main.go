package main

import (
	"fmt"
	"net/http"
	"radar/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	handlers.Handler(r)
	fmt.Println("Starting go api service")
	err := http.ListenAndServe("localhost:8000", r)
	if err != nil {
		fmt.Println(err.Error())
	}
}
