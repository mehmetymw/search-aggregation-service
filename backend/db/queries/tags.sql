-- name: EnsureTag :one
INSERT INTO tags (name)
VALUES (sqlc.arg(name))
ON CONFLICT (name)
DO UPDATE SET name = EXCLUDED.name
RETURNING id, name;

-- name: AssignTagToContent :exec
INSERT INTO content_tags (content_id, tag_id)
VALUES (sqlc.arg(content_id), sqlc.arg(tag_id))
ON CONFLICT (content_id, tag_id) DO NOTHING;

-- name: GetTagsByContentID :many
SELECT t.id, t.name
FROM tags t
INNER JOIN content_tags ct ON ct.tag_id = t.id
WHERE ct.content_id = sqlc.arg(content_id);

-- name: RemoveContentTags :exec
DELETE FROM content_tags
WHERE content_id = sqlc.arg(content_id);
