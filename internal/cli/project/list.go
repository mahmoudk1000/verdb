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

type listOptions struct {
	projectName string
	count       int32
}

func NewListCommand() *cobra.Command {
	opts := &listOptions{}
	var queries *database.Queries

	list := &cobra.Command{
		Use:     "list [project-name]",
		Aliases: []string{"ls"},
		Short:   "List projects",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			queries = db.Get()
			if len(args) > 0 {
				opts.projectName = args[0]
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runList(cmd.Context(), cmd, opts, queries)
		},
	}

	flags := list.Flags()
	flags.Int32VarP(&opts.count, flagNumber, "n", 0,
		"Number of projects to list (0 for all)")
	addOutputFlags(list)

	return list
}

func runList(
	ctx context.Context,
	cmd *cobra.Command,
	opts *listOptions,
	q *database.Queries,
) error {
	projects, err := fetchProjects(ctx, opts, q)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		if opts.projectName != "" {
			return fmt.Errorf(projectNotFoundErr, opts.projectName)
		}
		fmt.Println("No projects found")
		return nil
	}

	return formatAndPrint(cmd, projects)
}

func fetchProjects(
	ctx context.Context,
	opts *listOptions,
	q *database.Queries,
) ([]models.Project, error) {
	var (
		dbProjects []database.Project
		err        error
	)

	switch {
	case opts.projectName != "":
		var project database.Project
		project, err = q.GetProjectByName(ctx, opts.projectName)
		if err != nil {
			return nil, fmt.Errorf(failedToGetProjectErr, err)
		}
		dbProjects = []database.Project{project}

	case opts.count > 0:
		dbProjects, err = q.ListNProjects(ctx, opts.count)
		if err != nil {
			return nil, fmt.Errorf("failed to list projects: %w", err)
		}

	default:
		dbProjects, err = q.ListAllProjects(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list projects: %w", err)
		}
	}

	return models.ToProjects(dbProjects), nil
}
