# AGENTS.md

Context for coding agents working on `provider-coderd`.

## What this repository is

A [Crossplane](https://crossplane.io) provider for [Coder](https://coder.com),
**generated with [Upjet](https://github.com/crossplane/upjet)** from the
[coderd Terraform provider](https://github.com/coder/terraform-provider-coderd).

Two facts shape almost every task here:

1. **Most of the tree is generated output, not source.** Editing a generated
   file is always wrong — change the generator input and re-run `make generate`.
   CI's `check-diff` job re-generates and fails if the tree is dirty.
2. **The Terraform provider is linked into the binary, not exec'd.** The
   upstream provider's `fwprovider` package is a Go dependency, driven
   in-process through Upjet's Terraform plugin framework client. There is no
   Terraform CLI and no provider plugin in the image, and no subprocess per
   reconcile. Terraform is used *only* at code generation time, to dump the
   JSON schema.

Module: `github.com/axiansinfoma/crossplane-provider-coderd`. Go 1.26.
Requires Crossplane **v2** — every managed resource is served twice, as a
cluster-scoped type under `*.coderd.crossplane.io` and a namespaced type under
`*.coderd.m.crossplane.io`.

## Source vs. generated

**Hand-written (edit these):**

| Path | What it is |
| --- | --- |
| `config/external_name.go` | External-name config per Terraform resource. **A resource is only generated if it has an entry here** — this doubles as the include list. |
| `config/resources/resources.go` | API group, Kind and cross-resource references per resource. Shared by both provider variants. |
| `config/provider.go` | `GetProvider()` / `GetProviderNamespaced()` — the two Upjet provider configurations. |
| `config/fwprovider.go` | Constructs the embedded coderd provider instance. |
| `internal/clients/coderd.go` | `ProviderConfig` → Terraform provider configuration; credential resolution for both legacy (cluster) and modern (namespaced) managed resources. |
| `cmd/provider/main.go` | Controller manager entrypoint. |
| `cmd/generator/main.go` | Runs the Upjet pipeline. |
| `cmd/generatedlist/main.go` | Writes/verifies `config/generated.lst` without Terraform. |
| `internal/generated/list.go` | Helpers for that list. |
| `internal/features/features.go` | Feature flags. |
| `apis/*/v1alpha1/`, `apis/*/v1beta1/` (non-`zz_` files) | `ProviderConfig` / `ClusterProviderConfig` / `ProviderConfigUsage` types. |
| `apis/generate.go` | The `go:generate` chain the whole pipeline runs through. |
| `examples/` | Hand-written example manifests, one per kind, cluster + namespaced. Also what `make e2e` / uptest runs. |
| `docs/` | Hugo site (Hextra theme) with its own `go.mod`, deployed to GitHub Pages. |
| `hack/`, `scripts/`, `cluster/`, `.github/` | Automation. |

**Generated (never hand-edit):**

- `apis/**/zz_*.go` — managed resource types, deepcopy, managed methodsets.
- `internal/controller/**/zz_*.go` — one controller per kind per scope, plus `zz_setup.go`.
- `package/crds/` — CRD manifests.
- `examples-generated/` — example manifests emitted by Upjet.
- `config/schema.json` — `terraform providers schema` output.
- `config/provider-metadata.yaml` — scraped from the upstream provider's docs.
- `config/generated.lst` — sorted list of generated Terraform resource names.

All of these **are committed**; `check-diff` enforces that they match a fresh
generate.

`build/` is a **git submodule** (crossplane/build). Run
`git submodule update --init --recursive` (or just `make`, which self-heals via
the `fallthrough` target) before anything else on a fresh clone.

## Managed resources

15 Terraform resources across four short API groups, each generated into both
root groups:

| Group | Kinds |
| --- | --- |
| `iam` | `Organization`, `OrganizationGroupSync`, `OrganizationSyncSettings`, `User`, `Group` |
| `template` | `Template` |
| `deployment` | `License`, `WorkspaceProxy`, `ProvisionerKey`, `OAuth2ProviderSettings` |
| `agents` | `AIProvider`, `Model`, `DefaultModel`, `SystemPrompt`, `MCPServer` |

External names come in three flavours (see the comments in
`config/external_name.go`):

- `config.IdentifierFromProvider` — resources with a Coder-assigned `id`.
- `attributeAsExternalName("...")` — resources with **no** `id` attribute,
  addressed by another field (`coderd_organization_group_sync` by
  `organization_id`, `coderd_provisioner_key` by `name`).
- `singletonExternalName("...")` — deployment-wide singletons with a fixed
  placeholder ID matching the upstream `terraform import` ID
  (`coderd_agents_system_prompt`, `coderd_oauth2_provider_settings`,
  `coderd_organization_sync_settings`).

## Common tasks

### Add a managed resource for a new upstream resource

New upstream resources are **not** picked up automatically. The
`schema-diff-issues` workflow files an issue for each one that is missing.

1. Add an entry to `ExternalNameConfigs` in `config/external_name.go`.
2. Add the `groupKind` (and any `references`) in `config/resources/resources.go`.
3. `make generate`.
4. Add `examples/cluster/<group>/<kind>.yaml` and
   `examples/namespaced/<group>/<kind>.yaml`.
5. Update the table in `README.md` and `docs/content/docs/resources.md`.

### Bump the upstream Terraform provider

`Makefile`'s `TERRAFORM_PROVIDER_VERSION` and the `go.mod` dependency on
`github.com/coder/terraform-provider-coderd` **must move together** — one pins
the schema the CRDs come from, the other the implementation that is linked in.
Use the script, which does both plus the regeneration:

```console
./hack/update-terraform-provider.sh            # latest release
./hack/update-terraform-provider.sh 0.0.26     # a specific one
./hack/update-terraform-provider.sh --dry-run  # report only
```

The `update-terraform-provider` workflow runs this weekly and opens a PR with a
breaking-change report.

### Change credential handling or provider configuration

`internal/clients/coderd.go`. The credentials Secret holds a JSON document
mirroring the coderd provider's configuration schema: `url` and `token`
required, `default_organization_id` and `headers` optional.

## Commands

```console
make generate          # full regeneration (needs Terraform + docs checkout, slow)
make generated-lst     # refresh config/generated.lst only (no Terraform needed)
make generated-lst-check
make build             # build the provider binary
make lint              # golangci-lint v2
make test              # unit tests
make reviewable        # generate + lint + test — run before opening a PR
make check-diff        # generate, then fail if the tree is dirty (what CI runs)
make run               # run out-of-cluster against the current kubecontext
make local-deploy      # build the package and deploy it onto a kind cluster
make e2e               # local-deploy + uptest
make crddiff           # breaking CRD schema changes vs. the PR base branch
make schema-version-diff
```

`make generate` needs:

- **Terraform 1.5.7 or older.** 1.6+ is BSL-licensed and the Makefile refuses
  it (`check-terraform-version`). It is downloaded into `.cache/tools`.
- A sparse checkout of the upstream provider's `docs/resources` under
  `.work/coder/coderd` (`make pull-docs`).
- `goimports` on `PATH` (the Upjet generator shells out to it);
  `hack/update-terraform-provider.sh` builds it if missing.

## Conventions and gotchas

- **Commit messages use Conventional Commits** (`feat:`, `fix:`, `build:`,
  `docs:`, `chore:`), even though some of the early history does not.
- **There are currently no `_test.go` files in the repo.** `make test` passes
  trivially. If you add behaviour to `internal/clients` or `config`, adding the
  first tests is welcome, but don't assume an existing test suite covers you.
- Lint: golangci-lint **v2** config in `.golangci.yml`; `zz_*.go` and
  `examples/` are excluded. `goimports` local prefix is the module path, so
  local imports go in their own final group.
- Every source file carries the SPDX header
  (`// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>` /
  `// SPDX-License-Identifier: Apache-2.0`); `hack/boilerplate.go.txt` is what
  generated files get.
- Both provider configurations set `WithIncludeList([]string{})` deliberately —
  nothing is reconciled through the Terraform CLI, and the default include list
  would match every resource a second time. Don't "fix" this.
- `apis/generate.go` contains a `sed` that strips
  `" Resource - terraform-provider-coderd"` from `provider-metadata.yaml`,
  working around the scraper keying docs pages without examples by their full
  `page_title`. Keep it when touching the generate chain.
- Changes under `package/crds/` on a PR trigger `make crddiff`; a bump to
  `TERRAFORM_PROVIDER_VERSION` triggers `make schema-version-diff`. Both only
  report — they don't fail the build.
- Images and packages publish to `ghcr.io/axiansinfoma/provider-coderd`.

## CI

| Workflow | Trigger | What it does |
| --- | --- | --- |
| `ci.yml` | push to `main`/`release-*`, PRs | breaking-change reports, lint, `check-diff`, unit tests, `local-deploy`, then build and publish artifacts |
| `e2e.yaml` | `/test-examples` comment on a PR (write access required) | uptest against the example manifests |
| `update-terraform-provider.yml` | weekly Mon 06:00 UTC, manual | bump + regenerate + open PR |
| `schema-diff-issues.yml` | weekly Mon 06:00 UTC, `Makefile` push to `main`, manual | file an issue per unexposed upstream resource |
| `deploy-docs.yml` / `preview-docs.yml` | `docs/**` changes | Hugo build to GitHub Pages |
| `tag.yaml` | manual | cuts a release tag via the `crossplane-contrib/provider-workflows` reusable workflow |
| `publish-provider-package.yml` | manual | publishes the package to GHCR and the xpkg mirror, same reusable-workflow set |
| `backport.yml` | PR merged | opens backport PRs against `release-*` branches |

`make e2e` and the e2e workflow need a reachable Coder deployment; credentials
come from `UPTEST_CLOUD_CREDENTIALS` and are wired up by
`cluster/test/setup.sh`.
