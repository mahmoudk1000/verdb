package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
)

func NewDeleteCommand() *cobra.Command {
	var queries *database.Queries

	delete := &cobra.Command{
		Use:     "delete <project-name>",
		Aliases: []string{"del", "rm"},
		Short:   "Delete a project",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	flags := delete.Flags()
	flags.Bool("yes-i-am-sure", false, "Confirm project deletion without prompting")

	delete.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		ctx := cmd.Context()

		if yes, _ := cmd.Flags().GetBool("yes-i-am-sure"); !yes {
			fmt.Println("Please confirm project deletion with --yes-i-am-sure flag")
			return nil
		}

		return deleteProject(ctx, args[0], queries)
	}

	return delete
}

func deleteProject(ctx context.Context, name string, q *database.Queries) error {
	exists, err := q.CheckProjectExistsByName(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check if project exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("project '%s' does not exist", name)
	}

	return q.DeleteProjectByName(ctx, name)
}
