package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/AliRizaAynaci/rlaas/internal/logging"
	"github.com/AliRizaAynaci/rlaas/internal/service"
)

// RuleResponse represents the rule structure that will be sent to the frontend
type RuleResponse struct {
	ID            int    `json:"id"`
	Endpoint      string `json:"endpoint"`
	Strategy      string `json:"strategy"`
	KeyBy         string `json:"key_by"`
	LimitCount    int    `json:"limit_count"`
	WindowSeconds int    `json:"window_seconds"`
}

func AddRule(w http.ResponseWriter, r *http.Request) {
	// Get API key from header
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" {
		http.Error(w, "API key is required", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
		apiKey = apiKey[7:]
	}

	// Get rule data from request body
	var rule service.RateLimitRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if rule.Endpoint == "" {
		http.Error(w, "Endpoint is required", http.StatusBadRequest)
		return
	}
	if rule.Strategy == "" {
		http.Error(w, "Strategy is required", http.StatusBadRequest)
		return
	}
	if rule.KeyBy == "" {
		http.Error(w, "Key By is required", http.StatusBadRequest)
		return
	}
	if rule.LimitCount <= 0 {
		http.Error(w, "Limit Count must be greater than 0", http.StatusBadRequest)
		return
	}
	if rule.WindowSeconds <= 0 {
		http.Error(w, "Window Seconds must be greater than 0", http.StatusBadRequest)
		return
	}

	// Find the project
	project, err := service.GetProjectByAPIKey(apiKey)
	if err != nil {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	// Add project ID to the rule
	rule.ProjectID = project.ID

	// Save the rule
	err = service.AddRule(&rule)
	if err != nil {
		http.Error(w, "Failed to add rule", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := RuleResponse{
		ID:            rule.ID,
		Endpoint:      rule.Endpoint,
		Strategy:      rule.Strategy,
		KeyBy:         rule.KeyBy,
		LimitCount:    rule.LimitCount,
		WindowSeconds: rule.WindowSeconds,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetRules(w http.ResponseWriter, r *http.Request) {
	// Get API key from header
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" {
		http.Error(w, "API key is required", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
		apiKey = apiKey[7:]
	}

	// Find the project first
	project, err := service.GetProjectByAPIKey(apiKey)
	if err != nil {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	// Get rules for the project
	rules, err := service.GetRulesByProjectID(project.ID)
	if err != nil {
		http.Error(w, "Failed to fetch rules", http.StatusInternalServerError)
		return
	}

	// Convert from Model to Response
	ruleResponses := make([]RuleResponse, len(rules))
	for i, rule := range rules {
		ruleResponses[i] = RuleResponse{
			ID:            rule.ID,
			Endpoint:      rule.Endpoint,
			Strategy:      rule.Strategy,
			KeyBy:         rule.KeyBy,
			LimitCount:    rule.LimitCount,
			WindowSeconds: rule.WindowSeconds,
		}
	}

	// Prepare response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ruleResponses)
}

func DeleteRule(w http.ResponseWriter, r *http.Request) {
	// Get API key from header
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" {
		http.Error(w, "API key is required", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
		apiKey = apiKey[7:]
	}

	// Get Rule ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Rule ID is required", http.StatusBadRequest)
		return
	}
	ruleID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	// Find the project
	project, err := service.GetProjectByAPIKey(apiKey)
	if err != nil {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	// Delete the rule
	err = service.DeleteRule(ruleID, project.ID)
	if err != nil {
		// 1) Log the failure with context
		logging.L.Warn("Failed to delete rate limit rule",
			slog.String("component", "rules"),
			slog.Int("rule_id", ruleID),
			slog.String("error", err.Error()),
		)
		// 2) Return the HTTP error as before
		http.Error(w, "Failed to delete rule", http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Rule deleted successfully",
	})
}

func UpdateRule(w http.ResponseWriter, r *http.Request) {
	// Get API key from header
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" {
		http.Error(w, "API key is required", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
		apiKey = apiKey[7:]
	}

	// Get Rule ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Rule ID is required", http.StatusBadRequest)
		return
	}
	ruleID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	// Get rule data from request body
	var rule service.RateLimitRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find the project
	project, err := service.GetProjectByAPIKey(apiKey)
	if err != nil {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	// Update the rule
	rule.ID = ruleID
	rule.ProjectID = project.ID
	err = service.UpdateRule(&rule)
	if err != nil {
		http.Error(w, "Failed to update rule", http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Rule updated successfully",
	})
}
