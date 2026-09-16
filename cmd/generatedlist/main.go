// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Command generatedlist writes config/generated.lst, the sorted list of
// Terraform resource names this provider generates managed resources for.
//
// The full generator needs the Terraform CLI and a checkout of the upstream
// provider's documentation. This one does not: the provider configuration is
// built from the schema and metadata embedded in the config package, so the
// list can be refreshed from a bare checkout.
//
// Usage:
//
//	generatedlist [--check] <path>
//
// With --check it writes nothing and exits non-zero when the file on disk is
// stale.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/axiansinfoma/crossplane-provider-coderd/config"
	"github.com/axiansinfoma/crossplane-provider-coderd/internal/generated"
)

func main() {
	check := flag.Bool("check", false, "verify the list on disk is up to date instead of writing it")
	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		fmt.Fprintln(os.Stderr, "usage: generatedlist [--check] <path>")
		os.Exit(2)
	}

	names := generated.Names(config.GetProvider().Resources)

	if *check {
		ok, err := generated.UpToDate(names, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read %s: %s\n", path, err)
			os.Exit(2)
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "%s is stale, run `make generated-lst` to refresh it\n", path)
			os.Exit(1)
		}
		fmt.Printf("%s is up to date with %d resources\n", path, len(names))
		return
	}

	if err := generated.Write(names, path); err != nil {
		fmt.Fprintf(os.Stderr, "cannot write %s: %s\n", path, err)
		os.Exit(2)
	}
	fmt.Printf("wrote %d resources to %s\n", len(names), path)
}
