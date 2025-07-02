package handlers

import (
	"encoding/json"
	"github.com/AliRizaAynaci/rlaas/internal/service"
	"net/http"
)

type RegisterRequest struct {
	ProjectName string `json:"project_name"`
}

type RegisterResponse struct {
	ApiKey string `json:"api_key"`
}

func RegisterProjectHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProjectName == "" {
		http.Error(w, "Invalid request body or missing project_name", http.StatusBadRequest)
		return
	}

	apiKey, err := service.GenerateAPIKey()
	if err != nil {
		http.Error(w, "Failed to generate API key", http.StatusInternalServerError)
		return
	}

	err = service.CreateProject(req.ProjectName, apiKey)
	if err != nil {
		http.Error(w, "Failed to create project (maybe duplicate name)", http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{ApiKey: apiKey}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
