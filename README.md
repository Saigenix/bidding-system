<p align="center">
  <h1 align="center">⚡ Bidding System SDK</h1>
  <p align="center">
    A production-ready, pluggable real-time bidding engine built with Go.
    <br />
    Run as a standalone server or embed as a library in your platform.
    <br /><br />
    <a href="docs/ARCHITECTURE.md"><strong>Architecture »</strong></a>
    ·
    <a href="docs/PROJECT_INFO.md"><strong>Project Info »</strong></a>
    ·
    <a href="CONTRIBUTING.md"><strong>Contributing »</strong></a>
    ·
    <a href="https://discord.gg/7c9R2ttESZ"><strong>Discord »</strong></a>
  </p>
</p>

---

## What is this?

The Bidding System SDK lets you **add auctions and real-time bidding** to any platform. It handles:

- 🔐 **User auth** (JWT + bcrypt)
- 📦 **Product management**
- 🏷️ **Auction lifecycle** (create → start → bid → end)
- ⚡ **Real-time bid streaming** via REST, SSE, and WebSocket
- 🧩 **Pluggable SDK** — use as a Go library or standalone server
- 📖 **Swagger API docs** — interactive API documentation at `/swagger/index.html`

---

## Quick Start

```bash
# Clone
git clone https://github.com/saigenix/bidding-system.git
cd bidding-system

# First-time setup (installs deps, starts DB, runs migrations)
make setup

# Start the server
make run
```

Server runs at `http://localhost:8080`. See all commands with `make help`.

---

## Docker

The fastest way to get a fully working stack (PostgreSQL + migrations + app):

```bash
# Start everything
docker compose up -d

# Follow app logs
docker compose logs -f app

# Stop and clean up
docker compose down -v
```

This spins up:

| Container | Description | Port |
|-----------|-------------|------|
| `bidding-postgres` | PostgreSQL 16 (Alpine) | `5432` |
| `bidding-migrate` | Runs schema migrations then exits | — |
| `bidding-app` | Bidding System server | `8080` |

The app won't start until PostgreSQL is healthy and migrations have completed.

### Build image only

```bash
# Build
docker build -t bidding-system .

# Run (requires external PostgreSQL)
docker run -p 8080:8080 --env-file .env bidding-system
```

> **Tip:** Copy `.env.example` to `.env` and adjust values before running.

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.23+ |
| Framework | [Gin](https://github.com/gin-gonic/gin) |
| Database | PostgreSQL ([pgx](https://github.com/jackc/pgx)) |
| Auth | JWT + bcrypt |
| Real-time | SSE + WebSocket |
| Config | Viper (env-based) |
| Logging | Zerolog |
| API Docs | [Swagger](https://github.com/swaggo/gin-swagger) (OpenAPI 2.0) |
| Testing | Go standard `testing` + mock repositories |

Full details → [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

---

## Project Structure

```
├── cmd/server/          → Standalone server
├── config/              → Configuration (Viper)
├── internal/
│   ├── domain/          → Entities & interfaces (clean core)
│   │   └── *_test.go    → Domain unit tests
│   ├── repository/      → PostgreSQL adapters
│   ├── service/         → Business logic
│   │   └── *_test.go    → Service unit tests
│   ├── handler/         → REST + SSE + WebSocket handlers (Swagger annotated)
│   ├── auth/            → JWT middleware
│   └── mocks/           → Mock repository implementations for testing
├── pkg/                 → Shared packages (db, logger, router)
├── sdk/                 → Public SDK interface
├── migrations/          → SQL schema migrations
├── deploy/
│   └── helm/            → Helm chart for Kubernetes
├── docs/
│   ├── swagger/         → Generated Swagger/OpenAPI docs
│   ├── ARCHITECTURE.md  → Architecture deep dive
│   └── PROJECT_INFO.md  → LLM context file
├── .github/workflows/   → CI pipeline (lint, test, build)
├── Dockerfile           → Multi-stage production build
├── docker-compose.yml   → Full local stack (Postgres + app)
└── .env.example         → Environment variable template
```

---

## Swagger / API Docs

Interactive API documentation is available at:

```
http://localhost:8080/swagger/index.html
```

To regenerate Swagger docs after modifying handler annotations:

```bash
make swagger
```

---

## Testing

Run the full test suite:

```bash
# Run all tests
make test

# Run tests with coverage report
make test-cover
```

The project uses **mock repositories** (`internal/mocks/`) for unit testing the service and domain layers without requiring a database connection. Tests cover:

- **Domain layer** — Auction state helpers (`IsActive`, `HasEnded`)
- **Auth service** — Registration, login, JWT token validation
- **Product service** — CRUD operations and error handling
- **Auction service** — Lifecycle transitions and validation rules
- **Bid service** — Bid placement, amount validation, winning bid

---

## API Endpoints

### Auth (Public)

```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Login → returns JWT
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Products (Protected — pass `Authorization: Bearer <TOKEN>`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/products` | Create product |
| `GET` | `/products` | List products |
| `GET` | `/products/:id` | Get product |

### Auctions

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auctions` | Create auction |
| `GET` | `/auctions` | List auctions |
| `GET` | `/auctions/:id` | Get auction |
| `POST` | `/auctions/:id/start` | Start auction |
| `POST` | `/auctions/:id/end` | End auction |

### Bids & Real-Time

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auctions/:id/bids` | Place bid |
| `GET` | `/auctions/:id/bids` | Get all bids |
| `GET` | `/auctions/:id/bids/stream` | **SSE** live stream |
| `WS` | `/auctions/:id/bids/ws` | **WebSocket** |

> 📖 For full request/response schemas, visit the [Swagger UI](http://localhost:8080/swagger/index.html).

---

## Using as an SDK

Embed bidding into your Go application:

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/saigenix/bidding-system/sdk"
)

func main() {
    engine, err := sdk.NewEngine(
        sdk.WithJWTSecret("my-secret"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Stop()

    ctx := context.Background()

    // Create product → auction → bid
    product, _ := engine.CreateProduct(ctx, "Laptop", "Gaming laptop", "user-id")
    auction, _ := engine.CreateAuction(ctx, product.ID, time.Now(), time.Now().Add(24*time.Hour), 500.00)
    bid, _ := engine.PlaceBid(ctx, auction.ID, "bidder-id", 600.00)

    log.Printf("Bid placed: %+v", bid)
}
```

---

## Configuration

Set via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Server port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | DB user |
| `DB_PASSWORD` | `password` | DB password |
| `DB_NAME` | `bidding` | Database name |
| `JWT_SECRET` | — | **Set in production** |
| `JWT_EXPIRATION_HOUR` | `24` | Token lifetime |
| `LOG_LEVEL` | `info` | debug/info/warn/error |

---

## Make Commands

```
make run              Run the server
make build            Build binary
make test             Run tests
make test-cover       Run tests with coverage
make lint             Format + vet
make fmt-check        Check formatting (CI)
make ci               Run full CI pipeline locally
make swagger          Generate Swagger docs
make swagger-install  Install swag CLI
make db-up            Start PostgreSQL (Docker)
make db-down          Stop PostgreSQL
make migrate-up       Apply migrations
make migrate-down     Rollback migrations
make docker-build     Build Docker image
make docker-run       Run Docker container
make setup            Full first-time setup
make help             Show all commands
```

---

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**Join our community →** [Discord](https://discord.gg/7c9R2ttESZ)

---

## License

MIT License — see [LICENSE](LICENSE) for details.
