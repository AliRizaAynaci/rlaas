package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AliRizaAynaci/rlaas/internal/limiter"
	"github.com/AliRizaAynaci/rlaas/internal/logging"
	"github.com/AliRizaAynaci/rlaas/internal/service"
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

	config, err := service.GetRateLimitConfig(req.ApiKey, req.Endpoint)
	if errors.Is(err, service.ErrProjectNotFound) {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}
	if errors.Is(err, service.ErrEndpointNotOwned) {
		http.Error(w, "Endpoint does not belong to this project", http.StatusForbidden)
		return
	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
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
		// 1) Log the rate‐limit hit
		logging.L.Warn("Rate limit exceeded",
			slog.String("component", "limiter"),
			slog.String("api_key", req.ApiKey),
			slog.String("endpoint", req.Endpoint),
			slog.String("key", req.Key),
		)
		// 2) Return the HTTP 429
		http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
		return
	}

	resp := CheckResponse{Allowed: allowed}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
