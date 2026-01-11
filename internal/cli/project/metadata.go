/*
Copyright © 2026 mahmoudk1000 <mahmoudk1000@gmail.com>
*/
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

type metadataOptions struct {
	projectName string
	key         string
	setValue    string
	deleteKey   string
	isSet       bool
	isDelete    bool
}

func (o *metadataOptions) validate() error {
	if o.isSet && o.isDelete {
		return fmt.Errorf("cannot use --set and --delete flags together")
	}

	if o.isSet {
		parts := strings.SplitN(o.setValue, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf(invalidFormatErr, o.setValue)
		}
		if strings.TrimSpace(parts[0]) == "" {
			return fmt.Errorf(keyCannotBeEmptyErr)
		}
		if strings.TrimSpace(parts[1]) == "" {
			return fmt.Errorf(valueCannotBeEmptyErr)
		}
	}

	if o.isDelete && strings.TrimSpace(o.deleteKey) == "" {
		return fmt.Errorf(keyCannotBeEmptyErr)
	}

	return nil
}

func NewMetadataCommand() *cobra.Command {
	opts := &metadataOptions{}
	var queries *database.Queries

	metadata := &cobra.Command{
		Use:     "metadata <project-name> [key]",
		Aliases: []string{"md", "meta"},
		Short:   "Manage project metadata",
		Args:    cobra.RangeArgs(1, 2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			queries = db.Get()
			opts.projectName = args[0]

			if len(args) == 2 {
				opts.key = args[1]
			}

			if cmd.Flags().Changed("set") {
				opts.isSet = true
			}

			if cmd.Flags().Changed("delete") {
				opts.isDelete = true
			}

			return opts.validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := cmd.Context()

			if err := ensureProjectExists(ctx, opts.projectName, queries); err != nil {
				return err
			}

			if opts.isSet {
				return runMetadataSet(ctx, opts, queries)
			}
			if opts.isDelete {
				return runMetadataDelete(ctx, opts, queries)
			}
			return runMetadataGet(ctx, cmd, opts, queries)
		},
	}

	flags := metadata.Flags()
	flags.StringVarP(&opts.setValue, "set", "s", "",
		"Set metadata value in key=value format")
	flags.StringVarP(&opts.deleteKey, "delete", "d", "",
		"Delete metadata key")
	addOutputFlags(metadata)

	return metadata
}

func runMetadataGet(
	ctx context.Context,
	cmd *cobra.Command,
	opts *metadataOptions,
	q *database.Queries,
) error {
	metadataStr, err := q.GetProjectMetadata(ctx, opts.projectName)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	metadata, err := parseMetadataString(metadataStr)
	if err != nil {
		return err
	}

	if opts.key != "" {
		return displayMetadataKey(cmd, opts.key, metadata)
	}

	return displayAllMetadata(cmd, metadata)
}

func runMetadataSet(
	ctx context.Context,
	opts *metadataOptions,
	q *database.Queries,
) error {
	parts := strings.SplitN(opts.setValue, "=", 2)
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	metadataStr, err := q.GetProjectMetadata(ctx, opts.projectName)
	if err != nil {
		return fmt.Errorf("failed to get existing metadata: %w", err)
	}

	metadata, err := parseMetadataString(metadataStr)
	if err != nil {
		return err
	}

	metadata[key] = value

	metadataJSON, err := utils.MetadataToJSON(metadata)
	if err != nil {
		return fmt.Errorf(failedToParseMetadataErr, err)
	}

	params := database.UpdateProjectMetadataParams{
		Name:      opts.projectName,
		Metadata:  metadataJSON,
		UpdatedAt: time.Now().UTC(),
	}

	if err := q.UpdateProjectMetadata(ctx, params); err != nil {
		return fmt.Errorf(failedToUpdateProjectErr, err)
	}

	return nil
}

func runMetadataDelete(
	ctx context.Context,
	opts *metadataOptions,
	q *database.Queries,
) error {
	key := strings.TrimSpace(opts.deleteKey)

	metadataStr, err := q.GetProjectMetadata(ctx, opts.projectName)
	if err != nil {
		return fmt.Errorf("failed to get existing metadata: %w", err)
	}

	metadata, err := parseMetadataString(metadataStr)
	if err != nil {
		return err
	}

	if _, exists := metadata[key]; !exists {
		return fmt.Errorf(keyNotFoundInMetadataErr, key)
	}

	delete(metadata, key)

	metadataJSON, err := utils.MetadataToJSON(metadata)
	if err != nil {
		return fmt.Errorf(failedToParseMetadataErr, err)
	}

	params := database.UpdateProjectMetadataParams{
		Name:      opts.projectName,
		Metadata:  metadataJSON,
		UpdatedAt: time.Now().UTC(),
	}

	if err := q.UpdateProjectMetadata(ctx, params); err != nil {
		return fmt.Errorf(failedToUpdateProjectErr, err)
	}

	fmt.Printf("Successfully deleted metadata key %q from project %q\n",
		key, opts.projectName)
	return nil
}

func parseMetadataString(metadataStr string) (map[string]any, error) {
	metadata := make(map[string]any)

	if metadataStr == "" || metadataStr == "null" {
		return metadata, nil
	}

	if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
		return nil, fmt.Errorf(failedToParseMetadataErr, err)
	}

	return metadata, nil
}

func displayMetadataKey(cmd *cobra.Command, key string, metadata map[string]any) error {
	value, exists := metadata[key]
	if !exists {
		return fmt.Errorf(keyNotFoundInMetadataErr, key)
	}

	format, err := getOutputFormat(cmd)
	if err != nil {
		return err
	}

	if format == formatJSON {
		output, err := formatOutput(map[string]any{key: value}, formatJSON)
		if err != nil {
			return err
		}
		fmt.Println(output)
	} else if format == formatYAML {
		output, err := formatOutput(map[string]any{key: value}, formatYAML)
		if err != nil {
			return err
		}
		fmt.Println(output)
	} else {
		fmt.Printf("%s: %v\n", key, value)
	}

	return nil
}

func displayAllMetadata(cmd *cobra.Command, metadata map[string]any) error {
	if len(metadata) == 0 {
		fmt.Println("No metadata found")
		return nil
	}

	format, err := getOutputFormat(cmd)
	if err != nil {
		return err
	}

	if format == formatJSON || format == formatYAML {
		output, err := formatOutput(metadata, format)
		if err != nil {
			return err
		}
		fmt.Println(output)
	} else {
		for key, value := range metadata {
			fmt.Printf("%s: %v\n", key, value)
		}
	}

	return nil
}
