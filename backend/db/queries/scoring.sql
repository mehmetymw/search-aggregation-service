-- name: GetScoringRules :many
SELECT key, value FROM scoring_rules;

-- name: GetScoringRule :one
SELECT value FROM scoring_rules WHERE key = $1;

-- name: UpsertScoringRule :exec
INSERT INTO scoring_rules (key, value, description, updated_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (key) DO UPDATE
SET value = EXCLUDED.value,
    description = EXCLUDED.description,
    updated_at = NOW();
