-- name: GetAllEnabledProviders :many
SELECT id, name, code, format, base_url, is_enabled, created_at, updated_at
FROM providers
WHERE is_enabled = true;

-- name: GetProviderByCode :one
SELECT id, name, code, format, base_url, is_enabled, created_at, updated_at
FROM providers
WHERE code = sqlc.arg(code);

-- name: GetProviderByID :one
SELECT id, name, code, format, base_url, is_enabled, created_at, updated_at
FROM providers
WHERE id = sqlc.arg(provider_id);

-- name: UpsertProvider :exec
INSERT INTO providers (
    name,
    code,
    format,
    base_url,
    is_enabled
) VALUES (
    sqlc.arg(name),
    sqlc.arg(code),
    sqlc.arg(format),
    sqlc.arg(base_url),
    sqlc.arg(is_enabled)
)
ON CONFLICT (code)
DO UPDATE SET
    name = EXCLUDED.name,
    format = EXCLUDED.format,
    base_url = EXCLUDED.base_url,
    is_enabled = EXCLUDED.is_enabled,
    updated_at = NOW();

