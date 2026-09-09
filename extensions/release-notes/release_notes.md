# Initial release

First release of `provider-coderd`, generated from the coderd Terraform
provider v0.0.25.

Covers all 15 resources the upstream provider supports, in both the
cluster-scoped and the namespaced Crossplane v2 flavour:

- **IAM** — `Organization`, `OrganizationGroupSync`, `OrganizationSyncSettings`,
  `User`, `Group`
- **Templates** — `Template`
- **Deployment** — `License`, `WorkspaceProxy`, `ProvisionerKey`,
  `OAuth2ProviderSettings`
- **Agents** — `AIProvider`, `Model`, `DefaultModel`, `SystemPrompt`,
  `MCPServer`
