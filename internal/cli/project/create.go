package project

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/sqlc-dev/pqtype"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
	"github.com/mahmoudk1000/relen/internal/utils"
)

func NewCreateCommand() *cobra.Command {
	var (
		status      string
		link        string
		description string
		metadata    []string
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
	flags.StringVarP(&status, "status", "s", "active", "Project status")
	flags.StringVarP(&link, "link", "l", "", "Project link")
	flags.StringVarP(&description, "description", "d", "", "Project description")
	flags.StringArrayVarP(&metadata, "metadata", "m", []string{}, "Metadata key=value pairs")

	create.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		ctx := cmd.Context()

		metadataMap, err := utils.ParseMetadata(metadata)
		if err != nil {
			return fmt.Errorf("failed to parse metadata: %w", err)
		}

		metadataJSON, err := utils.MetadataToJSON(metadataMap)
		if err != nil {
			return err
		}

		return createJSONProject(ctx, args[0], status, link, description, metadataJSON, queries)
	}

	return create
}

func createJSONProject(
	ctx context.Context,
	name, status, link, desc string,
	md pqtype.NullRawMessage,
	q *database.Queries,
) error {
	exists, err := q.CheckProjectExistsByName(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check if project exists: %w", err)
	}
	if exists {
		return fmt.Errorf("project with name '%s' already exists", name)
	}

	now := time.Now().UTC()
	_, err = q.CreateProject(ctx, database.CreateProjectParams{
		Name:   name,
		Status: status,
		Link: sql.NullString{
			String: link,
			Valid:  link != "",
		},
		Description: sql.NullString{
			String: desc,
			Valid:  desc != "",
		},
		Metadata:  md,
		CreatedAt: now,
		UpdatedAt: now,
	})

	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}
