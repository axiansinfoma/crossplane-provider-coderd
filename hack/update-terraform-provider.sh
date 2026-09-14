#!/usr/bin/env bash
# Updates the coderd Terraform provider release this Crossplane provider is
# generated from and embeds, and regenerates everything that derives from it:
# the Go module dependency, the Terraform schema, the scraped provider
# documentation, the API types, the controllers, the CRDs and the examples.
#
# Usage:
#   hack/update-terraform-provider.sh            # update to the latest release
#   hack/update-terraform-provider.sh 0.0.26     # update to a specific version
#   hack/update-terraform-provider.sh --dry-run  # only report what would happen
#
# When running inside GitHub Actions the outcome is also written to
# $GITHUB_OUTPUT as the "current", "target", "available" and "updated" outputs.
# "available" is true whenever the target differs from the current version;
# "updated" is true only when the working tree was actually regenerated, so it
# is always false for a dry run.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

REPO="${TERRAFORM_PROVIDER_GITHUB_REPO:-coder/terraform-provider-coderd}"

log() { echo "==> $*"; }

set_output() {
  if [ -n "${GITHUB_OUTPUT:-}" ]; then
    echo "$1=$2" >> "${GITHUB_OUTPUT}"
  fi
}

current_version() {
  sed -nE 's/^export TERRAFORM_PROVIDER_VERSION \?= (.+)$/\1/p' Makefile
}

latest_version() {
  local args=(-fsSL -H "Accept: application/vnd.github+json")
  if [ -n "${GITHUB_TOKEN:-}" ]; then
    args+=(-H "Authorization: Bearer ${GITHUB_TOKEN}")
  fi
  curl "${args[@]}" "https://api.github.com/repos/${REPO}/releases/latest" \
    | sed -nE 's/.*"tag_name"[[:space:]]*:[[:space:]]*"v?([^"]+)".*/\1/p' \
    | head -1
}

# The upjet code generator shells out to goimports. Build the version pinned by
# go.mod if it is not already on PATH.
ensure_goimports() {
  if command -v goimports >/dev/null 2>&1; then
    return
  fi
  local bin
  bin="$(go env GOPATH)/bin"
  log "building goimports into ${bin}"
  mkdir -p "${bin}"
  go build -o "${bin}/goimports" golang.org/x/tools/cmd/goimports
  export PATH="${bin}:${PATH}"
}

DRY_RUN=false
TARGET=""
while [ $# -gt 0 ]; do
  case "$1" in
    --dry-run) DRY_RUN=true ;;
    -*)
      echo "unknown flag: $1" >&2
      exit 1
      ;;
    *) TARGET="$1" ;;
  esac
  shift
done

if [ -z "${TARGET}" ]; then
  log "resolving the latest release of ${REPO}"
  TARGET="$(latest_version)"
fi
TARGET="${TARGET#v}"
if [ -z "${TARGET}" ]; then
  echo "cannot determine the target Terraform provider version" >&2
  exit 1
fi

CURRENT="$(current_version)"
if [ -z "${CURRENT}" ]; then
  echo "cannot read TERRAFORM_PROVIDER_VERSION from the Makefile" >&2
  exit 1
fi

set_output current "${CURRENT}"
set_output target "${TARGET}"

if [ "${CURRENT}" = "${TARGET}" ]; then
  log "already generated from ${REPO} v${CURRENT}, nothing to do"
  set_output available false
  set_output updated false
  exit 0
fi

set_output available true

if [ "${DRY_RUN}" = true ]; then
  log "dry run: v${CURRENT} would be updated to v${TARGET}, leaving the tree untouched"
  set_output updated false
  exit 0
fi

log "updating from v${CURRENT} to v${TARGET}"
sed -i -E "s|^export TERRAFORM_PROVIDER_VERSION \?= .*$|export TERRAFORM_PROVIDER_VERSION ?= ${TARGET}|" Makefile

# The provider is linked into the controller, so the Go module must move to the
# same release as the schema the CRDs are generated from.
log "updating the embedded provider module to v${TARGET}"
go get "github.com/coder/terraform-provider-coderd@v${TARGET}"
go mod tidy

# The schema and the docs checkout are both tied to the provider version.
rm -rf .work/coder config/schema.json

ensure_goimports
log "regenerating the provider"
make generate

set_output updated true
log "regenerated from ${REPO} v${TARGET}"
