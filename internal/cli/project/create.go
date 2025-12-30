package project

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/models"
	"github.com/mahmoudk1000/relen/internal/utils"
)

func NewCreateCommand() *cobra.Command {
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

	flags := create.Flags()
	flags.SortFlags = false
	flags.StringVarP(&link, "link", "l", "", "link to the project")
	flags.StringVarP(&description, "description", "d", "", "description of the application")

	create.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		err := createJSONProject(args[0], link, description, configBuilder)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		fmt.Printf("Project '%s' was created\n", args[0])
		return nil
	}

	return create
}

func createJSONProject(
	name, link, desc string,
	cb *utils.ConfigBuilder[models.Projects],
) error {
	projects := cb.Model()

	for _, p := range projects.Project {
		if p.Name == name {
			return fmt.Errorf("project with name '%s' already exists", name)
		}
	}

	p := models.Project{
		Name:        name,
		Link:        link,
		Description: desc,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	projects.Project = append(projects.Project, p)
	cb.SetModel(projects)

	if err := cb.Save(); err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	return nil
}
