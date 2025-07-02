package handlers

import (
	"encoding/json"
	"github.com/AliRizaAynaci/rlaas/internal/middleware"
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
	// 1) JSON body parse & validation
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProjectName == "" {
		http.Error(w, "Invalid request body or missing project_name", http.StatusBadRequest)
		return
	}

	userID := middleware.FromContext(r.Context())

	apiKey, err := service.GenerateAPIKey()
	if err != nil {
		http.Error(w, "Failed to generate API key", http.StatusInternalServerError)
		return
	}

	proj, err := service.CreateProject(r.Context(), userID, req.ProjectName, apiKey)
	if err != nil {
		http.Error(w, "Failed to create project (maybe duplicate name)", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{ApiKey: proj.ApiKey})
}
