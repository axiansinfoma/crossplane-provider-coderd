// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/crossplane/upjet/v2/pkg/pipeline"

	"github.com/axiansinfoma/crossplane-provider-coderd/config"
	"github.com/axiansinfoma/crossplane-provider-coderd/internal/generated"
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
	// bumped, and what schema-diff-issues subtracts from the Terraform schema
	// to find resources we do not expose yet. `make generated-lst` refreshes it
	// on its own, without needing Terraform.
	names := generated.Names(provider.Resources)
	if err := generated.Write(names, filepath.Join(absRootDir, "config", "generated.lst")); err != nil {
		panic(fmt.Sprintf("cannot write the list of generated resources: %s", err))
	}
}
