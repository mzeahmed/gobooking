# Gobooking

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](go.mod)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://www.docker.com)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](#license)

A room and workspace booking platform, written in Go.

This project is a from-scratch Go rewrite of the
[Symfony BookingApp](https://github.com/mzeahmed/bookingapp), reusing the same domain and data
model, but idiomatic to Go: no ORM magic, explicit SQL, small interfaces,
standard library first.

---

## Project Goals

Same product scope as the original app:

Users can:

* Register and authenticate
* Browse available rooms
* Manage bookings
* Leave reviews
* Manage their profile

Administrators can:

* Manage rooms and equipment
* Manage bookings
* Manage users

This is also a learning project, aimed at practicing idiomatic Go backend
development: `net/http`, `sqlc`, structured concurrency, testing,
and clean layering without pulling in a full framework.

---

## Tech Stack

### Backend

* Go 1.23+
* `net/http` (stdlib) — HTTP router / middleware
* [sqlc](https://sqlc.dev) — typed Go code generated from SQL, on top of `database/sql`
* [goose](https://github.com/pressly/goose) — versioned SQL migrations
* [pgx](https://github.com/jackc/pgx) — PostgreSQL driver

### Infrastructure

* Docker Compose
* Traefik v3
* PostgreSQL 16
* Adminer (dev DB browser)
* Mailpit (dev email capture)

### Quality

* `go vet` / `staticcheck`
* `golangci-lint`
* `go test` (standard library testing, `testify` for assertions)

> The exact module layout and dependency choices will be finalized once
> implementation starts; this README documents the intended stack.

---

## Data Model

The Go version keeps the **same domain model** as the Symfony project:
`User`, `Room`, `Equipment`, `Reservation`, `Review`.

Two deliberate differences from the original PHP implementation:

* **Database engine**: PostgreSQL instead of MariaDB.
* **Identifiers**: integer primary keys (`SERIAL` / `BIGSERIAL`), matching
  what the PHP entities actually use today (not the UUIDs originally
  planned in the Symfony project's docs, which were never implemented).

See [docs/database.md](docs/database.md) for the full schema.

---

## Local Domains

The application uses HTTPS locally, the same way as the source Symfony
project.

Available domains:

```text
https://api.gobooking.local
https://mail.gobooking.local
https://db.gobooking.local
```

Traefik sits in front of the stack as a reverse proxy. It:

* Terminates HTTPS using local certificates generated with mkcert.
* Redirects all HTTP traffic to HTTPS.
* Routes each domain to the right container based on Host rules declared
  in `traefik/dynamic.yml` (e.g. `api.gobooking.local` -> app,
  `mail.gobooking.local` -> Mailpit, `db.gobooking.local` ->
  Adminer).

These domains do not resolve on their own: you need to edit `/etc/hosts`
to point them to `127.0.0.1` (see `make hosts` below).

Distinct domains (`*.gobooking.local` rather than `bookingapp.local`)
are used deliberately, so this stack can run alongside the original
Symfony project without `/etc/hosts` or port collisions.

---

## Requirements

* Go 1.23+
* Docker
* Docker Compose
* GNU Make
* [mkcert](https://github.com/FiloSottile/mkcert)
* [goose CLI](https://github.com/pressly/goose#install) (optional, for manual migration runs)
* [sqlc CLI](https://docs.sqlc.dev/en/latest/overview/install.html)

---

## Installation

### Clone the repository

```bash
git clone git@github.com:your-org/gobooking.git

cd gobooking
```

### Configure Hosts

```bash
make hosts
```

This adds the required domains to `/etc/hosts` (asks for your `sudo`
password) and is safe to run multiple times: entries already present are
left untouched.

### Start the Stack

```bash
make up
```

This generates local TLS certificates via mkcert (if missing), then
builds and starts Traefik, the app, PostgreSQL, Adminer, and Mailpit.

Application:

```text
https://api.gobooking.local
```

Traefik Dashboard:

```text
http://localhost:8080
```

Adminer:

```text
https://db.gobooking.local
```

Mailpit:

```text
https://mail.gobooking.local
```

> `make up` builds the `app` image from `Dockerfile`, which
> expects `api/go.mod` and a `cmd/api` entrypoint. Until the Go module and
> entrypoint exist (Phase 1 of the roadmap), only run `make up
> database mailpit adminer traefik` or comment out the `app` service.

### Run the Application Outside Docker

```bash
docker compose up -d database
go run ./cmd/api
```

---

## Make Commands

### Display Available Commands

```bash
make help
```

| Command          | Description                                |
|------------------|--------------------------------------------|
| `make run`       | Run the server locally                     |
| `make build`     | Build the binary into `api/bin/bookingapp` |
| `make fmt`       | Format the source code                     |
| `make vet`       | Run `go vet`                               |
| `make test`      | Run unit tests                             |
| `make check`     | Run `fmt`, `vet`, and `test`               |
| `make tidy`      | Clean up `go.mod` / `go.sum`               |
| `make update`    | Update dependencies                        |
| `make migrate-up`   | Apply database migrations                  |
| `make migrate-down` | Roll back the last migration               |
| `make sqlc`      | Regenerate Go code from SQL queries        |
| `make hosts`     | Add local domains to `/etc/hosts` (sudo)   |
| `make certs`     | Generate local TLS certificates if missing |
| `make up`        | Build and start the Docker containers      |
| `make down`      | Stop the Docker containers                 |
| `make restart`   | Restart the Docker containers              |
| `make logs`      | Show container logs                        |
| `make ps`        | List containers                            |
| `make bash`      | Open a shell in the app container          |
| `make clean`     | Remove generated files                     |
| `make doctor`    | Show the development environment           |

---

## Development Workflow

### Run Tests

```bash
make test
```

### Run Linter

```bash
golangci-lint run
```

### Create a New Migration

```bash
make migrate-create t="<migration_name>"
```

### Create a New Feature

1. Create a feature branch.
2. Write/update SQL queries under `api/db/queries/` and run `make sqlc`.
3. Implement the handler/service.
4. Add tests.
5. Run `make check`.
6. Submit a pull request.

---

## Project Structure

Planned layout, following standard Go project conventions:

```text
gobooking/

├── docker/
│   └── app/                 # Dockerfile for the Go application
│
├── traefik/
│   ├── traefik.yml           # static config (entrypoints, providers)
│   └── dynamic.yml            # routers/services, TLS
│
├── certs/                   # local mkcert certificates (gitignored)
│
├── docs/
│
├── api/                     # Go module (not yet implemented)
│   ├── cmd/
│   │   └── api/               # main package, application entrypoint
│   │
│   ├── internal/
│   │   ├── booking/            # reservation domain logic
│   │   ├── room/                 # room + equipment domain logic
│   │   ├── user/                  # user domain logic
│   │   ├── review/                # review domain logic
│   │   ├── http/                   # net/http router, handlers, middleware
│   │   └── db/                      # sqlc-generated code
│   │
│   ├── db/
│   │   ├── migrations/          # goose SQL migrations
│   │   └── queries/                # sqlc SQL query definitions
│   │
│   ├── sqlc.yaml
│   ├── go.mod
│   └── go.sum
│
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## Documentation

* [docs/architecture.md](docs/architecture.md) — layering and infrastructure overview
* [docs/database.md](docs/database.md) — schema, relationships, constraints
* [docs/roadmap.md](docs/roadmap.md) — implementation phases
* [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) — common local setup issues and fixes

---

## License

MIT
