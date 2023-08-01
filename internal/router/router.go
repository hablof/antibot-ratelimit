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
	IsRatelimitOK(ip net.IP) bool
	ResetLimit(prefix string) bool
}

func NewRouter() http.Handler {
	h := Handler{}
	sm := http.NewServeMux()
	// handlers registration
	// ...
	// ...

	return h.checkRatelimit(sm.ServeHTTP)
}
