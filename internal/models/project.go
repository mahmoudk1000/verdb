package models

import (
	"encoding/json"
	"time"

	"github.com/mahmoudk1000/relen/internal/database"
)

type Project struct {
	Name        string         `json:"name"`
	Status      string         `json:"status,omitempty"`
	Link        string         `json:"link,omitempty"`
	Description string         `json:"description,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	Created_At  string         `json:"created_at"`
	Updated_At  string         `json:"updated_at,omitempty"`
}

type FProject struct {
	Project
	Application []Application `json:"applications,omitempty"`
}

func ToProject(p database.Project) Project {
	var metadata map[string]any
	if p.Metadata.Valid && len(p.Metadata.RawMessage) > 0 {
		_ = json.Unmarshal(p.Metadata.RawMessage, &metadata)
	}

	return Project{
		Name:        p.Name,
		Status:      p.Status,
		Link:        p.Link.String,
		Description: p.Description.String,
		Metadata:    metadata,
		Created_At:  p.CreatedAt.Format(time.RFC1123),
		Updated_At:  p.UpdatedAt.Format(time.RFC1123),
	}
}

func ToProjects(ps []database.Project) []Project {
	result := make([]Project, 0, len(ps))
	for _, p := range ps {
		result = append(result, ToProject(p))
	}

	return result
}
