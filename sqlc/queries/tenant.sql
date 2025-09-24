-- name: CreateTenant :one
INSERT INTO tenants (
    name,
    email,
    domain
) VALUES (
    @name, @email, @domain
) RETURNING id, name, email, domain, is_active, created_at, updated_at;

-- name: GetTenantByID :one
SELECT id, name, email, domain, is_active, created_at, updated_at FROM tenants
WHERE id = @id;

-- name: GetTenantByEmail :one
SELECT id, name, email, domain, is_active, created_at, updated_at FROM tenants
WHERE email = @email;

-- name: GetTenantByDomain :one
SELECT id, name, email, domain, is_active, created_at, updated_at FROM tenants
WHERE domain = @domain;

-- name: UpdateTenant :one
UPDATE tenants
SET
    name = COALESCE(sqlc.narg(name), name),
    email = COALESCE(sqlc.narg(email), email),
    domain = COALESCE(sqlc.narg(domain), domain),
    is_active = COALESCE(sqlc.narg(is_active), is_active)
WHERE id = @id
RETURNING id, name, email, domain, is_active, created_at, updated_at;

-- name: DeleteTenant :exec
DELETE FROM tenants
WHERE id = @id;

-- name: ListTenants :many
SELECT id, name, email, domain, is_active, created_at, updated_at FROM tenants
ORDER BY created_at DESC;
