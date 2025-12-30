package project

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/models"
	"github.com/mahmoudk1000/relen/internal/utils"
)

func NewShowCommand() *cobra.Command {
	var configBuilder *utils.ConfigBuilder[models.Projects]

	show := &cobra.Command{
		Use:     "show <name>",
		Aliases: []string{"s"},
		Short:   "show details of a project",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			configBuilder = utils.NewConfigBuilder(projectsFileName, models.Projects{})
			err := configBuilder.BuildConfigDir()
			if err != nil {
				return err
			}
			return nil
		},
	}

	show.SilenceUsage = true

	flags := show.Flags()
	flags.Bool("json", false, "output in JSON format")

	show.RunE = func(cmd *cobra.Command, args []string) error {
		var (
			p   string
			err error
		)

		jsonFlag, _ := show.Flags().GetBool("json")

		switch {
		case jsonFlag:
			p, err = showProject(
				args[0],
				configBuilder,
				func(data any) (string, error) {
					return utils.FormatJSON(data)
				},
			)
		default:
			p, err = showProject(
				args[0],
				configBuilder,
				func(data any) (string, error) {
					return utils.Format(data)
				},
			)
		}

		if err != nil {
			return err
		}
		fmt.Println(p)

		return nil
	}

	return show
}

func showProject(
	name string,
	cb *utils.ConfigBuilder[models.Projects],
	output func(any) (string, error),
) (string, error) {
	projects := cb.Model()

	for _, p := range projects.Project {
		if p.Name == name {
			fmpP, err := output(p)
			if err != nil {
				return "", err
			}
			return fmpP, nil
		}
	}

	return "", fmt.Errorf("project '%s' not found", name)
}
