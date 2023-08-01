package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

const (
	resetLimitHandlerURL = "/resetlimit"
)

func (h *Handler) checkRatelimit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.EscapedPath() == resetLimitHandlerURL && r.Method == http.MethodPost {
			h.resetLimitHandle(w, r)
			return
		}

		ip := net.ParseIP(r.Header.Get("X-Forwarded-For"))
		if h.rl.IsRatelimitOK(ip) {
			f(w, r)
		} else {
			h.requestLimitExceededHandle(w, r)
		}
	}
}

func (h *Handler) requestLimitExceededHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte("no pressure, no pressure..."))
}

func (h *Handler) resetLimitHandle(w http.ResponseWriter, r *http.Request) {
	type jsonReq struct {
		Prefix string `json:"prefix"`
	}

	type jsonResp struct {
		Success bool `json:"success"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad body")

		return
	}
	r.Body.Close()

	req := jsonReq{}
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad json")

		return
	}

	resp := jsonResp{
		Success: h.rl.ResetLimit(req.Prefix),
	}

	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "error")

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
