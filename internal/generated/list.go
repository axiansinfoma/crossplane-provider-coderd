// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Package generated maintains config/generated.lst: the sorted list of
// Terraform resource names this provider generates managed resources for.
//
// It is what `make schema-version-diff` reports Terraform schema changes
// against when the upstream provider is bumped, and what the schema-diff-issues
// automation subtracts from config/schema.json to find resources the upstream
// provider exposes that we do not.
package generated

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
)

// Names returns the sorted keys of a provider's resource map.
func Names[T any](resources map[string]T) []string {
	names := make([]string, 0, len(resources))
	for name := range resources {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Marshal renders names in the on-disk representation of the list.
func Marshal(names []string) ([]byte, error) {
	out, err := json.Marshal(names)
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}

// Write writes names to path.
func Write(names []string, path string) error {
	out, err := Marshal(names)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o600)
}

// UpToDate reports whether the list at path already holds exactly names.
func UpToDate(names []string, path string) (bool, error) {
	want, err := Marshal(names)
	if err != nil {
		return false, err
	}
	// The path is a build-time argument from the Makefile, not user input.
	got, err := os.ReadFile(path) //nolint:gosec // developer tooling reading a path given on the command line
	if err != nil {
		return false, err
	}
	return bytes.Equal(bytes.TrimSpace(got), bytes.TrimSpace(want)), nil
}
