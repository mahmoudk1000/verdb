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
)

type deleteOptions struct {
	projectName string
	confirmed   bool
}

func NewDeleteCommand() *cobra.Command {
	opts := &deleteOptions{}
	var queries *database.Queries

	delete := &cobra.Command{
		Use:     "delete <project-name>",
		Aliases: []string{"del", "rm", "remove"},
		Short:   "Delete a project",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			queries = db.Get()
			opts.projectName = args[0]

			if !opts.confirmed {
				return fmt.Errorf(
					"deletion requires explicit confirmation: use --yes-i-am-sure flag",
				)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runDelete(cmd.Context(), opts, queries)
		},
	}

	flags := delete.Flags()
	flags.BoolVar(&opts.confirmed, "yes-i-am-sure", false,
		"Confirm project deletion without prompting")

	return delete
}

func runDelete(ctx context.Context, opts *deleteOptions, q *database.Queries) error {
	if err := ensureProjectExists(ctx, opts.projectName, q); err != nil {
		return err
	}

	return deleteProject(ctx, opts.projectName, q)
}

func deleteProject(ctx context.Context, name string, q *database.Queries) error {
	if err := q.DeleteProjectByName(ctx, name); err != nil {
		return fmt.Errorf(failedToDeleteProjectErr, err)
	}

	return nil
}
