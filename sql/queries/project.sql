-- name: CreateProject :one
INSERT INTO projects (name, link,  description, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CheckProjectExistsByName :one
SELECT EXISTS (
    SELECT 1 FROM projects WHERE name = $1
) AS exists;

-- name: GetProjectByName :one
SELECT * FROM projects
WHERE name = $1
LIMIT 1;

-- name: DeleteProjectByName :exec
DELETE FROM projects
WHERE name = $1;

-- name: ListAllProjects :many
SELECT name FROM projects;

-- name: CreateProjectVersion :one
INSERT INTO project_versions (id, project_id, version, description, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetLatestProjectVersionByProjectName :one
SELECT * FROM project_versions
  WHERE project_id = (
    SELECT id FROM projects WHERE name = $1
  )
  ORDER BY created_at DESC
LIMIT 1;

-- name: GetNVersionsByProjectName :many
SELECT * FROM project_versions
WHERE project_id = (
  SELECT id FROM projects WHERE name = $1
)
ORDER BY created_at DESC
LIMIT $2;
