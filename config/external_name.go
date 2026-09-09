// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/pkg/errors"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
//
// Every resource of the coderd Terraform provider whose schema carries an "id"
// attribute has that attribute assigned by the Coder deployment (a UUID, or a
// numeric ID for licenses), so config.IdentifierFromProvider is the right
// choice for them. The remaining resources are Terraform plugin framework
// resources with no "id" attribute at all, and get an external name derived
// from their own state below.
var ExternalNameConfigs = map[string]config.ExternalName{
	// Identified by a Coder-assigned UUID.
	"coderd_organization":         config.IdentifierFromProvider,
	"coderd_user":                 config.IdentifierFromProvider,
	"coderd_group":                config.IdentifierFromProvider,
	"coderd_template":             config.IdentifierFromProvider,
	"coderd_workspace_proxy":      config.IdentifierFromProvider,
	"coderd_ai_provider":          config.IdentifierFromProvider,
	"coderd_agents_model":         config.IdentifierFromProvider,
	"coderd_agents_default_model": config.IdentifierFromProvider,
	"coderd_agents_mcp_server":    config.IdentifierFromProvider,

	// Identified by the Coder-assigned numeric license ID.
	"coderd_license": config.IdentifierFromProvider,

	// No "id" attribute: the group sync configuration of an organization is
	// addressed by that organization's ID, which is also its import ID.
	"coderd_organization_group_sync": attributeAsExternalName("organization_id"),

	// No "id" attribute and no import support: provisioner keys are addressed
	// by their name within an organization.
	"coderd_provisioner_key": attributeAsExternalName("name"),

	// No "id" attribute: deployment-wide singletons. The external name is the
	// placeholder ID that "terraform import" expects upstream.
	"coderd_agents_system_prompt":       singletonExternalName("agents_system_prompt"),
	"coderd_oauth2_provider_settings":   singletonExternalName("oauth2_provider_settings"),
	"coderd_organization_sync_settings": singletonExternalName("organization_sync_settings"),
}

// attributeAsExternalName returns an external name configuration that reads the
// external name from the given attribute of the Terraform state instead of from
// the "id" attribute, which these resources do not have. The Terraform ID stays
// the external name, so it still matches the documented import ID.
func attributeAsExternalName(attr string) config.ExternalName {
	e := config.IdentifierFromProvider
	e.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
		v, ok := tfstate[attr]
		if !ok {
			return "", errors.Errorf("%q does not exist in tfstate", attr)
		}
		s, ok := v.(string)
		if !ok {
			return "", errors.Errorf("%q is not a string in tfstate", attr)
		}
		if s == "" {
			return "", errors.Errorf("%q is empty in tfstate", attr)
		}
		return s, nil
	}
	return e
}

// singletonExternalName returns an external name configuration for a
// deployment-wide singleton resource. There is only ever one instance of these
// on a Coder deployment, and their Terraform schema has no "id" attribute, so
// the external name is fixed.
func singletonExternalName(id string) config.ExternalName {
	e := config.IdentifierFromProvider
	e.GetExternalNameFn = func(map[string]any) (string, error) {
		return id, nil
	}
	e.GetIDFn = func(context.Context, string, map[string]any, map[string]any) (string, error) {
		return id, nil
	}
	return e
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
