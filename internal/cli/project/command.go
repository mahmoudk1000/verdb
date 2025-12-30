/*
Copyright © 2026 mahmoudk1000 <mahmoudk1000@gmail.com>
*/
package project

import (
	"github.com/spf13/cobra"
)

const projectsFileName string = "projects.json"

func NewProjectCommand() *cobra.Command {
	project := &cobra.Command{
		Use:     "project",
		Aliases: []string{"p", "proj"},
		Short:   "manage project",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Help()
			if err != nil {
				return err
			}
			return nil
		},
	}

	project.AddCommand(NewCreateCommand())
	project.AddCommand(NewDeleteCommand())
	project.AddCommand(NewShowCommand())

	return project
}
