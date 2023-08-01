package router

import (
	"net"
	"net/http"
)

func (h *Handler) checkRatelimit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ip := net.ParseIP(r.Header.Get("X-Forwarded-For"))

		if h.Ratelimiter.IsRatelimitOK(ip) {
			f(w, r)
		} else {
			h.requestLimitExceededHandler(w, r)
		}
	}
}

func (h *Handler) requestLimitExceededHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte("no pressure, no pressure..."))
}
