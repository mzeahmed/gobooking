-- name: CreateUser :one
INSERT INTO users (email, password, first_name, last_name, is_verified)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at;

-- name: AssignDefaultRole :exec
INSERT INTO user_roles (user_id, role_id)
SELECT $1, id
FROM roles
WHERE name = $2;

-- name: FindUserByEmail :one
SELECT u.id, u.email, u.password,
       COALESCE(u.first_name, '')::text AS first_name,
       COALESCE(u.last_name, '')::text  AS last_name,
       u.is_verified, u.created_at, u.updated_at,
       COALESCE(array_agg(r.name) FILTER (WHERE r.name IS NOT NULL), '{}')::text[] AS roles
FROM users u
         LEFT JOIN user_roles ur ON ur.user_id = u.id
         LEFT JOIN roles r ON r.id = ur.role_id
WHERE u.email = $1
GROUP BY u.id;

-- name: ListUsers :many
SELECT u.id, u.email, u.password,
       COALESCE(u.first_name, '')::text AS first_name,
       COALESCE(u.last_name, '')::text  AS last_name,
       u.is_verified, u.created_at, u.updated_at,
       COALESCE(array_agg(r.name) FILTER (WHERE r.name IS NOT NULL), '{}')::text[] AS roles
FROM users u
         LEFT JOIN user_roles ur ON ur.user_id = u.id
         LEFT JOIN roles r ON r.id = ur.role_id
GROUP BY u.id
ORDER BY u.id;

-- name: DeleteUser :execrows
DELETE
FROM users
WHERE id = $1;