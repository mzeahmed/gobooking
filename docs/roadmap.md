# Roadmap

## Phase 1 - Foundation

* [x] Go module setup
* [x] Docker Compose (PostgreSQL, Mailpit)
* [x] golang-migrate setup
* [ ] sqlc setup
* [x] HTTP server bootstrap (net/http router, middleware, config loading)

---

## Phase 2 - Authentication

* [x] `user` schema + migration
* [ ] Registration
* [ ] Login
* [ ] Logout
* [ ] Email verification
* [ ] Password reset

---

## Phase 3 - Room Management

* [ ] `room` schema + migration
* [ ] `equipment` schema + migration
* [ ] Create room
* [ ] Edit room
* [ ] Delete room
* [ ] Amenities management (equipment assignment)
* [ ] Room status (draft / published)

---

## Phase 4 - Marketplace

* [ ] Room listing endpoint
* [ ] Search
* [ ] Filters (capacity, city, price range, equipment)
* [ ] Room details endpoint

---

## Phase 5 - Booking

* [ ] `reservation` schema + migration
* [ ] Availability service
* [ ] Reservation creation
* [ ] Overlap prevention (`EXCLUDE` constraint + service-level check)
* [ ] Cancellation workflow
* [ ] Booking confirmation (email)

---

## Phase 6 - Reviews

* [ ] `review` schema + migration
* [ ] Review creation (one per user per room)
* [ ] Ratings aggregation on room listing/details

---

## Phase 7 - Notifications

* [ ] Mailer integration (SMTP, Mailpit in dev)
* [ ] Async email sending

---

## Phase 8 - Administration

* [ ] Admin-only room management endpoints
* [ ] Admin-only user management endpoints
* [ ] Admin-only reservation management endpoints

---

## Phase 9 - Quality

* [ ] Unit tests for domain services
* [ ] Integration tests for HTTP handlers
* [ ] `golangci-lint` in CI
* [ ] CI/CD pipeline

---

## Phase 10 - Advanced

* [ ] Structured logging + request tracing
* [ ] Rate limiting
* [ ] Caching layer
* [ ] API documentation (OpenAPI)
