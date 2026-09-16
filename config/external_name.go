// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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
	"coderd_organization":         uuidIdentified(),
	"coderd_user":                 uuidIdentified(),
	"coderd_group":                uuidIdentified(),
	"coderd_template":             uuidIdentified(),
	"coderd_workspace_proxy":      uuidIdentified(),
	"coderd_ai_provider":          uuidIdentified(),
	"coderd_agents_model":         uuidIdentified(),
	"coderd_agents_default_model": uuidIdentified(),
	"coderd_agents_mcp_server":    uuidIdentified(),

	// Identified by the Coder-assigned numeric license ID.
	"coderd_license": numericIdentified("0"),

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

// uuidIdentified returns the external name configuration for a resource whose
// Terraform "id" attribute is a Coder-assigned UUID.
//
// upjet's plugin framework Observe calls ReadResource unconditionally, even
// before the resource has been created, when the external name annotation --
// and therefore the "id" in the Terraform state -- is still empty. The coderd
// provider types these ids with a UUID-validating framework type, so the read
// fails with an error-severity diagnostic instead of returning a null state,
// and the managed resource never advances to Create. Treating that specific
// diagnostic as "resource not found" restores the intended flow.
func uuidIdentified() config.ExternalName {
	e := config.IdentifierFromProvider
	e.IsNotFoundDiagnosticFn = isEmptyUUIDDiagnostic
	return e
}

// numericIdentified returns the external name configuration for a resource
// whose Terraform "id" attribute is a Coder-assigned number rather than a
// string.
//
// These resources hit the same pre-create problem as uuidIdentified, but fail
// earlier and for a different reason. upjet seeds the Terraform state from the
// managed resource's parameters, and on the first reconcile copyParameters
// copies params["id"] -- the empty external name -- into that state. Decoding
// the empty string into a numeric attribute fails while the state is still
// being constructed ("error parsing number: EOF"), so the read never runs and
// no diagnostic hook can intercept it. Handing out a sentinel id instead keeps
// the state decodable; the sentinel matches no real license, so the provider's
// Read reports the resource as gone and upjet proceeds to Create.
func numericIdentified(sentinel string) config.ExternalName {
	e := config.IdentifierFromProvider
	e.GetIDFn = func(ctx context.Context, externalName string, params, tfstate map[string]any) (string, error) {
		if externalName == "" {
			return sentinel, nil
		}
		return config.ExternalNameAsID(ctx, externalName, params, tfstate)
	}
	return e
}

// diagSummaryInvalidUUID is the diagnostic summary the coderd provider returns
// when an attribute typed as a UUID cannot be parsed.
const diagSummaryInvalidUUID = "Invalid UUID"

// detailEmptyUUID is the parse error a UUID attribute produces when it holds
// the empty string, i.e. when the resource has not been created yet. We match
// on the zero length specifically so that a genuinely malformed UUID still
// surfaces as an error rather than being silently reported as not found.
const detailEmptyUUID = "invalid UUID length: 0"

func isEmptyUUIDDiagnostic(diags []*tfprotov6.Diagnostic) bool {
	for _, d := range diags {
		if d.Severity != tfprotov6.DiagnosticSeverityError {
			continue
		}
		if d.Summary == diagSummaryInvalidUUID && strings.Contains(d.Detail, detailEmptyUUID) {
			return true
		}
	}
	return false
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
