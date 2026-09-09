// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/crossplane/upjet/v2/pkg/pipeline"

	"github.com/axiansinfoma/crossplane-provider-coderd/config"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		panic("root directory is required to be given as argument")
	}
	rootDir := os.Args[1]
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		panic(fmt.Sprintf("cannot calculate the absolute path with %s", rootDir))
	}
	provider := config.GetProvider()
	pipeline.Run(provider, config.GetProviderNamespaced(), absRootDir)

	// The list of generated resources is what `make schema-version-diff`
	// reports Terraform schema changes against when the upstream provider is
	// bumped.
	if err := writeGeneratedList(provider.Resources, filepath.Join(absRootDir, "config", "generated.lst")); err != nil {
		panic(fmt.Sprintf("cannot write the list of generated resources: %s", err))
	}
}

func writeGeneratedList[T any](resources map[string]T, path string) error {
	names := make([]string, 0, len(resources))
	for name := range resources {
		names = append(names, name)
	}
	sort.Strings(names)
	out, err := json.Marshal(names)
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(out, '\n'), 0o600)
}
