package project

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/relen/internal/database"
	"github.com/mahmoudk1000/relen/internal/db"
	"github.com/mahmoudk1000/relen/internal/utils"
)

func NewMetadataCommand() *cobra.Command {
	var queries *database.Queries

	metadata := &cobra.Command{
		Use:     "metadata [project-name] [key]",
		Aliases: []string{"md"},
		Short:   "Manage project metadata",
		Args:    cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = db.Get()
		},
	}

	flags := metadata.Flags()
	flags.Bool("json", false, "Output in JSON format")
	flags.StringP("set", "s", "", "Set metadata value for the specified key in the format key=value")

	metadata.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		ctx := cmd.Context()

		projectName := args[0]
		flagSet, _ := cmd.Flags().GetString("set")
		flagJSON, _ := cmd.Flags().GetBool("json")

		if _, err := queries.GetProjectIdByName(ctx, projectName); err != nil {
			return fmt.Errorf(projectNotFoundErr, projectName)
		}

		if flagSet != "" {
			return setProjectMetadata(ctx, projectName, flagSet, queries)
		}

		var key string
		if len(args) > 1 {
			key = args[1]
		}

		return getProjectMetadata(ctx, projectName, key, flagJSON, queries)
	}

	return metadata
}

func getProjectMetadata(
	ctx context.Context,
	pName string,
	key string,
	asJSON bool,
	q *database.Queries,
) error {
	metadataStr, err := q.GetProjectMetadata(ctx, pName)
	if err != nil {
		return err
	}

	var metadata map[string]any
	if metadataStr != "" && metadataStr != "null" {
		if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
			return fmt.Errorf(failedToParseMetadataErr, err)
		}
	}

	if metadata == nil {
		metadata = make(map[string]any)
	}

	if key != "" {
		value, exists := metadata[key]
		if !exists {
			return fmt.Errorf(keyNotFoundInMetadataErr, key)
		}

		if asJSON {
			output, err := json.MarshalIndent(map[string]any{key: value}, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		} else {
			fmt.Printf("%s: %v\n", key, value)
		}
		return nil
	}

	if len(metadata) == 0 {
		fmt.Println("No metadata found")
		return nil
	}

	if asJSON {
		output, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	} else {
		for k, v := range metadata {
			fmt.Printf("%s: %v\n", k, v)
		}
	}

	return nil
}

func setProjectMetadata(
	ctx context.Context,
	pName, keyValue string,
	q *database.Queries,
) error {
	parts := strings.SplitN(keyValue, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf(invalidFormatErr, keyValue)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return fmt.Errorf(keyCannotBeEmptyErr)
	}

	metadataStr, err := q.GetProjectMetadata(ctx, pName)
	if err != nil {
		return err
	}

	var metadata map[string]any
	if metadataStr != "" && metadataStr != "null" {
		if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
			return fmt.Errorf(failedToParseMetadataErr, err)
		}
	}

	if metadata == nil {
		metadata = make(map[string]any)
	}

	metadata[key] = value

	metadataJSON, err := utils.MetadataToJSON(metadata)
	if err != nil {
		return err
	}

	if err := q.UpdateProjectMetadata(ctx, database.UpdateProjectMetadataParams{
		Name:      pName,
		Metadata:  metadataJSON,
		UpdatedAt: time.Now(),
	}); err != nil {
		return err
	}

	return nil
}
