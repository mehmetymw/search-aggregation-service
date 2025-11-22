-- name: GetAllContentTypeMetadata :many
SELECT id, display_name, is_enabled, sort_order, created_at
FROM content_type_metadata
WHERE is_enabled = true
ORDER BY sort_order;

-- name: GetContentTypeMetadataByID :one
SELECT id, display_name, is_enabled, sort_order, created_at
FROM content_type_metadata
WHERE id = $1;
