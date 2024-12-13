package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/victorthury/cep-api/configs"
	"github.com/victorthury/cep-api/internal/webserver/handlers"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	cepHandler := handlers.NewCepHandler(configs.BrasilApiUrl, configs.ViaCepUrl)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/cep/{cep}", cepHandler.GetCep)

	http.ListenAndServe(":8000", r)
}
