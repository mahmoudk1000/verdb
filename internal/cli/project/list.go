package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
)

func NewListCommand() *cobra.Command {
	var queries *database.Queries

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all projects",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			queries = db.Get()
			return nil
		},
	}

	flags := list.Flags()
	flags.Bool("json", false, "Output in JSON format")

	list.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		ctx := cmd.Context()

		jsonFlag, _ := flags.GetBool("json")

		ps, err := listProjects(ctx, queries)
		if err != nil {
			return err
		}

		switch {
		case jsonFlag:
			// TODO: implement JSON output
			return nil
		default:
			for _, pName := range ps {
				fmt.Println(pName)
			}
		}

		return nil
	}

	return list
}

func listProjects(ctx context.Context, q *database.Queries) ([]string, error) {
	ps, err := q.ListAllProjects(ctx)
	if err != nil {
		return nil, err
	}

	return ps, nil
}
