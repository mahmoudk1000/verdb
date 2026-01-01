package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
	"github.com/mahmoudk1000/relen/internal/models"
	"github.com/mahmoudk1000/relen/internal/utils"
)

func NewShowCommand() *cobra.Command {
	var queries *database.Queries

	show := &cobra.Command{
		Use:     "show <name>",
		Aliases: []string{"s"},
		Short:   "show details of a project",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			queries = db.Get()
			return nil
		},
	}

	flags := show.Flags()
	flags.Bool("json", false, "output in JSON format")

	show.RunE = func(cmd *cobra.Command, args []string) error {
		var fmtP string
		show.SilenceUsage = true
		ctx := cmd.Context()

		jsonFlag, _ := flags.GetBool("json")

		p, err := showProject(
			ctx,
			args[0],
			queries,
		)
		if err != nil {
			return err
		}

		switch {
		case jsonFlag:
			fmtP, err = utils.FormatJSON(models.DatabaseProjectToProject(p))
			if err != nil {
				return err
			}
		default:
			fmtP, err = utils.Format(models.DatabaseProjectToProject(p))
			if err != nil {
				return err
			}
		}

		fmt.Println(fmtP)

		return nil
	}

	return show
}

func showProject(
	ctx context.Context,
	name string,
	q *database.Queries,
) (database.Project, error) {

	exists, err := q.CheckProjectExistsByName(ctx, name)
	if err != nil {
		return database.Project{}, fmt.Errorf("failed to check if project exists: %w", err)
	}
	if !exists {
		return database.Project{}, fmt.Errorf("project '%s' does not exist", name)
	}

	p, err := q.GetProjectByName(ctx, name)
	if err != nil {
		return database.Project{}, fmt.Errorf("failed to get project: %w", err)
	}

	return p, nil
}
