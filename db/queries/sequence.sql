-- name: GetSequences :many
select 
    s.*, 
    json_agg(row_to_json(t))::jsonb steps 
from sequences s
left join steps t on t.sequence_id = s.id
group by
	s.id,
	s.external_id,
	s.sequence_name,
	s.open_tracking_enabled,
	s.click_tracking_enabled,
	s.created,
	s.updated
order by s.id
limit $1
offset $2;

-- name: GetSequenceById :one
select 
    s.*, 
	json_agg(row_to_json(t))::jsonb steps 
from sequences s
left join steps t on t.sequence_id = s.id
where s.external_id = $1
group by
	s.id,
	s.external_id,
	s.sequence_name,
	s.open_tracking_enabled,
	s.click_tracking_enabled,
	s.created,
	s.updated;

-- name: CreateSequence :one
INSERT INTO sequences (sequence_name, open_tracking_enabled, click_tracking_enabled) 
VALUES ($1, $2, $3) 
RETURNING *;

-- name: UpdateSequence :one
UPDATE sequences 
SET open_tracking_enabled = $2, click_tracking_enabled = $3 
WHERE id = $1 
RETURNING *;

-- name: DeleteSequence :exec
DELETE FROM sequences 
WHERE external_id = $1;