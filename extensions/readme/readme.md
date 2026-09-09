# Provider Coderd

`provider-coderd` is a Crossplane provider for [Coder](https://coder.com),
generated with [Upjet](https://github.com/crossplane/upjet) from the
[coderd Terraform provider](https://github.com/coder/terraform-provider-coderd).

It manages Coder organizations, users, groups, templates, licenses, workspace
proxies, provisioner keys and the Coder Agents / AI configuration as Crossplane
managed resources.

The provider requires Crossplane v2 and serves every resource both
cluster-scoped (`*.coderd.crossplane.io`) and namespaced
(`*.coderd.m.crossplane.io`).

## Authentication

A `ProviderConfig` points at a `Secret` whose `credentials` key holds a JSON
document mirroring the coderd Terraform provider's configuration schema:

```json
{
  "url": "https://coder.example.com",
  "token": "<a Coder API token>"
}
```

`default_organization_id` and `headers` are optional. Most resource types need
a token with elevated permissions.
