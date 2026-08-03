package handler

import (
	"encoding/json"
	"myapp/internal/domain"
	"myapp/internal/utils"
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
	w.WriteHeader(utils.MapError(domainErr.Code))
	json.NewEncoder(w).Encode(domainErr)
}
