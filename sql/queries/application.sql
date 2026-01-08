-- name: CreateApplication :one
INSERT INTO applications (project_id, name, description, repo_url, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CheckApplicationExistsByName :one
SELECT EXISTS (
    SELECT 1 FROM applications WHERE applications.name = $1 AND applications.project_id = (
      SELECT id FROM projects WHERE projects.id = $2
    )
) AS exists;

-- name: GetApplicationByName :one
SELECT * FROM applications
WHERE applications.name = $1 AND project_id = (
  SELECT id FROM projects WHERE projects.id = $2
) LIMIT 1;

-- name: DeleteProjectApplicationByName :one
DELETE FROM applications
WHERE applications.name = $1 AND project_id = (
  SELECT id FROM projects WHERE projects.id = $2
) RETURNING *;

-- name: ListAllProjectApplications :many
SELECT * FROM applications
WHERE project_id = (
  SELECT id FROM projects WHERE projects.name = $1
);

-- name: GetLatestApplicationVersionByApplicationName :one
SELECT av.* FROM application_versions av
JOIN applications a ON av.application_id = a.id
WHERE a.name = $1
ORDER BY av.created_at DESC
LIMIT 1;
