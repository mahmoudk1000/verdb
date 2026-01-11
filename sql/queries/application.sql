-- name: CreateApplication :one
INSERT INTO applications (project_id, name, status, description, repo_url, metadata, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateApplicationStatus :exec
UPDATE applications
SET status = $2, updated_at = $3
WHERE id = $1;

-- name: UpdateApplicationMetadata :exec
UPDATE applications
SET metadata = $2, updated_at = $3
WHERE id = $1;

-- name: GetApplicationByName :one
SELECT * FROM applications
WHERE name = $1 AND project_id = $2
LIMIT 1;

-- name: GetApplicationById :one
SELECT * FROM applications
WHERE id = $1;

-- name: ListAllProjectApplications :many
SELECT * FROM applications
WHERE project_id = (
  SELECT id FROM projects WHERE projects.id = sqlc.arg(project_id)
)
ORDER BY name;

-- name: ListApplicationsByStatus :many
SELECT * FROM applications
WHERE project_id = $1 AND status = $2
ORDER BY name;

-- name: DeleteProjectApplicationByName :one
DELETE FROM applications
WHERE name = $1 AND project_id = $2
RETURNING *;

-- name: CheckApplicationExistsByName :one
SELECT EXISTS (
    SELECT 1 FROM applications
    WHERE name = $1 AND project_id = $2
) AS exists;
