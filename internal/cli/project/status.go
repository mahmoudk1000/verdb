package project

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
	"github.com/mahmoudk1000/relen/internal/utils"
)

func NewStatusCommand() *cobra.Command {
	var queries *database.Queries

	status := &cobra.Command{
		Use:     "status",
		Aliases: []string{"st"},
		Args:    cobra.RangeArgs(1, 2),
		Short:   "Update project status",
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	flags := status.Flags()
	flags.Bool("json", false, "Output in JSON format")
	flags.BoolP("quiet", "q", false, "quiet mode, no output")

	status.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		ctx := cmd.Context()

		jsonFlag, _ := flags.GetBool("json")
		quietFlag, _ := flags.GetBool("quiet")

		if len(args) == 1 {
			s, err := getProjectStatus(ctx, args[0], queries)
			if err != nil {
				return err
			}
			if !quietFlag {
				if jsonFlag {
					fmtS, err := utils.FormatJSON(struct {
						Status string `json:"status"`
					}{
						Status: s,
					})
					if err != nil {
						return err
					}

					fmt.Println(fmtS)
					return nil
				}

				fmt.Println(s)
				return nil
			}
		} else {
			if err := updateProjectStatus(ctx, args[0], args[1], queries); err != nil {
				return err
			}
		}
		return nil
	}

	return status
}

func updateProjectStatus(ctx context.Context, pName string, s string, q *database.Queries) error {
	if _, err := q.CheckProjectExistsByName(ctx, pName); err != nil {
		return fmt.Errorf("project %s does not exist", pName)
	}

	if err := q.UpdateProjectStatusByName(ctx, database.UpdateProjectStatusByNameParams{
		Name:      pName,
		Status:    s,
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		return err
	}

	return nil
}

func getProjectStatus(ctx context.Context, pName string, q *database.Queries) (string, error) {
	pId, err := q.GetProjectIdByName(ctx, pName)
	if err != nil {
		return "", fmt.Errorf("project %s does not exist", pName)
	}

	s, err := q.GetProjectStatusById(ctx, pId)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve project status: %w", err)
	}

	return s, nil
}
