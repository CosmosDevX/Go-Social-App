package utils

import (
	"encoding/json"
	"myapp/internal/domain"
	"net"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "error during writing response", http.StatusBadRequest)
		return
	}
}

func WriteError(w http.ResponseWriter, domainErr domain.DomainError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(MapError(domainErr.Code))
	json.NewEncoder(w).Encode(domainErr)
}

func GetIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
