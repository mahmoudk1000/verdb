/*
Copyright © 2026 (mahmoudk1000) <mahmoudk1000@gmail.com>
*/
package application

import (
	"github.com/spf13/cobra"
)

var (
	// Error message templates
	projectNotFoundErr           = "project %q does not exist: %w"
	applicationNotFoundErr       = "application %q not found"
	applicationExistsErr         = "application with name %q already exists"
	checkApplicationExistsErr    = "failed checking application existence: %w"
	failedToCreateApplicationErr = "failed to create application: %w"
	failedToDeleteApplicationErr = "failed to delete application: %w"
	failedToListApplicationsErr  = "failed to list applications for project %q: %w"
)

func NewApplicationCommand() *cobra.Command {
	application := &cobra.Command{
		Use:     "application add|remove|list|show",
		Aliases: []string{"app", "applications"},
		Short:   "Manage applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	application.AddCommand(NewAddCommand())
	application.AddCommand(NewDeleteCommand())
	application.AddCommand(NewListCommand())

	return application
}
