# Architecture

## Overview

Gobooking is a room and workspace reservation platform, a Go port of
the Symfony BookingApp project. It reuses the same domain model and
product scope, implemented with explicit, idiomatic Go instead of a
full-stack framework.

## High Level Architecture

```text
Browser
    │
    ▼
Traefik
    │
    ▼
Go HTTP server (net/http + chi)
    │
 ┌──┴──────────────┐
 ▼                 ▼

PostgreSQL       Mailpit (dev)
```

No PHP-FPM/nginx layer is needed: the Go binary serves HTTP directly, and
Traefik sits in front of it as a reverse proxy, the same way it does in
the Symfony project. Traefik:

* Terminates HTTPS using local certificates generated with mkcert.
* Redirects all HTTP traffic to HTTPS.
* Routes `api.gobooking.local`, `mail.gobooking.local`, and
  `db.gobooking.local` to the app, Mailpit, and Adminer containers
  respectively, via `traefik/dynamic.yml`.

See [database.md](database.md) for schema details and the project
[README](../README.md) for the full Docker/Traefik setup instructions.

## Application Layers

```text
HTTP handler (chi)
    │
    ▼
Request decoding / validation
    │
    ▼
Domain service
    │
    ▼
sqlc-generated repository (database/sql over pgx)
    │
    ▼
PostgreSQL
```

Each domain package (`internal/user`, `internal/room`, `internal/booking`,
`internal/review`) owns its own handlers, service logic, and SQL queries.
There is no shared "God" repository layer — packages depend on narrow
interfaces, not on each other's concrete types.

## Core Modules

### Authentication

Responsible for:

* Registration
* Login (session or token based — to be decided at implementation time)
* Email verification
* Password reset

### Room Management

Responsible for:

* Room creation, update, deletion
* Equipment management

### Reservations

Responsible for:

* Availability checks
* Reservation creation
* Reservation cancellation

### Reviews

Responsible for:

* User reviews
* Ratings

### Notifications

Responsible for:

* Transactional emails (via Mailpit in development)

## Authorization

Role-based, same roles as the source app:

* `ROLE_USER`
* `ROLE_MANAGER`
* `ROLE_ADMIN`

Enforced with `chi` middleware that reads the authenticated user's roles
and gates access to handlers — the equivalent of Symfony Voters, without
the framework machinery.

## Infrastructure

### PostgreSQL

Main relational database. See [database.md](database.md) for schema.

### Mailpit

Used during development to inspect outgoing emails, same role as in the
Symfony project.

### Adminer

Used during development to browse the PostgreSQL database, the equivalent
of phpMyAdmin in the Symfony project.

## Testing Strategy

### Unit Tests

Focus on:

* Domain services
* Business rules (availability checks, pricing, validation)

### Integration Tests

Focus on:

* HTTP handlers, using `httptest`
* Repository queries, against a real PostgreSQL instance (e.g. via
  `testcontainers-go` or a dedicated test database)

## Future Improvements

* OpenAPI documentation generated from handler definitions
* Structured background job processing (async email sending, etc.)
* Observability: structured logging, metrics, tracing
