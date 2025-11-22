-- name: UpsertContent :one
INSERT INTO contents (
    provider_id,
    provider_content_id,
    title,
    content_type,
    published_at,
    is_active,
    updated_at
) VALUES (
    sqlc.arg(provider_id),
    sqlc.arg(provider_content_id),
    sqlc.arg(title),
    sqlc.arg(content_type),
    sqlc.arg(published_at),
    sqlc.arg(is_active),
    NOW()
)
ON CONFLICT (provider_id, provider_content_id)
DO UPDATE SET
    title = EXCLUDED.title,
    content_type = EXCLUDED.content_type,
    published_at = EXCLUDED.published_at,
    is_active = EXCLUDED.is_active,
    updated_at = NOW()
RETURNING id;
-- name: SearchContents :many
SELECT 
    id,
    provider_id,
    provider_content_id,
    title,
    content_type,
    published_at,
    is_active,
    created_at,
    updated_at
FROM contents
WHERE
    is_active = true
    AND (sqlc.narg(query)::text IS NULL OR title ILIKE '%' || sqlc.narg(query)::text || '%')
    AND (sqlc.narg(content_type)::varchar IS NULL OR content_type = sqlc.narg(content_type)::varchar)
ORDER BY id DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: CountContents :one
SELECT COUNT(*)
FROM contents
WHERE
    is_active = true
    AND (sqlc.narg(query)::text IS NULL OR title ILIKE '%' || sqlc.narg(query)::text || '%')
    AND (sqlc.narg(content_type)::varchar IS NULL OR content_type = sqlc.narg(content_type)::varchar);


-- name: GetContentByID :one
SELECT 
    id,
    provider_id,
    provider_content_id,
    title,
    content_type,
    published_at,
    is_active,
    created_at,
    updated_at
FROM contents
WHERE id = sqlc.arg(content_id);

-- name: GetContentsByIDs :many
SELECT 
    id,
    provider_id,
    provider_content_id,
    title,
    content_type,
    published_at,
    is_active,
    created_at,
    updated_at
FROM contents
WHERE id = ANY(sqlc.arg(content_ids)::bigint[]);
