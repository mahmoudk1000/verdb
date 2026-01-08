/*
Copyright © 2026 mahmoudk1000 <mahmoudk1000@gmail.com>
*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/cli/application"
	"github.com/mahmoudk1000/relen/internal/cli/config"
	"github.com/mahmoudk1000/relen/internal/cli/project"
	"github.com/mahmoudk1000/relen/internal/db"
)

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

func init() {
	relen.AddCommand(project.NewProjectCommand())
	relen.AddCommand(application.NewApplicationCommand())
	relen.AddCommand(config.NewConfigCommand())
}

func main() {
	dbENV := os.Getenv("RELEN_DATABASE_URL")
	if dbENV == "" {
		fmt.Println("RELEN_DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	if err := db.Init(dbENV); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("Failed to close database connection: %v\n", err)
			os.Exit(1)
		}
	}()

	if err := relen.Execute(); err != nil {
		os.Exit(1)
	}
}
