# Database Design

## Overview

Gobooking uses PostgreSQL as its primary database.

All entities use integer identifiers (`SERIAL` / `BIGSERIAL`), matching the
identifiers actually used by the original Symfony application's Doctrine
entities and migrations (the Symfony project's own docs originally
described UUIDs, but that was never implemented — this port follows the
real schema, not the stale docs).

Timestamps (`created_at`, `updated_at`) use `TIMESTAMPTZ` and are managed
by the application layer, mirroring the `#[PrePersist]` / `#[PreUpdate]`
lifecycle callbacks of the Doctrine entities.

## Status

Two groups of tables:

* **Implemented** in the source Symfony app today: `user`, `room`,
  `equipment`, `equipment_room`.
* **Planned** in the source app's domain model but not yet built there
  either (`reservation`, `review`). They are documented here so the Go
  schema is designed up front for the full domain, even though they will
  land in a later implementation phase — see [roadmap.md](roadmap.md).

## Entity Relationship Diagram

```text
User
 │
 ├── Reservations
 └── Reviews

Room
 │
 ├── Reservations
 ├── Reviews
 └── Equipments (many-to-many)

Reservation
 │
 ├── User
 └── Room

Review
 │
 ├── User
 └── Room
```

## Tables

### user

*Implemented.*

| Column                | Type         | Constraints                          |
|------------------------|--------------|---------------------------------------|
| id                     | BIGSERIAL    | PRIMARY KEY                           |
| email                  | VARCHAR(180) | NOT NULL, UNIQUE                      |
| password               | VARCHAR(255) | NOT NULL (bcrypt hash)                |
| reset_password_token   | VARCHAR(255) | NULL                                  |
| first_name             | VARCHAR(255) | NULL                                  |
| last_name              | VARCHAR(255) | NULL                                  |
| roles                  | JSONB        | NOT NULL, default `[]`                |
| is_verified            | BOOLEAN      | NOT NULL, default `false`             |
| created_at             | TIMESTAMPTZ  | NOT NULL                              |
| updated_at             | TIMESTAMPTZ  | NOT NULL                              |

`roles` stores a JSON array of role strings (e.g. `["ROLE_USER"]`,
`["ROLE_USER", "ROLE_ADMIN"]`), same convention as the Symfony
`UserInterface::getRoles()` contract. Every user implicitly has
`ROLE_USER` at the application layer, whether or not it's present in the
stored array.

---

### room

*Implemented.*

| Column       | Type          | Constraints    |
|--------------|---------------|-----------------|
| id           | BIGSERIAL     | PRIMARY KEY     |
| name         | VARCHAR(255)  | NOT NULL        |
| description  | TEXT          | NOT NULL        |
| capacity     | INT           | NOT NULL        |
| hourly_price | NUMERIC(10,2) | NOT NULL        |
| address      | VARCHAR(255)  | NULL            |
| city         | VARCHAR(255)  | NULL            |
| postal_code  | VARCHAR(255)  | NULL            |
| created_at   | TIMESTAMPTZ   | NOT NULL        |
| updated_at   | TIMESTAMPTZ   | NOT NULL        |

Two fixes versus the current Symfony schema, applied here since this is a
fresh implementation:

* `address` (the PHP column is misspelled `adress`).
* `hourly_price NUMERIC(10,2)` (the PHP column is `NUMERIC(10,0)`, i.e. no
  decimal places, which cannot represent cents — a bug in the source app).

---

### equipment

*Implemented.*

| Column      | Type         | Constraints |
|-------------|--------------|-------------|
| id          | BIGSERIAL    | PRIMARY KEY |
| name        | VARCHAR(255) | NOT NULL    |
| description | TEXT         | NULL        |
| icon        | VARCHAR(255) | NULL        |
| created_at  | TIMESTAMPTZ  | NOT NULL    |
| updated_at  | TIMESTAMPTZ  | NOT NULL    |

---

### equipment_room

*Implemented.* Join table for the `Room` ↔ `Equipment` many-to-many
relationship.

| Column       | Type      | Constraints                                      |
|--------------|-----------|----------------------------------------------------|
| equipment_id | BIGINT    | NOT NULL, FK → `equipment(id)` ON DELETE CASCADE   |
| room_id      | BIGINT    | NOT NULL, FK → `room(id)` ON DELETE CASCADE        |

Primary key: `(equipment_id, room_id)`.

---

### reservation

*Planned.*

| Column      | Type          | Constraints                                  |
|-------------|---------------|------------------------------------------------|
| id          | BIGSERIAL     | PRIMARY KEY                                    |
| room_id     | BIGINT        | NOT NULL, FK → `room(id)`                      |
| user_id     | BIGINT        | NOT NULL, FK → `user(id)`                      |
| start_date  | TIMESTAMPTZ   | NOT NULL                                       |
| end_date    | TIMESTAMPTZ   | NOT NULL                                       |
| total_price | NUMERIC(10,2) | NOT NULL                                       |
| status      | VARCHAR(20)   | NOT NULL, one of the values below              |
| created_at  | TIMESTAMPTZ   | NOT NULL                                       |
| updated_at  | TIMESTAMPTZ   | NOT NULL                                       |

`status` values:

* `pending`
* `confirmed`
* `cancelled`
* `completed`

Enforced with a `CHECK` constraint (Postgres has no native enum
requirement here, but a `CHECK (status IN (...))` keeps invalid states out
at the DB level, which the current MariaDB/Doctrine schema does not do).

---

### review

*Planned.*

| Column     | Type        | Constraints                     |
|------------|-------------|-----------------------------------|
| id         | BIGSERIAL   | PRIMARY KEY                        |
| room_id    | BIGINT      | NOT NULL, FK → `room(id)`          |
| user_id    | BIGINT      | NOT NULL, FK → `user(id)`          |
| rating     | SMALLINT    | NOT NULL, CHECK (rating BETWEEN 1 AND 5) |
| comment    | TEXT        | NULL                                |
| created_at | TIMESTAMPTZ | NOT NULL                            |
| updated_at | TIMESTAMPTZ | NOT NULL                            |

## Indexes

* `user.email` — unique index (login lookups)
* `reservation.room_id`
* `reservation.user_id`
* `reservation.start_date`, `reservation.end_date` (range queries for availability checks)
* `review.room_id`
* `review.user_id`

## Constraints

* One review per user per room: `UNIQUE (user_id, room_id)` on `review`.
* Reservation dates must not overlap for the same room. Not expressible as
  a plain SQL constraint across rows portably; enforced at the application
  layer, with an optional `EXCLUDE` constraint using the `btree_gist`
  extension as a defense-in-depth option:

  ```sql
  CREATE EXTENSION IF NOT EXISTS btree_gist;

  ALTER TABLE reservation
      ADD CONSTRAINT reservation_no_overlap
      EXCLUDE USING gist (
          room_id WITH =,
          tstzrange(start_date, end_date) WITH &&
      )
      WHERE (status IN ('pending', 'confirmed'));
  ```
* `reservation.end_date` must be after `start_date`: `CHECK (end_date > start_date)`.
* `review.rating` must be between 1 and 5: `CHECK (rating BETWEEN 1 AND 5)`.

## Migrations

Schema changes are managed with [goose](https://github.com/pressly/goose)
SQL migration files under `internal/adapters/postgresql/migrations/`, one file per change, with
`-- +goose Up` / `-- +goose Down` sections — no ORM-managed auto-migrations.
