package main

import (
	"log"
	"net/http"
	"time"

	"github.com/hablof/antibot-ratelimit/internal/config"
	"github.com/hablof/antibot-ratelimit/internal/router"
	"github.com/hablof/antibot-ratelimit/internal/service"
)

func main() {
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		log.Println(err)
		return
	}

	s := http.Server{
		Addr:         ":8050",
		Handler:      router.NewRouter(service.NewRatelimiter(cfg)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	s.ListenAndServe()
}
