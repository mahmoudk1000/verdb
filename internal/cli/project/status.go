/*
Copyright © 2026 mahmoudk1000 <mahmoudk1000@gmail.com>
*/
package project

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
)

type statusOptions struct {
	projectName string
	newStatus   string
	isUpdate    bool
}

func NewStatusCommand() *cobra.Command {
	opts := &statusOptions{}
	var queries *database.Queries

	status := &cobra.Command{
		Use:     "status <project-name> [new-status]",
		Aliases: []string{"st"},
		Short:   "Get or update project status",
		Args:    cobra.RangeArgs(1, 2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			queries = db.Get()
			opts.projectName = args[0]

			if len(args) == 2 {
				opts.newStatus = args[1]
				opts.isUpdate = true
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if opts.isUpdate {
				return runStatusUpdate(cmd.Context(), cmd, opts, queries)
			}
			return runStatusGet(cmd.Context(), cmd, opts, queries)
		},
	}

	addOutputFlags(status)
	status.Flags().BoolP(flagQuiet, "q", false, "Quiet mode, no output")

	return status
}

func runStatusGet(
	ctx context.Context,
	cmd *cobra.Command,
	opts *statusOptions,
	q *database.Queries,
) error {
	status, err := getProjectStatus(ctx, opts.projectName, q)
	if err != nil {
		return err
	}

	if isQuietMode(cmd) {
		return nil
	}

	format, err := getOutputFormat(cmd)
	if err != nil {
		return err
	}

	if format == formatJSON {
		output, err := formatOutput(struct {
			ProjectName string `json:"project_name"`
			Status      string `json:"status"`
		}{
			ProjectName: opts.projectName,
			Status:      status,
		}, formatJSON)
		if err != nil {
			return err
		}
		fmt.Println(output)
	} else if format == formatYAML {
		output, err := formatOutput(struct {
			ProjectName string `yaml:"project_name"`
			Status      string `yaml:"status"`
		}{
			ProjectName: opts.projectName,
			Status:      status,
		}, formatYAML)
		if err != nil {
			return err
		}
		fmt.Println(output)
	} else {
		fmt.Println(status)
	}

	return nil
}

func runStatusUpdate(
	ctx context.Context,
	cmd *cobra.Command,
	opts *statusOptions,
	q *database.Queries,
) error {
	if err := updateProjectStatus(ctx, opts.projectName, opts.newStatus, q); err != nil {
		return err
	}

	if !isQuietMode(cmd) {
		fmt.Printf("Successfully updated status of project %q to %q\n",
			opts.projectName, opts.newStatus)
	}

	return nil
}

func getProjectStatus(
	ctx context.Context,
	projectName string,
	q *database.Queries,
) (string, error) {
	projectID, err := q.GetProjectIdByName(ctx, projectName)
	if err != nil {
		return "", fmt.Errorf(projectNotFoundErr, projectName)
	}

	status, err := q.GetProjectStatusById(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf(failedToGetProjectErr, err)
	}

	return status, nil
}

func updateProjectStatus(
	ctx context.Context,
	projectName, newStatus string,
	q *database.Queries,
) error {
	if err := ensureProjectExists(ctx, projectName, q); err != nil {
		return err
	}

	params := database.UpdateProjectStatusByNameParams{
		Name:      projectName,
		Status:    newStatus,
		UpdatedAt: time.Now().UTC(),
	}

	if err := q.UpdateProjectStatusByName(ctx, params); err != nil {
		return fmt.Errorf(failedToUpdateProjectErr, err)
	}

	return nil
}
