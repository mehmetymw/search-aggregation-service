-- name: UpsertContentStats :exec
INSERT INTO content_stats (
    content_id,
    views,
    likes,
    duration_sec,
    reading_time,
    reactions,
    comments,
    last_sync_at
) VALUES (
    sqlc.arg(content_id),
    sqlc.arg(views),
    sqlc.arg(likes),
    sqlc.arg(duration_sec),
    sqlc.arg(reading_time),
    sqlc.arg(reactions),
    sqlc.arg(comments),
    NOW()
)
ON CONFLICT (content_id)
DO UPDATE SET
    views = EXCLUDED.views,
    likes = EXCLUDED.likes,
    duration_sec = EXCLUDED.duration_sec,
    reading_time = EXCLUDED.reading_time,
    reactions = EXCLUDED.reactions,
    comments = EXCLUDED.comments,
    last_sync_at = NOW();

-- name: GetContentStatsByID :one
SELECT 
    content_id,
    views,
    likes,
    duration_sec,
    reading_time,
    reactions,
    comments,
    last_sync_at
FROM content_stats
WHERE content_id = sqlc.arg(content_id);

-- name: GetContentStatsByIDs :many
SELECT 
    content_id,
    views,
    likes,
    duration_sec,
    reading_time,
    reactions,
    comments,
    last_sync_at
FROM content_stats
WHERE content_id = ANY(sqlc.arg(content_ids)::bigint[]);
