/*
Copyright © 2026 mahmoudk1000 <mahmoudk1000@gmail.com>
*/
package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/cli/application"
	"github.com/mahmoudk1000/relen/internal/cli/project"
)

func main() {
	var relen = &cobra.Command{
		Use:   "relen",
		Short: "A serious, well-scoped versioning tool.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	relen.AddCommand(project.NewProjectCommand())
	relen.AddCommand(application.NewApplicationCommand())

	if err := relen.Execute(); err != nil {
		os.Exit(1)
	}
}
