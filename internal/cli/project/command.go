/*
Copyright © 2026 mahmoudk1000 <mahmoudk1000@gmail.com>
*/
package project

import (
	"github.com/spf13/cobra"
)

var (
	projectNotFoundErr = "project %q not found"
)

func NewProjectCommand() *cobra.Command {
	project := &cobra.Command{
		Use:     "project create|delete|list|show|describe|metadata",
		Aliases: []string{"proj", "projects"},
		Short:   "Manage projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	project.AddCommand(NewCreateCommand())
	project.AddCommand(NewDeleteCommand())
	project.AddCommand(NewListCommand())
	project.AddCommand(NewMetadataCommand())
	project.AddCommand(NewShowCommand())
	project.AddCommand(NewStatusCommand())

	return project
}
