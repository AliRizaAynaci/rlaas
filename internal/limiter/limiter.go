package limiter

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/AliRizaAynaci/gorl"
	"github.com/AliRizaAynaci/gorl/core"
)

// RateLimitConfig holds the options for rate limiting (strategy, limit, window, vs.) :contentReference[oaicite:0]{index=0}
type RateLimitConfig struct {
	Strategy     core.StrategyType
	KeyBy        core.KeyFuncType
	Limit        int
	Window       time.Duration
	RedisCluster RedisClusterConfig
}

type RedisClusterConfig struct {
	Nodes    []string `json:"nodes"`
	Strategy string   `json:"strategy"` // "hash_mod", "consistent_hash"
}

// ConfigKey uniquely identifies a limiter configuration including rate and window :contentReference[oaicite:1]{index=1}
type ConfigKey struct {
	ApiKey   string        // client’s API key
	Endpoint string        // requested endpoint
	ShardKey string        // the Redis shard URL
	Limit    int           // number of allowed requests per window
	Window   time.Duration // time window duration (e.g. 1m, 10s)
}

type Limiter struct {
	gorlLimiter core.Limiter
}

var (
	limiterCache  = make(map[ConfigKey]*Limiter)
	shardSelector *ShardSelector
	mu            sync.Mutex
)

// Initialize shard selector
func InitSharding() {
	nodes := []string{
		getEnvOrDefault("REDIS_NODE_1", "redis://localhost:6379/0"),
		getEnvOrDefault("REDIS_NODE_2", "redis://localhost:6380/0"),
		getEnvOrDefault("REDIS_NODE_3", "redis://localhost:6381/0"),
	}

	// Boş node'ları filtrele
	var validNodes []string
	for _, node := range nodes {
		if node != "" {
			validNodes = append(validNodes, node)
		}
	}

	strategy := getEnvOrDefault("SHARDING_STRATEGY", "hash_mod")
	shardSelector = NewShardSelector(validNodes, strategy)
}

func GetLimiterForKey(apiKey, endpoint, userKey string, baseConfig RateLimitConfig) (*Limiter, error) {
	if shardSelector == nil {
		InitSharding()
	}

	shardKey := fmt.Sprintf("%s:%s", apiKey, endpoint)

	redisURL := shardSelector.GetRedisURL(shardKey)

	cfgKey := ConfigKey{
		ApiKey:   apiKey,
		Endpoint: endpoint,
		ShardKey: redisURL,
		Limit:    baseConfig.Limit,
		Window:   baseConfig.Window,
	}

	mu.Lock()
	defer mu.Unlock()

	if limiter, ok := limiterCache[cfgKey]; ok {
		return limiter, nil
	}

	gl, err := gorl.New(core.Config{
		Strategy: baseConfig.Strategy,
		KeyBy:    baseConfig.KeyBy,
		Limit:    baseConfig.Limit,
		Window:   baseConfig.Window,
		RedisURL: os.Getenv("REDISCLOUD_URL"),
	})
	if err != nil {
		return nil, err
	}

	limiter := &Limiter{gorlLimiter: gl}
	limiterCache[cfgKey] = limiter
	return limiter, nil
}

func (l *Limiter) Allow(key string) (bool, error) {
	return l.gorlLimiter.Allow(key)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
