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

func NewListCommand() *cobra.Command {
	var queries *database.Queries

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.RangeArgs(0, 1),
		Short:   "List all projects",
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	flags := list.Flags()
	flags.Bool("json", false, "Output in JSON format")
	flags.Bool("yaml", false, "Output in YAML format")
	flags.Int32P("number", "n", 0, "Number of projects to list (0 for all)")

	list.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		ctx := cmd.Context()

		jsonFlag, _ := flags.GetBool("json")
		yamlFlag, _ := flags.GetBool("yaml")
		count, _ := flags.GetInt32("number")

		var pName string
		if len(args) > 0 {
			pName = args[0]
		}

		ps, err := listProjects(ctx, pName, count, queries)
		if err != nil {
			return err
		}

		var fmtP string
		switch {
		case jsonFlag:
			fmtP, err = utils.FormatJSON(ps)
		case yamlFlag:
			fmtP, err = utils.FormatYAML(ps)
		default:
			fmtP, err = utils.Format(ps)
		}
		if err != nil {
			return err
		}

		fmt.Println(fmtP)
		return nil
	}

	return list
}

func listProjects(
	ctx context.Context,
	pName string,
	c int32,
	q *database.Queries,
) ([]models.Project, error) {
	var (
		ps  []database.Project
		err error
	)

	switch {
	case pName != "":
		var p database.Project
		p, err = q.GetProjectByName(ctx, pName)
		if err != nil {
			return nil, fmt.Errorf(failedToGetProjectErr, err)
		}
		ps = []database.Project{p}
	case c > 0:
		ps, err = q.ListNProjects(ctx, c)
	default:
		ps, err = q.ListAllProjects(ctx)
	}

	if err != nil {
		return nil, err
	}

	return models.ToProjects(ps), nil
}
