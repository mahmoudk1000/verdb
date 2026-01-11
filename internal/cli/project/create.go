/*
Copyright © 2026 mahmoudk1000 <mahmoudk1000@gmail.com>
*/
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
		Short:   "Create a new project",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			queries = db.Get()
			opts.name = args[0]
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runCreate(cmd.Context(), opts, queries)
		},
	}

	flags := create.Flags()
	flags.SortFlags = false
	flags.StringVarP(&opts.status, "status", "s", "active", "Project status")
	flags.StringVarP(&opts.link, "link", "l", "", "Project link or URL")
	flags.StringVarP(&opts.description, "description", "d", "", "Project description")
	flags.StringArrayVarP(&opts.metadata, "metadata", "m", []string{},
		"Metadata in key=value format (can be specified multiple times)")

	return create
}

func runCreate(ctx context.Context, opts *createOptions, q *database.Queries) error {
	if err := ensureProjectNotExists(ctx, opts.name, q); err != nil {
		return err
	}

	metadataMap, err := utils.ParseMetadata(opts.metadata)
	if err != nil {
		return fmt.Errorf(failedToParseMetadataErr, err)
	}

	metadataJSON, err := utils.MetadataToJSON(metadataMap)
	if err != nil {
		return fmt.Errorf(failedToParseMetadataErr, err)
	}

	if err := createProject(ctx, opts, metadataJSON, q); err != nil {
		return err
	}

	return nil
}

func createProject(
	ctx context.Context,
	opts *createOptions,
	metadata pqtype.NullRawMessage,
	q *database.Queries,
) error {
	now := time.Now().UTC()

	params := database.CreateProjectParams{
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
	}

	if _, err := q.CreateProject(ctx, params); err != nil {
		return fmt.Errorf(failedToCreateProjectErr, err)
	}

	return nil
}
