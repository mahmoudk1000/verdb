/*
Copyright © 2026 mahmoudk1000 <mahmoudk1000@gmail.com>
*/
package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/utils"
)

type outputFormat int

const (
	formatDefault outputFormat = iota
	formatJSON
	formatYAML
)

func getOutputFormat(cmd *cobra.Command) (outputFormat, error) {
	jsonFlag, err := cmd.Flags().GetBool(flagJSON)
	if err != nil {
		return formatDefault, err
	}

	yamlFlag, err := cmd.Flags().GetBool(flagYAML)
	if err != nil {
		return formatDefault, err
	}

	if jsonFlag && yamlFlag {
		return formatDefault, fmt.Errorf("cannot use both --json and --yaml flags")
	}

	if jsonFlag {
		return formatJSON, nil
	}
	if yamlFlag {
		return formatYAML, nil
	}

	return formatDefault, nil
}

func formatOutput(data any, format outputFormat) (string, error) {
	switch format {
	case formatJSON:
		return utils.FormatJSON(data)
	case formatYAML:
		return utils.FormatYAML(data)
	default:
		return utils.Format(data)
	}
}

func formatAndPrint(cmd *cobra.Command, data any) error {
	format, err := getOutputFormat(cmd)
	if err != nil {
		return err
	}

	output, err := formatOutput(data, format)
	if err != nil {
		return err
	}

	fmt.Println(output)
	return nil
}

func ensureProjectExists(ctx context.Context, name string, q *database.Queries) error {
	exists, err := q.CheckProjectExistsByName(ctx, name)
	if err != nil {
		return fmt.Errorf(checkProjectExistsErr, err)
	}

	if !exists {
		return fmt.Errorf(projectNotFoundErr, name)
	}

	return nil
}

func ensureProjectNotExists(ctx context.Context, name string, q *database.Queries) error {
	exists, err := q.CheckProjectExistsByName(ctx, name)
	if err != nil {
		return fmt.Errorf(checkProjectExistsErr, err)
	}

	if exists {
		return fmt.Errorf(projectExistsErr, name)
	}

	return nil
}

func addOutputFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.Bool(flagJSON, false, "Output in JSON format")
	flags.Bool(flagYAML, false, "Output in YAML format")
}

func isQuietMode(cmd *cobra.Command) bool {
	quiet, err := cmd.Flags().GetBool(flagQuiet)
	if err != nil {
		return false
	}
	return quiet
}
