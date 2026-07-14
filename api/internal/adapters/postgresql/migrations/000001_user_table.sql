-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id          INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email       VARCHAR(180) NOT NULL UNIQUE,
    password    TEXT         NOT NULL,
    first_name  VARCHAR(255),
    last_name   VARCHAR(255),
    is_verified BOOLEAN      NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE users IS 'Application user accounts.';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd