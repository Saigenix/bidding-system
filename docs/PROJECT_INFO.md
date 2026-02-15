# Project Info — LLM Context File

> This file provides structured context about the Bidding System SDK project for LLMs
> assisting with development. Read this before making any changes.

---

## Project Identity

- **Name**: Bidding System SDK
- **Module**: `github.com/saigenix/bidding-system`
- **Language**: Go 1.23+
- **Type**: Library (SDK) + Standalone Server
- **License**: MIT

## What This Project Does

A pluggable, real-time bidding system that can be:
1. **Run as a standalone server** via `cmd/server/main.go`
2. **Embedded as a Go SDK** via `import "github.com/saigenix/bidding-system/sdk"`

Core user flows:
- Register/Login → Get JWT token
- Create product → Create auction with time window → Start auction
- Place bids (must exceed current price, auction must be active)
- Stream bid updates via SSE or WebSocket
- End auction → Determine winner (highest bid)

---

## Directory Map

```
bidding-system/
│
├── cmd/server/main.go              ← Entry point. Wires SDK + Gin router.
│
├── config/config.go                ← Viper config. Reads from env vars.
│
├── internal/
│   ├── domain/                     ← CORE. Entities + repository interfaces.
│   │   ├── user.go                   User entity
│   │   ├── product.go                Product entity
│   │   ├── auction.go                Auction entity (has Status enum + helpers)
│   │   ├── auction_test.go           Tests for IsActive/HasEnded helpers
│   │   ├── bid.go                    Bid entity
│   │   └── repository.go            All repository interfaces (ports)
│   │
│   ├── repository/postgres/        ← PostgreSQL implementations of domain interfaces.
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── auction.go
│   │   └── bid.go
│   │
│   ├── service/                    ← Business logic. Depends ONLY on domain interfaces.
│   │   ├── auth.go                   JWT generation, bcrypt, login/register
│   │   ├── auth_test.go              Auth service unit tests
│   │   ├── product.go                Product CRUD
│   │   ├── product_test.go           Product service unit tests
│   │   ├── auction.go                Auction lifecycle (create/start/end)
│   │   ├── auction_test.go           Auction service unit tests
│   │   ├── bid.go                    Bid placement with validation
│   │   └── bid_test.go               Bid service unit tests
│   │
│   ├── handler/                    ← Gin HTTP handlers (Swagger annotated).
│   │   ├── auth.go                   POST /auth/register, POST /auth/login
│   │   ├── product.go                GET/POST /products
│   │   ├── auction.go                GET/POST /auctions, start/end
│   │   └── bid.go                    POST bids, SSE stream, WebSocket
│   │
│   ├── mocks/                      ← Mock repositories for unit testing.
│   │   └── repositories.go          In-memory implementations of all repo interfaces
│   │
│   └── auth/middleware.go          ← JWT middleware for Gin
│
├── pkg/
│   ├── db/db.go                    ← pgxpool connection manager
│   ├── logger/logger.go            ← Zerolog logger factory
│   └── web/router.go               ← Gin router setup (routes, CORS, Swagger, middleware)
│
├── sdk/
│   ├── engine.go                   ← Public SDK. Creates repos + services. Has convenience methods.
│   └── options.go                  ← Functional options (WithDBPool, WithLogger, etc.)
│
├── migrations/
│   ├── 001_init.up.sql             ← Creates users, products, auctions, bids tables
│   └── 001_init.down.sql           ← Drops all tables
│
├── docs/
│   ├── swagger/                    ← Generated Swagger/OpenAPI documentation
│   │   ├── docs.go                   Go package for embedding Swagger specs
│   │   ├── swagger.json              OpenAPI 2.0 JSON spec
│   │   └── swagger.yaml              OpenAPI 2.0 YAML spec
│   ├── ARCHITECTURE.md             ← Architecture deep dive
│   └── PROJECT_INFO.md             ← This file
│
├── Makefile                        ← make run, make build, make test, make swagger, etc.
├── CONTRIBUTING.md                 ← Contribution guidelines
└── README.md                       ← User-facing documentation
```

---

## Key Patterns

