package service

import (
	"context"
	"github.com/AliRizaAynaci/gorl/core"
	"os"
	"rlaas/internal/database"
	"rlaas/internal/limiter"
	"time"
)

// GetRateLimitConfig retrieves rate limit configuration for a given API key and endpoint
func GetRateLimitConfig(apiKey, endpoint string) (limiter.RateLimitConfig, bool) {

	db, err := database.GetInstance()
	if err != nil {
		return limiter.RateLimitConfig{}, false
	}

	// First get the project ID
	var projectID int
	err = db.Pool().QueryRow(context.Background(), "SELECT id FROM projects WHERE api_key=$1", apiKey).Scan(&projectID)
	if err != nil {
		return limiter.RateLimitConfig{}, false
	}

	var cfg limiter.RateLimitConfig
	query := `SELECT strategy, key_by, limit_count, window_seconds FROM rate_limit_rules WHERE project_id=$1 AND endpoint=$2`
	row := db.Pool().QueryRow(context.Background(), query, projectID, endpoint)
	var strategy, keyBy string
	var limitCount, windowSeconds int
	err = row.Scan(&strategy, &keyBy, &limitCount, &windowSeconds)
	if err != nil {
		return limiter.RateLimitConfig{}, false
	}

	cfg.Strategy = core.StrategyType(strategy)
	cfg.KeyBy = core.KeyFuncType(keyBy)
	cfg.Limit = limitCount
	cfg.Window = time.Duration(windowSeconds) * time.Second

	// Redis cluster config
	cfg.RedisCluster = limiter.RedisClusterConfig{
		Nodes: []string{
			getEnvOrDefault("REDIS_NODE_1", "redis://localhost:6379/0"),
			getEnvOrDefault("REDIS_NODE_2", "redis://localhost:6380/0"),
			getEnvOrDefault("REDIS_NODE_3", "redis://localhost:6381/0"),
		},
		Strategy: getEnvOrDefault("SHARDING_STRATEGY", "hash_mod"),
	}

	return cfg, true
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
