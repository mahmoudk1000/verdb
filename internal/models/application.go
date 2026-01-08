package models

import "github.com/mahmoudk1000/relen/internal/database"

type Application struct {
	Name        string `json:"name"`
	Repo_Url    string `json:"repo,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}

func ToApplication(a database.Application) Application {
	return Application{
		Name:        a.Name,
		Repo_Url:    a.RepoUrl.String,
		Description: a.Description.String,
		CreatedAt:   a.CreatedAt.Format("2006-01-02T15:04:05 -07:00:00"),
	}
}

func ToApplications(apps []database.Application) []Application {
	results := make([]Application, 0, len(apps))
	for _, a := range apps {
		results = append(results, ToApplication(a))
	}

	return results
}
