# Architecture & Tech Stack

## Overview

The Bidding System SDK is built following **Clean Architecture** (also known as Hexagonal / Ports & Adapters). The core principle: **business logic has zero knowledge of infrastructure**. Repositories are defined as interfaces in the domain layer, and concrete implementations (PostgreSQL) live in the adapter layer.

---

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────┐
│                        cmd/server/main.go                        │
│                     (Application Entry Point)                    │
└──────────────────────────┬───────────────────────────────────────┘
                           │ wires
┌──────────────────────────▼───────────────────────────────────────┐
│                          sdk/engine.go                            │
│              (SDK Engine — Dependency Injection Root)             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ AuthService  │  │AuctionSvc   │  │  BidService  │  ...        │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└──────────────────────────┬───────────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────────┐
        │                  │                      │
┌───────▼──────┐  ┌────────▼────────┐  ┌──────────▼──────────┐
│  Transport   │  │  Service Layer  │  │  Repository Layer   │
│   (Gin)      │  │  (Use Cases)    │  │  (PostgreSQL)       │
│              │  │                 │  │                     │
│ • REST API   │  │ • Auth logic    │  │ • UserRepo          │
│ • SSE Stream │  │ • Bid rules     │  │ • ProductRepo       │
│ • WebSocket  │  │ • Auction mgmt  │  │ • AuctionRepo       │
│ • JWT MW     │  │ • Product CRUD  │  │ • BidRepo           │
└──────────────┘  └────────┬────────┘  └──────────┬──────────┘
                           │                      │
                    ┌──────▼──────────────────────▼────────┐
                    │           Domain Layer                │
                    │                                      │
                    │  Entities:   User, Product,           │
                    │              Auction, Bid             │
                    │  Interfaces: UserRepository,          │
                    │              ProductRepository,       │
                    │              AuctionRepository,       │
                    │              BidRepository            │
                    └──────────────────────────────────────┘
```

---

## Layer Responsibilities

### 1. Domain Layer (`internal/domain/`)

The **innermost layer**. Contains:

- **Entities** — Pure data structures with no external dependencies
  - `User` — ID, Email, PasswordHash, CreatedAt
  - `Product` — ID, Name, Description, OwnerID, CreatedAt
  - `Auction` — ID, ProductID, StartTime, EndTime, StartingPrice, CurrentPrice, Status, CreatedAt
  - `Bid` — ID, AuctionID, UserID, Amount, CreatedAt

- **Repository Interfaces (Ports)** — Contracts that the outer layers must implement
  - `UserRepository` — Create, GetByEmail, GetByID
  - `ProductRepository` — Create, GetByID, List
  - `AuctionRepository` — Create, GetByID, List, Update
  - `BidRepository` — Create, GetByAuctionID, GetHighestBid

**Rule**: This layer imports nothing from the project. It defines the language of the system.

### 2. Repository Layer (`internal/repository/postgres/`)

Implements the domain interfaces using PostgreSQL (`pgx/pgxpool`).

- Uses parameterized queries (SQL injection safe)
- Connection pooling via `pgxpool`
- Each repo gets a `*pgxpool.Pool` via constructor injection

**Swappable**: To use MySQL or MongoDB, create `internal/repository/mysql/` implementing the same interfaces.

### 3. Service Layer (`internal/service/`)

Business logic and use cases:

| Service | Responsibilities |
|---------|-----------------|
| `AuthService` | Register (bcrypt), Login (JWT), Token validation |
| `ProductService` | Create, Get, List products |
| `AuctionService` | Create, Get, List, Start, End auctions with validation |
| `BidService` | Place bids with amount/status validation, update auction price |

**Rule**: Services only depend on domain interfaces, never on concrete repos.

### 4. Transport Layer (`internal/handler/`, `internal/auth/`, `pkg/web/`)

HTTP layer using Gin:

- **Handlers** — Parse HTTP requests, call services, return JSON responses
- **Middleware** — JWT authentication (extract token → validate → set userID in context)
- **Router** — Route registration, CORS, endpoint grouping

### 5. SDK Layer (`sdk/`)

Public interface for embedding:

- `Engine` — Main struct, holds all services, provides convenience methods
- `Option` — Functional options (WithDBPool, WithLogger, WithJWTSecret, WithConfig)
- Users can `import "github.com/saigenix/bidding-system/sdk"` and use the engine directly

### 6. Infrastructure (`config/`, `pkg/`)

Cross-cutting concerns:

| Package | Technology | Purpose |
|---------|-----------|---------|
| `config/` | Viper | Environment-based configuration |
| `pkg/db/` | pgxpool | PostgreSQL connection pool management |
| `pkg/logger/` | Zerolog | Structured logging with levels |
| `pkg/web/` | Gin | Router setup and CORS |

---

## Tech Stack

| Category | Technology | Package |
|----------|-----------|---------|
| **Language** | Go 1.23+ | — |
| **HTTP Framework** | Gin | `github.com/gin-gonic/gin` |
| **Database** | PostgreSQL 14+ | `github.com/jackc/pgx/v5/pgxpool` |
| **Authentication** | JWT | `github.com/golang-jwt/jwt/v5` |
| **Password Hashing** | bcrypt | `golang.org/x/crypto/bcrypt` |
| **ID Generation** | UUID v4 | `github.com/google/uuid` |
| **Configuration** | Viper | `github.com/spf13/viper` |
| **Logging** | Zerolog | `github.com/rs/zerolog` |
| **Migrations** | golang-migrate | `github.com/golang-migrate/migrate/v4` |
| **WebSocket** | Gorilla WebSocket | `github.com/gorilla/websocket` |
| **SSE** | Gin built-in | `c.SSEvent()` |

---

## Communication Protocols

### REST API

Standard JSON over HTTP for CRUD operations.

```
POST   /auth/register          — Register user
POST   /auth/login             — Login, get JWT
POST   /products               — Create product
GET    /products                — List products
GET    /products/:id            — Get product
POST   /auctions               — Create auction
GET    /auctions                — List auctions
GET    /auctions/:id            — Get auction
POST   /auctions/:id/start     — Start auction
POST   /auctions/:id/end       — End auction
POST   /auctions/:id/bids      — Place bid
GET    /auctions/:id/bids       — Get bids
```

### Server-Sent Events (SSE)

One-way server→client stream for live bid updates.

```
GET /auctions/:auction_id/bids/stream

