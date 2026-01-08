package models

import (
	"time"

	"github.com/mahmoudk1000/relen/internal/database"
)

type Project struct {
	Name        string `json:"name"`
	Link        string `json:"link,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type FProject struct {
	Project
	Application []Application `json:"applications,omitempty"`
}

func ToProject(p database.Project) Project {
	return Project{
		Name:        p.Name,
		Link:        p.Link.String,
		Description: p.Description.String,
		CreatedAt:   p.CreatedAt.Format(time.RFC1123),
	}
}

func ToProjects(ps []database.Project) []Project {
	result := make([]Project, 0, len(ps))
	for _, p := range ps {
		result = append(result, ToProject(p))
	}

	return result
}
