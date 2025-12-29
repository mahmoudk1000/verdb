package project

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/verdb/internal/models"
	"github.com/mahmoudk1000/verdb/internal/utils"
)

func showCommand() *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			err := showJSONProject(configBuilder, args[0])
			if err != nil {
				return fmt.Errorf("failed to show project: %w", err)
			}
			return nil
		},
	}

	return show
}

func showJSONProject(configBuilder *utils.ConfigBuilder[models.Projects], name string) error {
	return nil
}
