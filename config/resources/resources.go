// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Package resources configures the API group, kind and cross-resource
// references of every managed resource this provider generates. It is shared
// by the cluster-scoped and the namespaced provider configurations, which only
// differ in their root API group.
package resources

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// The short API groups resources are generated into. They follow the
// subcategories the upstream provider documents its resources under.
const (
	groupIAM        = "iam"
	groupTemplate   = "template"
	groupDeployment = "deployment"
	groupAgents     = "agents"
)

// fieldOrganizationID is the Terraform field naming the organization a
// resource belongs to.
const fieldOrganizationID = "organization_id"

// groupKind is the API group and kind a Terraform resource is generated into.
type groupKind struct {
	group string
	kind  string
}

// groupKinds maps every Terraform resource of the coderd provider to its API
// group and kind.
var groupKinds = map[string]groupKind{
	"coderd_organization":               {groupIAM, "Organization"},
	"coderd_organization_group_sync":    {groupIAM, "OrganizationGroupSync"},
	"coderd_organization_sync_settings": {groupIAM, "OrganizationSyncSettings"},
	"coderd_user":                       {groupIAM, "User"},
	"coderd_group":                      {groupIAM, "Group"},

	"coderd_template": {groupTemplate, "Template"},

	"coderd_license":                  {groupDeployment, "License"},
	"coderd_workspace_proxy":          {groupDeployment, "WorkspaceProxy"},
	"coderd_provisioner_key":          {groupDeployment, "ProvisionerKey"},
	"coderd_oauth2_provider_settings": {groupDeployment, "OAuth2ProviderSettings"},

	"coderd_ai_provider":          {groupAgents, "AIProvider"},
	"coderd_agents_model":         {groupAgents, "Model"},
	"coderd_agents_default_model": {groupAgents, "DefaultModel"},
	"coderd_agents_system_prompt": {groupAgents, "SystemPrompt"},
	"coderd_agents_mcp_server":    {groupAgents, "MCPServer"},
}

// organizationRef references the organization a resource belongs to. Almost
// every Coder resource is scoped to one.
var organizationRef = ujconfig.Reference{TerraformName: "coderd_organization"}

// references maps a Terraform resource to the cross-resource references of its
// fields, keyed by Terraform field path.
var references = map[string]map[string]ujconfig.Reference{
	"coderd_group": {
		fieldOrganizationID: organizationRef,
		"members":           {TerraformName: "coderd_user"},
	},
	"coderd_organization_group_sync": {
		fieldOrganizationID: organizationRef,
	},
	"coderd_template": {
		fieldOrganizationID: organizationRef,
	},
	"coderd_provisioner_key": {
		fieldOrganizationID: organizationRef,
	},
	"coderd_agents_model": {
		fieldOrganizationID: organizationRef,
		"ai_provider_id":    {TerraformName: "coderd_ai_provider"},
	},
	"coderd_agents_default_model": {
		fieldOrganizationID: organizationRef,
		"model_id":          {TerraformName: "coderd_agents_model"},
	},
	"coderd_agents_mcp_server": {
		fieldOrganizationID: organizationRef,
	},
}

// Configure configures every resource of the coderd provider.
func Configure(p *ujconfig.Provider) {
	for name, gk := range groupKinds {
		refs := references[name]
		p.AddResourceConfigurator(name, func(r *ujconfig.Resource) {
			r.ShortGroup = gk.group
			r.Kind = gk.kind
			for field, ref := range refs {
				r.References[field] = ref
			}
		})
	}
}
