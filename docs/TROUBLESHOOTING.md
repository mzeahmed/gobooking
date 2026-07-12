# Troubleshooting

Common issues encountered while running Gobooking locally, and how to fix them.

---

## `dial tcp 127.0.0.1:5432: connect: connection refused`

**Symptom**: an API request that touches the database fails with an error
like:

```text
failed to connect to `user=gobooking database=gobooking`:
    127.0.0.1:5432 (localhost): dial error: dial tcp 127.0.0.1:5432: connect: connection refused
    [::1]:5432 (localhost): dial error: dial tcp [::1]:5432: connect: connection refused
```

**Cause**: `internal/config/config.go` builds the database DSN with this
precedence: use `DATABASE_URL` if it is set, otherwise fall back to the
`DB_HOST` / `DB_PORT` / `DB_USER` / `DB_PASSWORD` / `DB_NAME` variables (see
`internal/db/db.go`).

The root `.env` file defines `DATABASE_URL` pointing at `localhost`, because
that variable is primarily meant for the `golang-migrate` CLI (`make
migrate-up` / `make migrate-down`), which runs on the host machine, outside
Docker.

Config loading uses [`godotenv.Load`](https://github.com/joho/godotenv),
which **never overrides a variable that is already set in the process
environment**. Inside the `app` container, `docker-compose.yml` sets
`DB_HOST=database` (and friends) directly as container environment
variables, so those are correctly picked up. But until this was fixed,
`DATABASE_URL` was not set in the container's environment, so `godotenv`
loaded it from `.env` — with `localhost` — and that value won.

From inside the `app` container, `localhost` refers to the container
itself, not the `database` service, hence the connection refusal.

**Fix**: `docker-compose.yml` now also sets `DATABASE_URL` in the `app`
service's environment, pointing at the `database` service host instead of
`localhost`. Docker Compose automatically loads the root `.env` file for
`${...}` interpolation inside `docker-compose.yml`, so the credentials are
sourced from there rather than duplicated — only the host is hardcoded,
since that's the one part that must differ from `.env`:

```yaml
environment:
  DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@database:${DB_PORT}/${DB_NAME}?sslmode=disable
```

This takes precedence over the `.env` value inside the container, while the
`.env` file's `localhost` DSN remains correct for `migrate` invoked from the
host.

**If you hit this again**: check whether the value actually being used
differs between "run via Docker" and "run via `migrate` / `go run` on the
host". Any variable that should differ between those two contexts needs to
be set explicitly in `docker-compose.yml`'s `environment:` block for the
`app` service, since container env vars always win over `.env`.

---

## `password authentication failed for user "'gobooking'"` (quotes included)

**Symptom**: a request that touches the database fails with a Postgres
auth error where the username/database name visibly include quote
characters, e.g.:

```text
failed to connect to `user='gobooking' database='gobooking'`:
    ... failed SASL auth: FATAL: password authentication failed for user "'gobooking'"
```

**Cause**: `.env` values were wrapped in single quotes (e.g.
`DB_USER='gobooking'`). Three different tools read this file, and only two
of them strip surrounding quotes:

* `godotenv` (Go) — strips quotes. Fine for `go run` / `air` reading `.env`
  directly.
* Docker Compose's built-in `.env` parser — strips quotes. Fine when
  Compose reads `.env` itself for `${...}` interpolation (e.g. via
  `docker compose config`).
* `make`'s `include .env` + `export` (used by every `make` target, see the
  top of the `Makefile`) — does **not** strip quotes. It exports
  `DB_USER` into the shell environment as the literal string
  `'gobooking'`, quotes included.

Docker Compose gives **shell environment variables precedence over its own
`.env` file parsing**. Since `make up` / `make restart` run `docker compose`
as a subprocess of `make` (which already exported the quoted values into
the shell), Compose uses the quoted shell values instead of its own
correctly-parsed ones — and the literal quotes leak into the container's
environment and then into the DSN.

This is why `docker compose config`, run directly, showed clean values
during debugging, while the actual running container (started via `make
restart`) had quotes — same `docker-compose.yml`, different invocation
path.

**Fix**: removed the surrounding single quotes from simple values in
`.env` (`DB_HOST`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `JWT_SECRET`).
Unquoted values parse identically across `godotenv`, Compose, and `make`.

**If you hit this again**: keep `.env` values unquoted unless a value
actually contains characters that require quoting (spaces, `#`, etc.) —
and if it does, verify it round-trips correctly through `make`'s
`include`/`export`, not just through Compose or `godotenv` alone.