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

type createOptions struct {
	name        string
	status      string
	link        string
	description string
	metadata    []string
}

func NewCreateCommand() *cobra.Command {
	opts := &createOptions{}
	var queries *database.Queries

	create := &cobra.Command{
		Use:     "create <name>",
		Aliases: []string{"c", "new"},
		Short:   "add a new application to the project",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
			opts.name = args[0]
		},
	}

	flags := create.Flags()
	flags.SortFlags = false
	flags.StringVarP(&opts.status, "status", "s", "active", "Project status")
	flags.StringVarP(&opts.link, "link", "l", "", "Project link")
	flags.StringVarP(&opts.description, "description", "d", "", "Project description")
	flags.StringArrayVarP(&opts.metadata, "metadata", "m", []string{}, "Metadata key=value pairs")

	create.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		metadataMap, err := utils.ParseMetadata(opts.metadata)
		if err != nil {
			return fmt.Errorf(failedToParseMetadataErr, err)
		}

		metadata, err := utils.MetadataToJSON(metadataMap)
		if err != nil {
			return err
		}

		return createProject(cmd.Context(), opts, metadata, queries)
	}

	return create
}

func createProject(
	ctx context.Context,
	opts *createOptions,
	metadata pqtype.NullRawMessage,
	q *database.Queries,
) error {
	exists, err := q.CheckProjectExistsByName(ctx, opts.name)
	if err != nil {
		return fmt.Errorf(checkProjectExistsErr, err)
	}
	if exists {
		return fmt.Errorf(projectExistsErr, opts.name)
	}

	now := time.Now().UTC()
	if _, err = q.CreateProject(ctx, database.CreateProjectParams{
		Name:   opts.name,
		Status: opts.status,
		Link: sql.NullString{
			String: opts.link,
			Valid:  opts.link != "",
		},
		Description: sql.NullString{
			String: opts.description,
			Valid:  opts.description != "",
		},
		Metadata:  metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}); err != nil {
		return fmt.Errorf(failedToCreateProjectErr, err)
	}

	return nil
}
