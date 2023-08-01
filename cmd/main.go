package main

import (
	"net/http"
	"time"

	"github.com/hablof/antibot-ratelimit/internal/router"
	"github.com/hablof/antibot-ratelimit/internal/service"
)

func main() {
	s := http.Server{
		Addr:         ":8050",
		Handler:      router.NewRouter(service.NewRatelimiter()),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	s.ListenAndServe()
}
