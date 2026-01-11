/*
Copyright © 2026 mahmoudk1000 <mahmoudk1000@gmail.com>
*/
package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
	"github.com/mahmoudk1000/relen/internal/models"
)

func NewShowCommand() *cobra.Command {
	var queries *database.Queries

	show := &cobra.Command{
		Use:     "show <name>",
		Aliases: []string{"s", "describe", "get"},
		Short:   "Show details of a project",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			queries = db.Get()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runShow(cmd.Context(), cmd, args[0], queries)
		},
	}

	addOutputFlags(show)

	return show
}

func runShow(
	ctx context.Context,
	cmd *cobra.Command,
	projectName string,
	q *database.Queries,
) error {
	project, err := getProjectByName(ctx, projectName, q)
	if err != nil {
		return err
	}

	return formatAndPrint(cmd, project)
}

func getProjectByName(
	ctx context.Context,
	name string,
	q *database.Queries,
) (models.Project, error) {
	if err := ensureProjectExists(ctx, name, q); err != nil {
		return models.Project{}, err
	}

	dbProject, err := q.GetProjectByName(ctx, name)
	if err != nil {
		return models.Project{}, fmt.Errorf(failedToGetProjectErr, err)
	}

	return models.ToProject(dbProject), nil
}
