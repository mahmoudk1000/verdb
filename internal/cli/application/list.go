package application

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
	"github.com/mahmoudk1000/relen/internal/models"
	"github.com/mahmoudk1000/relen/internal/utils"
)

func NewListCommand() *cobra.Command {
	var queries *database.Queries

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all applications of a project",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	flags := list.Flags()
	flags.Bool("json", false, "output in JSON format")
	flags.Bool("yaml", false, "output in YAML format")

	list.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		var fmtA string
		ctx := cmd.Context()
		jsonFlag, _ := flags.GetBool("json")
		yamlFlag, _ := flags.GetBool("yaml")

		applications, err := listApplications(ctx, args[0], queries)
		if err != nil {
			return err
		}

		switch {
		case jsonFlag:
			fmtA, err = utils.FormatJSON(applications)
			if err != nil {
				return err
			}
		case yamlFlag:
			fmtA, err = utils.FormatYAML(applications)
			if err != nil {
				return err
			}
		default:
			fmtA, err = utils.Format(applications)
			if err != nil {
				return err
			}
		}

		fmt.Println(fmtA)

		return nil
	}

	return list
}

func listApplications(
	ctx context.Context,
	pName string,
	q *database.Queries,
) ([]models.Application, error) {
	pId, err := q.GetProjectIdByName(ctx, pName)
	if err != nil {
		return nil, fmt.Errorf(projectNotFoundErr, pName, err)
	}

	ps, err := q.ListAllProjectApplications(ctx, pId)
	if err != nil {
		return nil, fmt.Errorf(failedToListApplicationsErr, pName, err)
	}

	return models.ToApplications(ps), nil
}
