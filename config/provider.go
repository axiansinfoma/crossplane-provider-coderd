// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	"github.com/axiansinfoma/crossplane-provider-coderd/config/resources"
)

const (
	resourcePrefix = "coderd"
	modulePath     = "github.com/axiansinfoma/crossplane-provider-coderd"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration.
//
// Every resource is reconciled through upjet's Terraform plugin framework
// client, which drives the embedded coderd provider in-process. The Terraform
// JSON schema is still what the CRDs are generated from; the embedded provider
// supplies the resource implementations.
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("coderd.crossplane.io"),
		ujconfig.WithTerraformPluginFrameworkProvider(NewFrameworkProvider()),
		ujconfig.WithTerraformPluginFrameworkIncludeList(ExternalNameConfigured()),
		// Nothing is reconciled through the Terraform CLI; the default include
		// list would otherwise match every resource a second time.
		ujconfig.WithIncludeList([]string{}),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		resources.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns the namespaced provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("coderd.m.crossplane.io"),
		ujconfig.WithTerraformPluginFrameworkProvider(NewFrameworkProvider()),
		ujconfig.WithTerraformPluginFrameworkIncludeList(ExternalNameConfigured()),
		// Nothing is reconciled through the Terraform CLI; the default include
		// list would otherwise match every resource a second time.
		ujconfig.WithIncludeList([]string{}),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}))

	for _, configure := range []func(provider *ujconfig.Provider){
		resources.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
