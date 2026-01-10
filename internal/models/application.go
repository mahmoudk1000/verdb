package models

import (
	"encoding/json"

	"github.com/mahmoudk1000/relen/internal/database"
)

type Application struct {
	Name        string         `json:"name"`
	Status      string         `json:"status,omitempty"`
	Repo_Url    string         `json:"repo,omitempty"`
	Description string         `json:"description,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	Created_At  string         `json:"created_at"`
	Updated_At  string         `json:"updated_at,omitempty"`
}

func ToApplication(a database.Application) Application {
	var metadata map[string]any
	if a.Metadata.Valid && len(a.Metadata.RawMessage) > 0 {
		_ = json.Unmarshal(a.Metadata.RawMessage, &metadata)
	}

	return Application{
		Name:        a.Name,
		Repo_Url:    a.RepoUrl.String,
		Description: a.Description.String,
		Created_At:  a.CreatedAt.Format("2006-01-02T15:04:05 -07:00:00"),
	}
}

func ToApplications(apps []database.Application) []Application {
	results := make([]Application, 0, len(apps))
	for _, a := range apps {
		results = append(results, ToApplication(a))
	}

	return results
}
