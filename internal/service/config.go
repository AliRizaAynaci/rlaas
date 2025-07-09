package service

import (
	"context"
	"errors"
	"github.com/AliRizaAynaci/gorl/core"
	"github.com/AliRizaAynaci/rlaas/internal/database"
	"github.com/AliRizaAynaci/rlaas/internal/limiter"
	"os"
	"time"
)

var (
	ErrProjectNotFound  = errors.New("project not found for given API key")
	ErrEndpointNotOwned = errors.New("endpoint does not belong to this project")
)

// GetRateLimitConfig retrieves rate limit configuration for a given API key and endpoint
func GetRateLimitConfig(apiKey, endpoint string) (limiter.RateLimitConfig, error) {
	db, err := database.GetInstance()
	if err != nil {
		return limiter.RateLimitConfig{}, err
	}

	// 1) Find project ID
	var projectID int
	err = db.Pool().QueryRow(context.Background(), "SELECT id FROM projects WHERE api_key=$1", apiKey).Scan(&projectID)
	if err != nil {
		return limiter.RateLimitConfig{}, ErrProjectNotFound
	}

	// 2) Check if endpoint belongs to project
	query := `SELECT strategy, key_by, limit_count, window_seconds
			  FROM rate_limit_rules
			  WHERE project_id = $1 AND endpoint = $2`
	row := db.Pool().QueryRow(context.Background(), query, projectID, endpoint)

	var cfg limiter.RateLimitConfig
	var strategy, keyBy string
	var limitCount, windowSeconds int
	err = row.Scan(&strategy, &keyBy, &limitCount, &windowSeconds)
	if err != nil {
		return limiter.RateLimitConfig{}, ErrEndpointNotOwned
	}

	cfg.Strategy = core.StrategyType(strategy)
	cfg.KeyBy = core.KeyFuncType(keyBy)
	cfg.Limit = limitCount
	cfg.Window = time.Duration(windowSeconds) * time.Second
	cfg.RedisCluster = limiter.RedisClusterConfig{
		Nodes: []string{
			getEnvOrDefault("REDIS_NODE_1", "redis://localhost:6379/0"),
			getEnvOrDefault("REDIS_NODE_2", "redis://localhost:6380/0"),
			getEnvOrDefault("REDIS_NODE_3", "redis://localhost:6381/0"),
		},
		Strategy: getEnvOrDefault("SHARDING_STRATEGY", "hash_mod"),
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
