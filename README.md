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

Full details → [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

---

## Project Structure

```
├── cmd/server/          → Standalone server
├── config/              → Configuration (Viper)
├── internal/
│   ├── domain/          → Entities & interfaces (clean core)
│   ├── repository/      → PostgreSQL adapters
│   ├── service/         → Business logic
│   ├── handler/         → REST + SSE + WebSocket handlers
│   └── auth/            → JWT middleware
├── pkg/                 → Shared packages (db, logger, router)
├── sdk/                 → Public SDK interface
├── migrations/          → SQL schema migrations
└── docs/                → Architecture & project docs
```

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
make lint             Format + vet
make db-up            Start PostgreSQL (Docker)
make db-down          Stop PostgreSQL
make migrate-up       Apply migrations
make migrate-down     Rollback migrations
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
