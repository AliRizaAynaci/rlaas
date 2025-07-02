# RLaaS: Rate Limiting as a Service

**RLaaS** is an enterprise‑grade, multi‑tenant rate limiting service built with Go, Redis, and PostgreSQL. It provides your applications with fine‑grained, scalable API rate limiting without the overhead of managing complex infrastructure.

---

## Table of Contents

1. [Key Features](#key-features)
2. [Database Schema](#database-schema)
3. [Architecture Overview](#architecture-overview)
4. [Project Structure](#project-structure)
5. [Getting Started](#getting-started)

  * [Prerequisites](#prerequisites)
  * [Clone & Install](#clone--install)
  * [Environment Configuration](#environment-configuration)
  * [Start Dependencies](#start-dependencies)
  * [Database Migrations](#database-migrations)
  * [Launch the Server](#launch-the-server)
6. [API Reference](#api-reference)

  * [Authentication](#authentication)
  * [Project Management](#project-management)
  * [Rate Limit Rules](#rate-limit-rules)
  * [Rate Limit Check](#rate-limit-check)
7. [Local Testing Examples](#local-testing-examples)
8. [Contributing](#contributing)
9. [License](#license)

---

## Key Features

* **Multi‑Tenant**: Isolate rate limits by project and user via API keys and Google OAuth2.
* **Dynamic Rule Management**: Create, list, update, and delete rate limit rules at runtime.
* **Flexible Strategies**: Support for fixed window, sliding window, token bucket, and leaky bucket via the `gorl` library.
* **Distributed Redis Sharding**: Horizontally scale rate limiting across multiple Redis nodes using hash\_mod or consistent hashing.
* **RESTful API**: Intuitive endpoints for administration and real‑time rate limit checks.
* **JWT‑Based Security**: Secure user authentication with Google OAuth2 and JWT.

---

## Database Schema

<details>
<summary>Click&nbsp;to&nbsp;expand&nbsp;🎛️</summary>

<table>
<tr>
<td valign="top" width="33%">

### users

| column      | type         | constraints      |
| ----------- | ------------ | ---------------- |
| id          | serial       | PK               |
| google\_id  | varchar(64)  | UNIQUE, NOT NULL |
| email       | varchar(128) | UNIQUE, NOT NULL |
| name        | varchar(128) | NOT NULL         |
| created\_at | timestamptz  | DEFAULT now()    |

</td>
<td valign="top" width="33%">

### projects

| column      | type        | constraints                      |
| ----------- | ----------- | -------------------------------- |
| id          | serial      | PK                               |
| user\_id    | int         | FK → users.id, ON DELETE CASCADE |
| name        | varchar(64) | UNIQUE per user, NOT NULL        |
| api\_key    | char(64)    | UNIQUE, NOT NULL                 |
| created\_at | timestamptz | DEFAULT now()                    |

</td>
<td valign="top" width="34%">

### rate\_limit\_rules

| column          | type         | constraints                     |
| --------------- | ------------ | ------------------------------- |
| id              | serial       | PK                              |
| project\_id     | int          | FK → projects.id                |
| endpoint        | varchar(128) | NOT NULL                        |
| strategy        | varchar(32)  | NOT NULL (fixed, sliding, etc.) |
| key\_by         | varchar(32)  | NOT NULL (api\_key, ip, user)   |
| limit\_count    | int          | NOT NULL                        |
| window\_seconds | int          | NOT NULL                        |
| created\_at     | timestamptz  | DEFAULT now()                   |

</td>
</tr>
</table>

</details>

---

## Architecture Overview

Below is a simplified system design sketch using ASCII art:

```text
      +--------------+     +--------------------------+
      |              |     |                          |
      |    Client    +----->   HTTP Server & Router   |
      |  (HTTP/CLI)  |     |      (Go net/http)       |
      |              |     |                          |
      +--------------+     +---------+----------------+
                                     |
                                     |
        +----------------------------v----------------------------+
        |                     Middleware Layer                    |
        |  - Google OAuth2 / JWT Authentication                   |
        |  - Request Logging & Metrics (Prometheus)               |
        +----------------------------+----------------------------+
                                     |
                                     |
          +--------------------------v------------------------+   
          |                   Service Core                    |   
          |  - Project & Rule Management                      |
          |  - RateLimit Check Handler                        |
          |  - Health & Admin Endpoints                       |
          +--------------------------+------------------------+   
                                     |
                                     |
                +--------------------v-----------------+      
                |           Limiter Module             |      
                | (gorl strategies + Redis Sharding)   |      
                +--------------------+-----------------+      
                                     |
                                     |                      
                +--------------------v-----------------+      
                |      Data Stores & Configuration     |      
                |  +---------------------------+       |      
                |  | PostgreSQL (Users, Rules) |       |      
                |  +---------------------------+       |      
                |  +---------------------------+       |      
                |  | Redis Cluster (Counters)  |       |      
                |  +---------------------------+       |      
                +--------------------------------------+      
```

This layout reflects a monolithic Go application with clear layers rather than a microservices mesh.

---

## Project Structure

```rlaas/
├── cmd/api/           # Application entry point
├── internal/         
│   ├── database/      # Database operations
│   │   ├── database.go
│   │   └── database_test.go
|   |── middleware/
|   |   └── auth.go       # Google OAuth2 operations
│   ├── limiter/       # Rate limiting core
│   │   ├── limiter.go    # Main rate limiting logic
│   │   └── shard.go      # Redis sharding implementation
│   ├── models/        # Data models
│   │   ├── project.go    # Project entity
│   │   ├── rule.go       # Rate limit rules
│   │   └── user.go       # User entity
│   ├── server/        # HTTP server and handlers
│   │   ├── handlers/     # Request handlers
│   │   ├── routes.go     # API routes
│   │   └── server.go     # Server configuration
│   └── service/       # Business logic layer
│       ├── apikey.go     # API key management
│       ├── auth.go       # DB operations 
│       ├── config.go     # Configuration
│       ├── project.go    # Project logic
│       └── rule.go       # Rule management
├── docker-compose.yml # Docker composition
└── Makefile          # Build and development commands
```
---

## Getting Started

> *The steps below assume a Unix‑like environment; adjust commands for Windows as needed.*

### Prerequisites

* **Go 1.19+**
* **Docker & Docker Compose**
* **migrate CLI**:

```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

* **Google OAuth2 Credentials** (Client ID & Secret)

### Clone & Install

```bash
git clone https://github.com/AliRizaAynaci/rlaas.git
cd rlaas
go mod tidy
```

### Environment Configuration

1. Copy the example:

   ```bash
   cp .env.example .env
   ```
2. Edit `.env` and provide values:

```dotenv
PORT=8080
APP_ENV=local

DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=rlaas
DB_USERNAME=postgres
DB_PASSWORD=password
DB_SCHEMA=public

SHARDING_STRATEGY=hash_mod
REDIS_NODE_1=redis://localhost:6379/0
REDIS_NODE_2=redis://localhost:6380/0
REDIS_NODE_3=redis://localhost:6381/0

GOOGLE_CLIENT_ID=<YOUR_GOOGLE_CLIENT_ID>
GOOGLE_CLIENT_SECRET=<YOUR_GOOGLE_CLIENT_SECRET>
OAUTH_REDIRECT_URL=http://localhost:8080/auth/google/callback

JWT_SECRET=<YOUR_BASE64_URL_SAFE_SECRET>
```

### Start Dependencies

Bring up PostgreSQL and Redis cluster:

```bash
docker-compose up -d
```

### Database Migrations

Set the `DATABASE_URL` environment variable and run migrations:

```bash
export DATABASE_URL="postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_DATABASE?sslmode=disable&search_path=$DB_SCHEMA"
migrate -path migrations -database "$DATABASE_URL" up
```

*On Windows (PowerShell):*

```powershell
$env:DATABASE_URL = "postgres://$($env:DB_USERNAME):$($env:DB_PASSWORD)@$($env:DB_HOST):$($env:DB_PORT)/$($env:DB_DATABASE)?sslmode=disable&search_path=$($env:DB_SCHEMA)"
migrate -path migrations -database $env:DATABASE_URL up
```

### Launch the Server

```bash
go run cmd/api/main.go
```

The API server will listen on `http://localhost:8080`.

---

## API Reference

---

### 🔑 Authentication

#### `GET /auth/google/login`

**Description:** Redirects users to Google OAuth2 consent screen to initiate authentication.

**Response:**

* **302 Found**: Redirects client to Google OAuth2 URL.

```http
GET /auth/google/login HTTP/1.1
Host: localhost:8080
```

---

#### `GET /auth/google/callback`

**Description:** Handles the OAuth2 callback from Google, validates the code, creates or retrieves the user, and issues a JWT session cookie.

**Query Parameters:**

| Name  | Type   | Required | Description               |
| ----- | ------ | -------- | ------------------------- |
| code  | string | yes      | OAuth2 authorization code |
| state | string | no       | CSRF protection token     |

**Responses:**

* **200 OK**: Sets `session_token` cookie; body contains user info.
* **400 Bad Request**: Missing or invalid code.
* **500 Internal Server Error**: OAuth or DB error.

```http
GET /auth/google/callback?code=abc123 HTTP/1.1
Host: localhost:8080
```

---

### 🗂 Project Management

All endpoints below require a valid JWT `session_token` cookie.

#### `POST /register`

Create a new project and receive a unique API key.

**Request Body:**

| Field         | Type   | Required | Description                 |
| ------------- | ------ | -------- | --------------------------- |
| project\_name | string | yes      | Unique name for the project |

```http
POST /register HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Cookie: session_token=<JWT>

{
  "project_name": "my-project"
}
```

**Responses:**

* **200 OK**

  ```json
  {
    "api_key": "e3f1a2..."
  }
  ```
* **400 Bad Request**: Invalid name or missing field.
* **409 Conflict**: Project name already exists.

---

### 📏 Rate Limit Rules

All rule operations require the header `Authorization: Bearer <API_KEY>`.

#### `GET /rules`

List all rate-limit rules for the current project.

**Headers:**

```
Authorization: Bearer <API_KEY>
```

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "endpoint": "/api/v1/check",
    "strategy": "fixed_window",
    "key_by": "api_key",
    "limit_count": 10,
    "window_seconds": 60
  }
]
```

---

#### `POST /rule/add`

Add a new rule under your project.

**Request Body:**

| Field           | Type   | Required | Description                            |
| --------------- | ------ | -------- | -------------------------------------- |
| endpoint        | string | yes      | URI path to be rate‑limited            |
| strategy        | string | yes      | `fixed_window`, `sliding_window`, etc. |
| key\_by         | string | yes      | `api_key`, `ip`, or custom identifier  |
| limit\_count    | int    | yes      | Allowed requests per window            |
| window\_seconds | int    | yes      | Window duration in seconds             |

```http
POST /rule/add HTTP/1.1
Host: localhost:8080
Authorization: Bearer <API_KEY>
Content-Type: application/json

{
  "endpoint": "/api/v1/check",
  "strategy": "fixed_window",
  "key_by": "api_key",
  "limit_count": 10,
  "window_seconds": 60
}
```

**Responses:**

* **201 Created**: Returns the created rule object.
* **400 Bad Request**: Validation error.
* **401 Unauthorized**: Missing/invalid API key.

---

#### `PUT /rule/{id}`

Update attributes of an existing rule.

**Path Parameter:**

| Name | Type | Description     |
| ---- | ---- | --------------- |
| id   | int  | Rule identifier |

**Request Body:** (at least one field)

```json
{ "limit_count": 20 }
```

**Responses:**

* **200 OK**: `{ "message": "Rule updated successfully" }`
* **400 Bad Request**: Invalid update payload.
* **404 Not Found**: Rule does not exist.

---

#### `DELETE /rule/{id}`

Remove a rate-limit rule.

**Path Parameter:**

| Name | Type | Description     |
| ---- | ---- | --------------- |
| id   | int  | Rule identifier |

**Response:**

* **200 OK**: `{ "message": "Rule deleted successfully" }`
* **404 Not Found**: Rule not found.

---

### ✅ Rate Limit Check

Public endpoint to verify if a request falls within defined limits.

#### `POST /check`

**Request Body:**

| Field    | Type   | Required | Description                      |
| -------- | ------ | -------- | -------------------------------- |
| api\_key | string | yes      | API key for the project          |
| endpoint | string | yes      | Target endpoint (e.g., `/check`) |
| key      | string | yes      | Identifier (e.g., user ID or IP) |

```http
POST /check HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "api_key": "<API_KEY>",
  "endpoint": "/api/v1/check",
  "key": "user-123"
}
```

**Responses:**

* **200 OK**: `{ "allowed": true }`
* **429 Too Many Requests**: `{ "allowed": false, "retry_after": 45 }`

---

## Local Testing Examples

Use Insomnia, Postman, or `curl` to exercise the full flow:

1. **Login & get JWT**:

   ```bash
   curl -v http://localhost:8080/auth/google/login
   ```
2. **Register a project** (with JWT cookie)

   ```bash
   curl -X POST http://localhost:8080/register \
     -H "Cookie: session_token=<JWT>" \
     -H "Content-Type: application/json" \
     -d '{"project_name":"test"}'
   ```
3. **Manage rules** (with API key)

   ```bash
   curl -X POST http://localhost:8080/rule/add \
     -H "Authorization: Bearer <API_KEY>" \
     -H "Content-Type: application/json" \
     -d '{"endpoint":"/api/v1/check","strategy":"fixed_window","key_by":"api_key","limit_count":5,"window_seconds":60}'
   ```
4. **Rate limit check**:

   ```bash
   curl -X POST http://localhost:8080/check \
     -H "Content-Type: application/json" \
     -d '{"api_key":"<API_KEY>","endpoint":"/api/v1/check","key":"user-1"}'
   ```

---

## Contributing

Thank you for your interest! To contribute:

1. Fork the repo and create a branch: `git checkout -b feature/foo`
2. Implement your feature or fix.
3. Write tests and ensure existing tests pass: `make test`
4. Submit a Pull Request with a clear description.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
