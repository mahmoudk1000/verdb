package application

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
	"github.com/mahmoudk1000/relen/internal/utils"
)

func NewAddCommand() *cobra.Command {
	var (
		link        string
		description string
		queries     *database.Queries
	)

	add := &cobra.Command{
		Use:     "add <project_name> <application_name>",
		Aliases: []string{"a", "new"},
		Short:   "add a new application to the project",
		Args:    cobra.RangeArgs(1, 2),
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	flags := add.Flags()
	flags.StringVarP(&link, "link", "l", "", "application's link")
	flags.StringVarP(&description, "description", "d", "", "application's description")

	add.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		pName, aName, err := utils.ParseProjectSlashApplication(args)
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		linkFlag, _ := flags.GetString("link")
		descFlag, _ := flags.GetString("description")

		return addApplication(ctx,
			pName,
			aName,
			linkFlag,
			descFlag,
			queries,
		)
	}

	return add
}

func addApplication(
	ctx context.Context,
	prjName, appName, link, desc string,
	q *database.Queries,
) error {
	pID, err := q.GetProjectIdByName(ctx, prjName)
	if err != nil {
		return fmt.Errorf(projectNotFoundErr, prjName, err)
	}

	exists, err := q.CheckApplicationExistsByName(ctx, database.CheckApplicationExistsByNameParams{
		Name: prjName,
		ProjectID: pID,
	})
	if err != nil {
		return fmt.Errorf(checkApplicationExistsErr, err)
	}
	if exists {
		return fmt.Errorf(applicationExistsErr, appName)
	}

	if _, err := q.CreateApplication(ctx, database.CreateApplicationParams{
		Name:      appName,
		ProjectID: pID,
		RepoUrl: sql.NullString{
			String: link,
			Valid:  link != "",
		},
		Description: sql.NullString{
			String: desc,
			Valid:  desc != "",
		},
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		return fmt.Errorf(failedToCreateApplicationErr, err)
	}

	return nil
}
