package project

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/models"
	"github.com/mahmoudk1000/relen/internal/utils"
)

func NewDeleteCommand() *cobra.Command {
	var (
		configBuilder *utils.ConfigBuilder[models.Projects]
	)

	delete := &cobra.Command{
		Use:     "delete <project-name>",
		Aliases: []string{"del", "rm"},
		Short:   "Delete a project",
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

	flags := delete.Flags()
	flags.Bool("yes-i-am-sure", false, "Confirm project deletion without prompting")

	delete.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		if yes, _ := cmd.Flags().GetBool("yes-i-am-sure"); !yes {
			fmt.Println("Please confirm project deletion with --yes-i-am-sure flag")
			return nil
		}

		err := deleteProject(args[0], configBuilder)
		if err != nil {
			return fmt.Errorf("failed to delete project: %w", err)
		}
		fmt.Printf("Project '%s' deleted successfully\n", args[0])

		return nil
	}

	return delete
}

func deleteProject(name string, cb *utils.ConfigBuilder[models.Projects]) error {
	projects := cb.Model()

	for i, p := range projects.Project {
		if p.Name == name {
			projects.Project = append(projects.Project[:i], projects.Project[i+1:]...)

			cb.SetModel(projects)
			if err := cb.Save(); err != nil {
				return fmt.Errorf("failed to save updated projects: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("project '%s' not found", name)
}
