package project

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/verdb/internal/models"
	"github.com/mahmoudk1000/verdb/internal/utils"
)

func createCommand() *cobra.Command {
	var (
		link          string
		description   string
		configBuilder *utils.ConfigBuilder[models.Projects]
	)

	create := &cobra.Command{
		Use:     "create <name>",
		Aliases: []string{"c", "new"},
		Short:   "add a new application to the project",
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

	create.SilenceUsage = true

	flags := create.Flags()
	flags.StringVarP(&link, "link", "l", "", "link to the project")
	flags.StringVarP(&description, "description", "d", "", "description of the application")

	create.RunE = func(cmd *cobra.Command, args []string) error {
		err := createJSONProject(configBuilder, args[0], link, description)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		fmt.Printf("Project '%s' was created\n", args[0])
		return nil
	}

	return create
}

func createJSONProject(
	configBuilder *utils.ConfigBuilder[models.Projects],
	name, link, desc string,
) error {
	projects := configBuilder.Model()

	if _, exists := projects.Project[name]; exists {
		return fmt.Errorf("project '%s' already exists", name)
	}

	projects.Project[name] = models.Project{
		Link:        link,
		Description: desc,
		CreatedAt:   time.Now().UTC(),
	}

	if err := configBuilder.Save(); err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	return nil
}
