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
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	flags := show.Flags()
	flags.Bool("json", false, "output in JSON format")
	flags.Bool("yaml", false, "output in YAML format")

	show.RunE = func(cmd *cobra.Command, args []string) error {
		show.SilenceUsage = true
		ctx := cmd.Context()

		var (
			fmtP string
			err  error
		)

		jsonFlag, _ := flags.GetBool("json")
		yamlFlag, _ := flags.GetBool("yaml")

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
			fmtP, err = utils.FormatJSON(p)
		case yamlFlag:
			fmtP, err = utils.FormatYAML(p)
		default:
			fmtP, err = utils.Format(p)
		}
		if err != nil {
			return err
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
) (models.Project, error) {
	exists, err := q.CheckProjectExistsByName(ctx, name)
	if err != nil {
		return models.Project{}, fmt.Errorf(checkProjectExistsErr, err)
	}
	if !exists {
		return models.Project{}, fmt.Errorf(projectNotFoundErr, name)
	}

	p, err := q.GetProjectByName(ctx, name)
	if err != nil {
		return models.Project{}, fmt.Errorf(failedToGetProjectErr, err)
	}

	return models.ToProject(p), nil
}
