package config

import (
	"github.com/spf13/cobra"
)

func NewConfigCommand() *cobra.Command {
	config := &cobra.Command{
		Use:   "config",
		Short: "Initialize database for the CLI tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}

			return nil
		},
	}

	config.AddCommand(NewInitCommand())

	return config
}
