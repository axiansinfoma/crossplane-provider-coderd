# provider-coderd

`provider-coderd` is a [Crossplane](https://crossplane.io/) provider for
[Coder](https://coder.com). It is generated with
[Upjet](https://github.com/crossplane/upjet) from the
[coderd Terraform provider](https://github.com/coder/terraform-provider-coderd)
and exposes every resource that provider supports as a Crossplane managed
resource.

It requires **Crossplane v2** and serves each managed resource twice: as a
cluster-scoped type under `*.coderd.crossplane.io` and as a namespaced type
under `*.coderd.m.crossplane.io`.

## Managed resources

| API group | Kinds |
| --- | --- |
| `iam.coderd.crossplane.io` | `Organization`, `OrganizationGroupSync`, `OrganizationSyncSettings`, `User`, `Group` |
| `template.coderd.crossplane.io` | `Template` |
| `deployment.coderd.crossplane.io` | `License`, `WorkspaceProxy`, `ProvisionerKey`, `OAuth2ProviderSettings` |
| `agents.coderd.crossplane.io` | `AIProvider`, `Model`, `DefaultModel`, `SystemPrompt`, `MCPServer` |

Each group also exists in its namespaced form, e.g.
`iam.coderd.m.crossplane.io`.

## Getting started

Install the provider:

```console
kubectl apply -f examples/install.yaml
```

Create a `Secret` holding the credentials for your Coder deployment. The keys
mirror the configuration schema of the coderd Terraform provider — `url` and
`token` are required, `default_organization_id` and `headers` are optional:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: coderd-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    {
      "url": "https://coder.example.com",
      "token": "<a Coder API token>"
    }
```

Most resource types need a token with elevated permissions. Then point a
`ProviderConfig` at it and create resources — see [examples/](examples/) for a
manifest per kind, in both the cluster-scoped and the namespaced flavour.

## How this provider is generated

The provider runs the coderd Terraform provider under the Terraform CLI. Both
binaries are baked into the controller image, and the controller drives a
Terraform workspace per managed resource.

Everything under `apis/`, `internal/controller/`, `package/crds/` and
`examples-generated/` is generated. The inputs are:

- `Makefile` — `TERRAFORM_PROVIDER_VERSION` pins the upstream release.
- `config/schema.json` — the Terraform schema, from `terraform providers schema`.
- `config/provider-metadata.yaml` — scraped from the upstream provider's docs.
- `config/external_name.go` — how each resource's external name maps to its
  Terraform ID.
- `config/resources/resources.go` — the API group, kind and cross-resource
  references of each resource.

Regenerate with:

```console
make generate
```

## Staying in step with the Terraform provider

`.github/workflows/update-terraform-provider.yml` checks weekly for a new
release of the coderd Terraform provider, regenerates the provider and opens a
pull request. Run the same thing locally with:

```console
# update to the latest upstream release
./hack/update-terraform-provider.sh

# or to a specific one
./hack/update-terraform-provider.sh 0.0.26
```

A new **resource** in the upstream provider is not picked up automatically: add
it to `config/external_name.go` and `config/resources/resources.go` first, then
regenerate.

## Developing

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

Run code generation, linters and tests:

```console
make reviewable
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/axiansinfoma/crossplane-provider-coderd/issues).
