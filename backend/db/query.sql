-- name: InsertParticipant :one
INSERT INTO participants (
    name,
    email,
    wa_number
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetParticipantByID :one
SELECT *
FROM participants
WHERE id = $1;

-- name: GetParticipantByExternalID :one
SELECT *
FROM participants
WHERE external_id = $1;

-- name: ListParticipants :many
SELECT *
FROM participants
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateParticipantAccessed :one
UPDATE participants
SET accessed = TRUE
WHERE email = $1
   OR wa_number = $2
   AND id = $3
RETURNING *;

-- name: UpdateParticipantAccessedByExternalID :one
UPDATE participants
SET accessed = TRUE
WHERE external_id = $1
RETURNING *;

-- name: DeleteParticipant :exec
DELETE FROM participants
WHERE id = $1;

-- name: GetUnsentInvites :many
SELECT *
FROM participants
WHERE sent = FALSE
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: MarkInvitesAsSent :exec
UPDATE participants
SET sent = TRUE
WHERE id = ANY($1::int[]);
