package cli

// DONTCOVER

import (
	"strings"
)

func ParseStringSliceFromString(s string, separator string) ([]string, error) {
	if s == "" {
		return []string{}, nil
	}

	stringSlice := strings.Split(s, separator)

	parsedStrings := make([]string, 0, len(stringSlice))
	for _, s := range stringSlice {
		s = strings.TrimSpace(s)

		parsedStrings = append(parsedStrings, s)
	}
	return parsedStrings, nil
}
