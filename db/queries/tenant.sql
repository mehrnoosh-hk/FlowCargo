-- name: CreateTenant :one
INSERT INTO tenants (
    name,
    email,
    domain
) VALUES (
    $1, $2, $3
) RETURNING id, name, email, domain, is_active, created_at, updated_at;

-- name: GetTenantByID :one
SELECT id, name, email, domain, is_active, created_at, updated_at FROM tenants
WHERE id = $1;

-- name: GetTenantByEmail :one
SELECT id, name, email, domain, is_active, created_at, updated_at FROM tenants
WHERE email = $1;

-- name: GetTenantByDomain :one
SELECT id, name, email, domain, is_active, created_at, updated_at FROM tenants
WHERE domain = $1;

-- name: UpdateTenant :one
UPDATE tenants
SET
    name = CASE
        WHEN sqlc.narg('name')::VARCHAR IS NOT NULL THEN sqlc.narg('name')
        ELSE name
    END,
    email = CASE
        WHEN sqlc.narg('email')::VARCHAR IS NOT NULL THEN sqlc.narg('email')
        ELSE email
    END,
    domain = CASE
        WHEN sqlc.narg('domain')::VARCHAR IS NOT NULL THEN sqlc.narg('domain')
        ELSE domain
    END,
    is_active = CASE
        WHEN sqlc.narg('is_active')::BOOL IS NOT NULL THEN sqlc.narg('is_active')
        ELSE is_active
    END
WHERE id = $1
RETURNING id, name, email, domain, is_active, created_at, updated_at;

-- name: DeleteTenant :exec
DELETE FROM tenants
WHERE id = $1::UUID;

-- name: ListTenants :many
SELECT id, name, email, domain, is_active, created_at, updated_at FROM tenants
ORDER BY created_at DESC;