### Clean Architecture
- Domain layer has ZERO imports from the project
- Services depend on interfaces, not concrete types
- Repository implementations are injected via constructors

### Dependency Injection
- All wiring happens in `sdk/engine.go` (`NewEngine`)
- No global state, no init() magic
- Functional options pattern for configuration

### Auction State Machine
```
pending → active → ended
```
- `pending`: Created but not yet open for bidding
- `active`: Open for bidding, bids must exceed current_price
- `ended`: No more bids accepted

### Bid Validation Rules
1. Auction must be in `active` status
2. Current time must be between start_time and end_time
3. Bid amount must be strictly greater than auction.current_price
4. After placing bid, auction.current_price is updated

### Testing Strategy
- **Mock repositories** (`internal/mocks/`) — in-memory implementations of all repository interfaces
- **Service-layer tests** — test business logic without a database, using mock repos
- **Domain-layer tests** — test entity helper methods (e.g., `IsActive`, `HasEnded`)
- Run: `make test` or `go test ./... -v`

---

## Important Dependencies

| Package | Import Path | Used For |
|---------|-------------|----------|
| Gin | `github.com/gin-gonic/gin` | HTTP routing, middleware, SSE |
| pgx | `github.com/jackc/pgx/v5/pgxpool` | PostgreSQL driver + pool |
| JWT | `github.com/golang-jwt/jwt/v5` | Token generation/validation |
| bcrypt | `golang.org/x/crypto/bcrypt` | Password hashing |
| UUID | `github.com/google/uuid` | Entity ID generation |
| Viper | `github.com/spf13/viper` | Config from env vars |
| Zerolog | `github.com/rs/zerolog` | Structured logging |
| WebSocket | `github.com/gorilla/websocket` | Real-time bidding |
| Migrate | `github.com/golang-migrate/migrate/v4` | DB schema migrations |
| Swagger | `github.com/swaggo/gin-swagger` | OpenAPI docs generation & UI |
| Swag | `github.com/swaggo/swag` | Swagger annotation parser |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | PostgreSQL user |
| `DB_PASSWORD` | `password` | PostgreSQL password |
| `DB_NAME` | `bidding` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `DB_MAX_CONNS` | `25` | Max pool connections |
| `DB_MIN_CONNS` | `5` | Min pool connections |
| `JWT_SECRET` | `your-secret-key...` | HMAC signing key |
| `JWT_EXPIRATION_HOUR` | `24` | Token TTL in hours |
| `LOG_LEVEL` | `info` | debug/info/warn/error |

---

## Common Tasks

```bash
make run            # Start dev server
make build          # Compile binary
make test           # Run tests
make test-cover     # Run tests with coverage report
make lint           # Format + vet
make swagger        # Regenerate Swagger docs
make db-up          # Start PostgreSQL in Docker
make migrate-up     # Apply migrations
make setup          # Full first-time setup
```

---

## API Quick Reference

**Swagger UI**: `http://localhost:8080/swagger/index.html`

**Public:**
- `POST /auth/register` — `{"email":"...","password":"..."}`
- `POST /auth/login` — `{"email":"...","password":"..."}` → `{"token":"..."}`

**Protected (Bearer token):**
- `POST /products` — `{"name":"...","description":"..."}`
- `GET /products`, `GET /products/:id`
- `POST /auctions` — `{"product_id":"...","start_time":"...","end_time":"...","starting_price":100}`
- `GET /auctions`, `GET /auctions/:id`
- `POST /auctions/:id/start`, `POST /auctions/:id/end`
- `POST /auctions/:id/bids` — `{"auction_id":"...","amount":150}`
- `GET /auctions/:id/bids`
- `GET /auctions/:id/bids/stream` — SSE
- `GET /auctions/:id/bids/ws` — WebSocket

---

## Known Limitations / TODOs

- SSE and WebSocket currently use polling (2s interval) — needs pub/sub (Redis/NATS)
- No background worker to auto-end expired auctions
- No rate limiting
- No integration tests (only service + domain unit tests)
- WebSocket `CheckOrigin` allows all origins (restrict in production)
- No pagination on list endpoints
- Product ownership is not validated when creating auctions
