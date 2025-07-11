package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AliRizaAynaci/rlaas/internal/middleware"
	"github.com/AliRizaAynaci/rlaas/internal/service"
)

type MeResponse struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.FromContext(r.Context())

	// Get user info from database
	user, err := service.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := MeResponse{
		UserID: user.ID,
		Email:  user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
