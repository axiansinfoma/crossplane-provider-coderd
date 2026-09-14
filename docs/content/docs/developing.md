---
title: Developing
weight: 4
---

## How the provider is generated

The provider is generated with [Upjet](https://github.com/crossplane/upjet) from
the coderd Terraform provider. Everything under `apis/`, `internal/controller/`,
`package/crds/` and `examples-generated/` is output, not source.

The coderd provider is linked into the controller binary through its
[`fwprovider`](https://github.com/coder/terraform-provider-coderd/tree/main/fwprovider)
package and driven in-process by Upjet's Terraform plugin framework client.
There is no Terraform CLI and no provider plugin in the image, and no subprocess
per reconcile. Terraform is used only at code generation time, to dump the
provider's JSON schema.

The inputs are:

| Input | What it controls |
| --- | --- |
| `Makefile` | `TERRAFORM_PROVIDER_VERSION` pins the upstream release the schema and docs come from. |
| `go.mod` | Pins the same release of `github.com/coder/terraform-provider-coderd` as the embedded implementation. It must move with `TERRAFORM_PROVIDER_VERSION`. |
| `config/schema.json` | The Terraform schema, from `terraform providers schema`. |
| `config/provider-metadata.yaml` | Scraped from the upstream provider's documentation. |
| `config/external_name.go` | How each resource's external name maps to its Terraform ID. |
| `config/resources/resources.go` | The API group, kind and cross-resource references of each resource. |
| `config/generated.lst` | The resources that are actually generated. Written by `make generate`, refreshable on its own with `make generated-lst`. |

Regenerate with:

```console
make generate
```

## Adding a resource

A resource in `config/schema.json` is only generated once it has an entry in
`config/external_name.go`. The `schema-diff-issues` workflow files an issue for
every upstream resource that does not, so the work is usually already tracked.

1. Add an entry for the resource to `config/external_name.go`.
2. Give it an API group, kind and any cross-resource references in
   `config/resources/resources.go`.
3. Run `make generate`.
4. Add an example manifest under `examples/cluster/<group>/` and
   `examples/namespaced/<group>/`.
5. Add it to the [managed resources](../resources) page and the table in
   `README.md`.

## Keeping up with the upstream provider

The `update-terraform-provider` workflow runs weekly. When a new coderd provider
release appears it bumps `TERRAFORM_PROVIDER_VERSION` and the `go.mod`
dependency together, regenerates everything, files a tracking issue with the
remaining checklist and opens a pull request whose body carries a
breaking-change report from `make crddiff` and `make schema-version-diff`.

Run it by hand against a specific version with:

```console
hack/update-terraform-provider.sh 0.0.26
```

or see what it would do without touching the tree:

```console
hack/update-terraform-provider.sh --dry-run
```

## Testing

```console
make lint          # golangci-lint
make test          # unit tests
make local-deploy  # build and run the provider on a kind cluster
make e2e           # local-deploy, then the uptest suite
```

`make e2e` needs a reachable Coder deployment and credentials in
`UPTEST_DATASOURCE_PATH`. On a pull request, comment `/test-examples` to run it
in CI.
