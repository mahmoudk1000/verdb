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

func NewDescribeCommand() *cobra.Command {
	var queries *database.Queries

	describe := &cobra.Command{
		Use:   "describe <project_name> <application_name>",
		Short: "Describe an application",
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	flags := describe.Flags()
	flags.Bool("json", false, "Output in JSON format")

	describe.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		ctx := cmd.Context()
		jsonFlag, _ := cmd.Flags().GetBool("json")

		var fmtA string
		a, err := describeApplication(ctx, args[0], args[1], queries)
		if err != nil {
			return err
		}

		switch {
		case jsonFlag:
			fmtA, err = utils.FormatJSON(a)
			if err != nil {
				return err
			}
		default:
			fmtA, err = utils.Format(a)
			if err != nil {
				return err
			}
		}

		fmt.Println(fmtA)

		return nil
	}

	return describe
}

func describeApplication(
	ctx context.Context,
	pName, aName string,
	q *database.Queries,
) (models.Application, error) {
	pId, err := q.GetProjectIdByName(ctx, pName)
	if err != nil {
		return models.Application{}, fmt.Errorf("failed to get project id: %w", err)
	}

	a, err := q.GetApplicationByName(ctx, database.GetApplicationByNameParams{
		Name: aName,
		ID:   pId,
	})
	if err != nil {
		return models.Application{}, fmt.Errorf("failed to find application: %w", err)
	}

	return models.ToApplication(a), nil
}
