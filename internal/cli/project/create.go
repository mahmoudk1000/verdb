package project

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
)

func NewCreateCommand() *cobra.Command {
	var (
		link        string
		description string
		queries     *database.Queries
	)

	create := &cobra.Command{
		Use:     "create <name>",
		Aliases: []string{"c", "new"},
		Short:   "add a new application to the project",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	flags := create.Flags()
	flags.SortFlags = false
	flags.StringVarP(&link, "link", "l", "", "link to the project")
	flags.StringVarP(&description, "description", "d", "", "description of the application")

	create.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		ctx := cmd.Context()

		return createJSONProject(ctx, args[0], link, description, queries)
	}

	return create
}

func createJSONProject(
	ctx context.Context,
	name, link, desc string,
	q *database.Queries,
) error {

	exists, err := q.CheckProjectExistsByName(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check if project exists: %w", err)
	}
	if exists {
		return fmt.Errorf("project with name '%s' already exists", name)
	}

	_, err = q.CreateProject(ctx, database.CreateProjectParams{
		Name: name,
		Link: sql.NullString{
			String: link,
			Valid:  link != "",
		},
		Description: sql.NullString{
			String: desc,
			Valid:  desc != "",
		},
		CreatedAt: time.Now().UTC(),
	})

	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}
