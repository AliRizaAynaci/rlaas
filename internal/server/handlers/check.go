package handlers

import (
	"encoding/json"
	"github.com/AliRizaAynaci/rlaas/internal/limiter"
	"github.com/AliRizaAynaci/rlaas/internal/service"
	"net/http"
)

// CheckRequest represents the rate limit check request structure
type CheckRequest struct {
	ApiKey   string `json:"api_key"`
	Endpoint string `json:"endpoint"`
	Key      string `json:"key"` // e.g., IP, userID, etc.
}

// CheckResponse represents the rate limit check response structure
type CheckResponse struct {
	Allowed bool `json:"allowed"`
}

// RateLimitCheck handles the rate limit check request and returns whether the request is allowed
func RateLimitCheck(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config, ok := service.GetRateLimitConfig(req.ApiKey, req.Endpoint)
	if !ok {
		http.Error(w, "Config not found", http.StatusNotFound)
		return
	}

	lim, err := limiter.GetLimiterForKey(req.ApiKey, req.Endpoint, req.Key, config)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	allowed, err := lim.Allow(req.Key)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
		return
	}

	resp := CheckResponse{Allowed: allowed}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
