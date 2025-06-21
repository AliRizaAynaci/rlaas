# Rate Limiting as a Service (RLaaS)

A scalable and flexible rate limiting service that provides API rate limiting capabilities as a service. This project allows you to implement rate limiting for your applications without managing the complexity of rate limiting infrastructure. Built with Go and utilizing the `gorl` library for efficient rate limiting operations.

## Features

- **Dynamic Rate Limiting Rules**: Create and manage rate limiting rules per project and endpoint
- **Flexible Rate Limiting Strategies**: Support for various rate limiting strategies via `gorl`
- **Distributed Redis Sharding**: 
  - Multiple Redis node support
  - Configurable sharding strategies (hash_mod, consistent_hash)
  - High availability and scalability
- **RESTful API**: Simple and intuitive API for managing rate limits
- **Real-time Rate Limit Checking**: Fast and efficient rate limit verification
- **Project-based Management**: Organize rate limits by projects
- **Custom Key Generation**: Flexible key generation for rate limiting (IP, user ID, custom keys)

## Architecture

The service is built using:
- Go (Golang) for the backend service
- `gorl` library for rate limiting implementation
- Redis for distributed rate limiting with sharding support
- PostgreSQL for storing configuration and rules
- Docker for containerization

## Project Structure

```
rlaas/
├── cmd/api/           # Application entry point
├── internal/         
│   ├── database/      # Database operations
│   │   ├── database.go
│   │   └── database_test.go
│   ├── limiter/       # Rate limiting core
│   │   ├── limiter.go    # Main rate limiting logic
│   │   └── shard.go      # Redis sharding implementation
│   ├── models/        # Data models
│   │   ├── project.go    # Project entity
│   │   └── rule.go       # Rate limit rules
│   ├── server/        # HTTP server and handlers
│   │   ├── handlers/     # Request handlers
│   │   ├── routes.go     # API routes
│   │   └── server.go     # Server configuration
│   └── service/       # Business logic layer
│       ├── apikey.go     # API key management
│       ├── config.go     # Configuration
│       ├── project.go    # Project logic
│       └── rule.go       # Rule management
├── docker-compose.yml # Docker composition
└── Makefile          # Build and development commands
```

## Prerequisites

- Go 1.19 or higher
- Docker and Docker Compose
- Redis (for distributed rate limiting)
- PostgreSQL (for configuration storage)


## API Endpoints

### Base Endpoints
- `GET /` - Service health check
- `GET /health` - Database health status

### Project Management
- `POST /register` - Register a new project
  ```json
  {
    "name": "project_name",
    "api_key": "your_api_key"
  }
  ```

### Rate Limit Rules
- `GET /rules` - List all rules
- `POST /rule/add` - Create a new rate limit rule
  ```json
  {
    "project_id": 1,
    "endpoint": "/api/resource",
    "strategy": "fixed_window",
    "key_by": "ip",
    "limit_count": 100,
    "window_seconds": 3600
  }
  ```
- `PUT /rule/` - Update a rule
- `DELETE /rule/` - Delete a rule

### Rate Limit Checking
- `POST /check` - Check if a request is within rate limits
  ```json
  {
    "api_key": "your_api_key",
    "endpoint": "/api/resource",
    "key": "127.0.0.1"
  }
  ```

## Rate Limiting Configuration

The service supports various rate limiting configurations:

```go
type RateLimitConfig struct {
    Strategy     core.StrategyType    // Rate limiting strategy
    KeyBy        core.KeyFuncType     // Key generation method
    Limit        int                  // Rate limit count
    Window       time.Duration        // Time window
    RedisCluster RedisClusterConfig   // Redis configuration
}
```

### Supported Strategies
- Fixed Window
- Sliding Window
- Token Bucket
- Leaky Bucket

### Key Generation Methods
- IP Address
- User ID
- Custom Keys
- Combined Keys


## Redis Sharding

The service supports two sharding strategies:
1. **Hash Modulo**: Simple distribution using modulo operation
2. **Consistent Hashing**: More balanced distribution with minimal redistribution

Configure sharding via environment variables:
```bash
SHARDING_STRATEGY=consistent_hash
REDIS_NODE_1=redis://localhost:6379/0
REDIS_NODE_2=redis://localhost:6380/0
REDIS_NODE_3=redis://localhost:6381/0
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
