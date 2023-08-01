package router

import (
	"fmt"
	"net/http"
)

func (h *Handler) staticResource(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "static content")
}
