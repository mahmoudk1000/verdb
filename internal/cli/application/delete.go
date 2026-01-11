package application

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
	"github.com/mahmoudk1000/relen/internal/utils"
)

func NewDeleteCommand() *cobra.Command {
	var queries *database.Queries

	delete := &cobra.Command{
		Use:     "delete <project_name> <application_name>",
		Args:    cobra.RangeArgs(1, 2),
		Aliases: []string{"remove", "rm"},
		Short:   "Delete an application",
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	delete.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		ctx := cmd.Context()
		pName, aName, err := utils.ParseProjectSlashApplication(args)
		if err != nil {
			return err
		}

		return deleteApplication(ctx, pName, aName, queries)
	}
	return delete
}

func deleteApplication(ctx context.Context, pName, aName string, q *database.Queries) error {
	pId, err := q.GetProjectIdByName(ctx, pName)
	if err != nil {
		return err
	}

	if _, err := q.CheckApplicationExistsByName(ctx, database.CheckApplicationExistsByNameParams{
		Name: aName,
		ProjectID: pId,
	}); err != nil {
		return err
	}

	if _, err = q.DeleteProjectApplicationByName(ctx, database.DeleteProjectApplicationByNameParams{
		Name: aName,
		ProjectID: pId,
	}); err != nil {
		return fmt.Errorf(failedToDeleteApplicationErr, err)
	}

	return err
}
