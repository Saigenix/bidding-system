# Contributing to Bidding System SDK

First off, thank you for considering contributing! Every contribution helps make this project better for everyone.

---

## 💬 Join the Community

Have questions, ideas, or just want to chat? **Join our Discord:**

👉 **[Discord Server](https://discord.gg/7c9R2ttESZ)** 👈

---

## How Can I Contribute?

### 🐛 Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates.

**When filing a bug, include:**

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior vs actual behavior
- Your Go version (`go version`)
- Your OS and PostgreSQL version
- Relevant logs or error messages

### 💡 Suggesting Features

Open an issue with the `enhancement` label and describe:

- The problem you're trying to solve
- Your proposed solution
- Alternative approaches you've considered

### 🔧 Pull Requests

1. **Fork** the repository
2. **Clone** your fork:
   ```bash
   git clone https://github.com/<your-username>/bidding-system.git
   cd bidding-system
   ```
3. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. **Make your changes** following the coding guidelines below
5. **Run checks:**
   ```bash
   make lint
   make test
   ```
6. **Commit** with a meaningful message:
   ```bash
   git commit -m "feat: add support for auction categories"
   ```
7. **Push** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
8. Open a **Pull Request** against `main`

---

## Coding Guidelines

### Project Structure

This project follows **Clean Architecture**. Please place your code in the correct layer:

| Layer | Directory | Purpose |
|-------|-----------|---------|
| Domain | `internal/domain/` | Entities, interfaces, business rules |
| Repository | `internal/repository/` | Database implementations |
| Service | `internal/service/` | Use cases, business logic |
| Handler | `internal/handler/` | HTTP handlers, request/response |
| SDK | `sdk/` | Public SDK interface |

### Code Style

- Run `make fmt` before committing
- Run `make vet` to check for issues
- Follow standard Go conventions ([Effective Go](https://go.dev/doc/effective_go))
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: fix a bug
docs: documentation changes
refactor: code refactoring
test: add or update tests
chore: maintenance tasks
```

### Testing

- Write unit tests for service layer logic
- Use interfaces for mocking dependencies
- Place tests in `_test.go` files alongside the code
- Run `make test-cover` to check coverage

---

## Development Setup

### Prerequisites

- Go 1.23+
- PostgreSQL 14+ or Docker
- `make` (optional but recommended)

### Quick Start

```bash
# Clone and setup
git clone https://github.com/saigenix/bidding-system.git
cd bidding-system

# Full setup (installs deps, starts DB, runs migrations)
make setup

# Start development server
make run
```

### Environment Variables

Copy and configure:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=bidding
export JWT_SECRET=dev-secret-key
export SERVER_PORT=8080
export LOG_LEVEL=debug
```

---

## Need Help?

- 💬 **Discord**: [Join our server](https://discord.gg/7c9R2ttESZ)
- 📝 **Issues**: [Open an issue](https://github.com/saigenix/bidding-system/issues)

We're happy to help with any questions!

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
