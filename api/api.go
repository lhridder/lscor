package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"log"
	"lscor/config"
	"lscor/corero"
	"net/http"
	"time"
)

func Start() {
	cfg := config.GlobalConfig

	r := chi.NewRouter()
	r.Use(middleware.RealIP)

	if cfg.Debug {
		r.Use(middleware.Logger)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	r.Route("/api", func(r chi.Router) {
		r.Use(httprate.LimitByIP(100, 1*time.Minute))
		r.Get("/", GetHome)
		r.Get("/recent", GetRecent)
		r.Get("/top", GetTop)
	})

	listen := cfg.Listen
	log.Printf("Starting web listener on %s", listen)
	panic(http.ListenAndServe(listen, r))
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("api"))
}

func GetRecent(w http.ResponseWriter, r *http.Request) {
	timeframe := r.URL.Query().Get("timeframe")
	target := r.URL.Query().Get("target")
	duration, err := time.ParseDuration(timeframe)
	if err != nil {
		log.Printf("Invalid timeframe received: %s", err)
		http.Error(w, "invalid timeframe", http.StatusBadRequest)
		return
	}

	attacks, err := corero.GetRecent(duration, target)
	if err != nil {
		log.Printf("Failed to get recent attacks: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(attacks)
	if err != nil {
		log.Printf("Failed to serialize attacks: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(data)
}

func GetTop(w http.ResponseWriter, r *http.Request) {
	timeframe := r.URL.Query().Get("timeframe")
	duration, err := time.ParseDuration(timeframe)
	if err != nil {
		log.Printf("Invalid timeframe received: %s", err)
		http.Error(w, "invalid timeframe", http.StatusBadRequest)
		return
	}

	top, err := corero.GetTop(duration)
	if err != nil {
		log.Printf("Failed to get top ips: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(top)
	if err != nil {
		log.Printf("Failed to top ips: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(data)
}
