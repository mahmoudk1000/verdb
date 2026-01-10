package project

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
)

func NewStatusCommand() *cobra.Command {
	var queries *database.Queries

	status := &cobra.Command{
		Use:     "status",
		Aliases: []string{"st"},
		Args:    cobra.RangeArgs(1, 2),
		Short:   "Update project status",
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	status.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		ctx := cmd.Context()

		if len(args) == 1 {
			s, err := getProjectStatus(ctx, args[0], queries)
			if err != nil {
				return err
			}
			fmt.Println(s)
		} else {
			if err := updateProjectStatus(ctx, args[0], args[1], queries); err != nil {
				return err
			}
		}
		return nil
	}

	return status
}

func updateProjectStatus(ctx context.Context, pName string, s string, q *database.Queries) error {
	if err := q.UpdateProjectStatus(ctx, database.UpdateProjectStatusParams{
		Name:      pName,
		Status:    s,
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		return err
	}

	return nil
}

func getProjectStatus(ctx context.Context, pName string, q *database.Queries) (string, error) {
	pId, err := q.GetProjectIdByName(ctx, pName)
	if err != nil {
		return "", nil
	}

	s, err := q.GetProjectStatusById(ctx, pId)
	if err != nil {
		return "", err
	}

	return s, nil
}
