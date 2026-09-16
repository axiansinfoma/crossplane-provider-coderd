---
title: Managed resources
weight: 3
---

Every resource the coderd Terraform provider supports is exposed as a managed
resource. Each kind below is served twice: cluster-scoped under the group
shown, and namespaced under the same group with `.m.` inserted — for example
`iam.coderd.crossplane.io` and `iam.coderd.m.crossplane.io`.

All managed resource kinds are at `v1alpha1`; the `ProviderConfig` kinds are at
`v1beta1`.

## `iam.coderd.crossplane.io`

| Kind | Terraform resource |
| --- | --- |
| `Organization` | `coderd_organization` |
| `OrganizationGroupSync` | `coderd_organization_group_sync` |
| `OrganizationSyncSettings` | `coderd_organization_sync_settings` |
| `Group` | `coderd_group` |
| `User` | `coderd_user` |

## `template.coderd.crossplane.io`

| Kind | Terraform resource |
| --- | --- |
| `Template` | `coderd_template` |

## `deployment.coderd.crossplane.io`

| Kind | Terraform resource |
| --- | --- |
| `License` | `coderd_license` |
| `WorkspaceProxy` | `coderd_workspace_proxy` |
| `ProvisionerKey` | `coderd_provisioner_key` |
| `OAuth2ProviderSettings` | `coderd_oauth2_provider_settings` |

## `agents.coderd.crossplane.io`

| Kind | Terraform resource |
| --- | --- |
| `AIProvider` | `coderd_ai_provider` |
| `Model` | `coderd_agents_model` |
| `DefaultModel` | `coderd_agents_default_model` |
| `SystemPrompt` | `coderd_agents_system_prompt` |
| `MCPServer` | `coderd_agents_mcp_server` |

## Field schemas

This page is deliberately a map, not a reference. The complete, always-current
field schema of every kind lives in the CRDs themselves:

```console
kubectl explain organization.iam.coderd.crossplane.io --recursive
```

or read
[`package/crds/`](https://github.com/axiansinfoma/crossplane-provider-coderd/tree/main/package/crds)
in the repository. For what each field *means*, the
[coderd Terraform provider documentation](https://registry.terraform.io/providers/coder/coderd/latest/docs)
is the upstream source the schema is generated from.

## Cross-resource references

Where a field takes the ID of another Coder object, the generated kind also
accepts a reference to the managed resource that owns it, so Crossplane resolves
the ID for you:

```yaml
apiVersion: iam.coderd.crossplane.io/v1alpha1
kind: Group
metadata:
  name: platform
spec:
  forProvider:
    name: platform
    membersRefs:
      - name: alice
      - name: bob
```

`membersRefs` resolves against `User` resources. Which fields have a reference
is decided in
[`config/resources/resources.go`](https://github.com/axiansinfoma/crossplane-provider-coderd/blob/main/config/resources/resources.go).
