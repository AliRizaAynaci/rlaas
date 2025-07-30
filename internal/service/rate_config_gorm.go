package service

import (
	"os"
	"time"

	"github.com/AliRizaAynaci/gorl/core"
	"github.com/AliRizaAynaci/rlaas/internal/limiter"
	"github.com/AliRizaAynaci/rlaas/internal/rule"
	"gorm.io/gorm"
)

type RateConfigService struct{ db *gorm.DB }

func NewRateConfigService(db *gorm.DB) *RateConfigService { return &RateConfigService{db} }

func (s *RateConfigService) Get(apiKey, endpoint string) (limiter.RateLimitConfig, error) {
	/* 1) project id */
	var pid uint
	if err := s.db.Raw(`SELECT id FROM projects WHERE api_key = ?`, apiKey).
		Scan(&pid).Error; err != nil || pid == 0 {
		return limiter.RateLimitConfig{}, ErrProjectNotFound
	}

	/* 2) rule */
	var rl rule.Rule
	if err := s.db.Where("project_id=? AND endpoint=?", pid, endpoint).
		First(&rl).Error; err != nil {
		return limiter.RateLimitConfig{}, ErrEndpointNotOwned
	}

	return limiter.RateLimitConfig{
		Strategy: core.StrategyType(rl.Strategy),
		KeyBy:    core.KeyFuncType(rl.KeyBy),
		Limit:    rl.LimitCount,
		Window:   time.Duration(rl.WindowSeconds) * time.Second,
		RedisCluster: limiter.RedisClusterConfig{
			Nodes: []string{
				getEnvOrDefault("REDIS_NODE_1", "redis://localhost:6379/0"),
				getEnvOrDefault("REDIS_NODE_2", "redis://localhost:6380/0"),
				getEnvOrDefault("REDIS_NODE_3", "redis://localhost:6381/0"),
			},
			Strategy: getEnvOrDefault("SHARDING_STRATEGY", "hash_mod"),
		},
		FailOpen: rl.FailOpen,
	}, nil
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
