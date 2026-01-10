package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sqlc-dev/pqtype"
)

func ParseProjectSlashApplication(args []string) (project, app string, err error) {
	if len(args) == 0 {
		return "", "", fmt.Errorf("no arguments provided")
	}

	if len(args) == 1 && strings.Contains(args[0], "/") {
		parts := strings.SplitN(args[0], "/", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid format: use 'project/app' or 'project app'")
		}

		return parts[0], parts[1], nil
	}

	if len(args) == 2 {
		return args[0], args[1], nil
	}

	return "", "", fmt.Errorf("invalid format: use 'project/app' or 'project app'")
}

func ParseMetadata(metadataSlice []string) (map[string]any, error) {
	if len(metadataSlice) == 0 {
		return nil, nil
	}

	metadata := make(map[string]any)

	for _, item := range metadataSlice {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid metadata format '%s': expected key=value", item)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("invalid metadata: empty key in '%s'", item)
		}

		metadata[key] = value
	}

	return metadata, nil
}

func MetadataToJSON(metadata map[string]any) (pqtype.NullRawMessage, error) {
	if len(metadata) == 0 {
		return pqtype.NullRawMessage{Valid: false}, nil
	}

	jsonBytes, err := json.Marshal(metadata)
	if err != nil {
		return pqtype.NullRawMessage{}, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return pqtype.NullRawMessage{
		RawMessage: jsonBytes,
		Valid:      true,
	}, nil
}
