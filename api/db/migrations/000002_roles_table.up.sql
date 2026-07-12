CREATE TABLE roles
(
    id          INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name        VARCHAR(180) NOT NULL UNIQUE,
    description TEXT         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
) ;

COMMENT ON TABLE roles IS 'User roles available in the application.';

INSERT INTO roles ( name, description)
VALUES ('admin', 'Administrator role with full access'),
       ('user', 'Regular user role with limited access'),
       ('moderator', 'Moderator role with permissions to manage content'),
       ('guest', 'Guest role with minimal access');
