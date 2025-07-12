<p align="center">
  <img src="logo.png" width="140" alt="RLaaS logo" />
</p>

<h1 align="center">RLaaS — <em>Rate‑Limiting&nbsp;as&nbsp;a&nbsp;Service</em></h1>
<p align="center">
  Scalable, multi‑tenant rate‑limiting middleware you can drop into <strong>any</strong> stack.
</p>
<p align="center">
  <a href="https://go.dev/doc"><img src="https://img.shields.io/badge/Go-1.24%2B-00ADD8?logo=go&style=flat" alt="Go version" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-green?style=flat" alt="MIT license" /></a>
</p>

---

## ✨ Key Highlights

|                            |                                                                 |
| -------------------------- | --------------------------------------------------------------- |
| 🔐 **OAuth 2.0 Auth**      | Google OAuth + secure JWT session cookies                       |
| 🗝️ **Project & API Keys** | One‑click project creation, UUID v4 keys                        |
| 🕹️ **Fine‑grained Rules** | Token‑Bucket & Sliding‑Window strategies per endpoint           |
| 🚀 **Redis Sharding**      | Consistent‑Hash or Mod‑Hash selector across N nodes             |
| 📈 **Stateless Check API** | `POST /check` returns <code>{"allowed"\:true}</code> or **429** |
| 🐳 **Container‑First**     | Single `docker‑compose` spins up Postgres + 3×Redis + RLaaS     |

---

## 🏗️ Tech Stack & Architecture

<details>
<summary>Click to expand</summary>

```text
┌──────────────┐   Auth            ┌────────────┐           ████████  Redis 3‑node ring
│   Frontend   │──────────────────▶│  RLaaS API │──────────▶  Node A   (rate counters)
└──────────────┘   cookie (HTTP)   └────────────┘           ████████
                       ▲            ▲                       Node B
                       │  REST/JSON │                       ████████
                       │            └───▶  PostgreSQL       Node C
                       │                 (users, projects, rules)
                       └───────────────────────────────────────────
```

* **API layer** — <strong>Fiber</strong> + middlewares (auth, logging, recovery)
* **Persistence** — <strong>GORM</strong> + PostgreSQL
* **Limiter Core** — <strong>gorl</strong> + bespoke shard selector
* **Observability** — <code>log/slog</code> JSON logs ▪︎ Prom metrics *coming soon*

</details>

---

## 📂 Project Layout

```txt
rlaas/
├─ cmd/api/              # entrypoint (main.go) & DI wiring
├─ internal/
│  ├─ app/               # builder/bootstrapper
│  ├─ auth/              # Google login & logout handlers
│  ├─ check/             # /check endpoint (stateless limiter)
│  ├─ config/            # env loader → typed struct
│  ├─ database/          # GORM init + migrations
│  ├─ limiter/           # shard selector & limiter facade
│  ├─ logging/           # slog logger factory
│  ├─ middleware/        # auth, request logger, recovery
│  ├─ project/           # project domain (model, repo, service, handler)
│  ├─ rule/              # rule    domain (model, repo, service, handler)
│  └─ user/              # user    domain (model, repo, service, handler)
├─ docker-compose.yml    # Postgres + 3×Redis + RLaaS API
├─ Makefile              # build / run / test tasks
└─ README.md             # you are here
```

---

## ⚡ Quick Start (Docker)

```bash
# start Postgres + Redis‑cluster + RLaaS API
make docker-run   # or: docker compose up -d

# stop & clean
make docker-down
```

Local API will be available at **[http://localhost:8080](http://localhost:8080)**.

---

## 🛠 Local Development

```bash
cp .env.example .env   # configure DB creds & OAuth keys
make run               # go run + (optional) live reload
```

Set <code>MIGRATE\_ON\_START=true</code> to auto‑apply DB migrations on boot.

---

## 📑 REST API Reference

### Authentication

| Method | Path                    | Description                                       |
| ------ | ----------------------- | ------------------------------------------------- |
| `GET`  | `/auth/google/login`    | Redirects to Google consent screen                |
| `GET`  | `/auth/google/callback` | Completes OAuth flow, sets `session_token` cookie |
| `POST` | `/logout`               | Clears session cookie                             |

### User

\| `GET /me` | Returns <code>{id,email,name,avatar\_url}</code> (JWT required) |

### Projects

| Method   | Path             | Body / Params                  |
| -------- | ---------------- | ------------------------------ |
| `POST`   | `/projects`      | `{ "project_name": "My API" }` |
| `GET`    | `/projects`      | –                              |
| `DELETE` | `/projects/:pid` | –                              |

### Rules

| Method   | Path                        | Body      |
| -------- | --------------------------- | --------- |
| `GET`    | `/projects/:pid/rules`      | –         |
| `POST`   | `/projects/:pid/rules`      | see below |
| `PUT`    | `/projects/:pid/rules/:rid` | –         |
| `DELETE` | `/projects/:pid/rules/:rid` | –         |

```jsonc
{
  "endpoint": "/api/v1/resource",
  "strategy": "token_bucket",   // or "sliding_window"
  "key_by":   "ip",             // ip | api_key | user_id
  "limit_count": 100,
  "window_seconds": 60
}
```

### Rate‑Limit Check

```http
POST /check
{
  "api_key":   "<project-key>",
  "endpoint":  "/api/v1/resource",
  "key":       "client-ip or user-id"
}
```

*200* → `{ "allowed": true }`   |   *429* → **Too Many Requests**

---

## 🏃 Make Targets

| Target             | Purpose                 |
| ------------------ | ----------------------- |
| `make build`       | Compile RLaaS binary    |
| `make run`         | Run with live reload    |
| `make docker-run`  | Compose up all services |
| `make docker-down` | Stop & clean containers |

---

## 🔧 Configuration (.env)

| Key                         | Default                         | Description              |
| --------------------------- | ------------------------------- | ------------------------ |
| `PORT`                      | `8080`                          | HTTP listen port         |
| `DB_HOST` / …               | –                               | Postgres credentials     |
| `JWT_SECRET`                | –                               | HMAC secret for sessions |
| `GOOGLE_CLIENT_ID / SECRET` | –                               | OAuth 2.0 app creds      |
| `REDIS_NODE_1..3`           | `redis://localhost:6379/0` etc. | Redis shard URLs         |
| `SHARDING_STRATEGY`         | `hash_mod`                      | or `consistent_hash`     |
| `MIGRATE_ON_START`          | `false`                         | Auto‑migrate on boot     |


## License

Released under the **MIT License** – see [`LICENSE`](LICENSE) for full text.
