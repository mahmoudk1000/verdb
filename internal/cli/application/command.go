/*
Copyright © 2026 (mahmoudk1000) <mahmoudk1000@gmail.com>
*/
package application

import (
	"github.com/spf13/cobra"
)

func NewApplicationCommand() *cobra.Command {
	application := &cobra.Command{
		Use:   "application add|remove|list|show",
		Short: "Manage applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	application.AddCommand(NewAddCommand())
	application.AddCommand(NewDeleteCommand())
	application.AddCommand(NewDescribeCommand())

	return application
}
