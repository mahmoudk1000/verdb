package config

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
)

func NewInitCommand() *cobra.Command {
	var queries *database.Queries

	init := &cobra.Command{
		Use:   "init",
		Short: "Initialize database for the CLI tool",
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := queries.InitSchema(ctx); err != nil {
				return err
			}

			return nil
		},
	}

	return init
}
