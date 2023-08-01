package router

import (
	"net"
	"net/http"
)

type Handler struct {
	rl Ratelimiter
}

// middleware calls service func, so Ratelimiter is not coupled with http router libs / frameworks
type Ratelimiter interface {
	IsLimitOK(ip net.IP) bool
	ResetLimit(prefix string) bool
}

func NewRouter(rl Ratelimiter) http.Handler {
	h := Handler{
		rl: rl,
	}
	sm := http.NewServeMux()
	// handlers registration
	// ...
	// ...
	sm.HandleFunc("/", h.staticResource)

	return h.checkRatelimit(sm.ServeHTTP)
}