Response: text/event-stream
event: bid
data: {"id":"...","user_id":"...","amount":150.00,"created_at":"..."}
```

### WebSocket

Bi-directional real-time communication for interactive bidding.

```
WS /auctions/:auction_id/bids/ws

Server sends:
  {"type":"initial","bids":[...]}
  {"type":"update","latest_bid":{...}}
```

---

## Database Schema

```sql
users
├── id           UUID (PK)
├── email        VARCHAR(255) UNIQUE
├── password_hash VARCHAR(255)
└── created_at   TIMESTAMPTZ

products
├── id           UUID (PK)
├── name         VARCHAR(255)
├── description  TEXT
├── owner_id     UUID (FK → users)
└── created_at   TIMESTAMPTZ

auctions
├── id             UUID (PK)
├── product_id     UUID (FK → products)
├── start_time     TIMESTAMPTZ
├── end_time       TIMESTAMPTZ
├── starting_price DECIMAL(10,2)
├── current_price  DECIMAL(10,2)
├── status         VARCHAR(20) [pending|active|ended]
└── created_at     TIMESTAMPTZ

bids
├── id          UUID (PK)
├── auction_id  UUID (FK → auctions)
├── user_id     UUID (FK → users)
├── amount      DECIMAL(10,2)
└── created_at  TIMESTAMPTZ
```

### Indexes

- `idx_products_owner` — products(owner_id)
- `idx_auctions_product` — auctions(product_id)
- `idx_auctions_status` — auctions(status)
- `idx_bids_auction` — bids(auction_id)
- `idx_bids_user` — bids(user_id)
- `idx_bids_created` — bids(created_at DESC)

---

## Dependency Flow

```
main.go → sdk.Engine → services → domain interfaces ← repositories → PostgreSQL
                ↓
          handlers → Gin router → HTTP
```

**Key rule**: Arrows point inward. Outer layers depend on inner layers, never the reverse.
