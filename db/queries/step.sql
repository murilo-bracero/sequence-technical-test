-- name: CreateSteps :copyfrom
INSERT INTO steps (external_id, mail_subject, mail_content, sequence_id) 
VALUES ($1, $2, $3, $4);

-- name: CreateStep :one
INSERT INTO steps (mail_subject, mail_content, sequence_id) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetStepById :one
SELECT steps.* FROM steps
JOIN sequences ON steps.sequence_id = sequences.id AND sequences.external_id = $2
WHERE steps.external_id = $1;

-- name: UpdateStep :one
UPDATE steps 
SET mail_subject = $2, mail_content = $3 
WHERE external_id = $1 
RETURNING *;

-- name: DeleteStep :exec
DELETE FROM steps 
WHERE external_id = $1;
