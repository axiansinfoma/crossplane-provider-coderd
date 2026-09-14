#!/usr/bin/env python3
"""
Find Terraform provider resources that are present in config/schema.json but
not yet generated (per config/generated.lst), and file one GitHub issue per
resource that isn't already tracked by an existing open issue.

A resource only reaches config/generated.lst once it has an entry in
config/external_name.go, so anything the upstream coderd provider grew that we
have not wired up yet shows up here.

Usage:
  schema_diff_issues.py --repo <owner/repo> [--generated-lst PATH] [--schema PATH]
                        [--dry-run] [--label LABEL]

Requires the `gh` CLI to be installed and authenticated (GH_TOKEN/GITHUB_TOKEN)
unless --dry-run is used.

Exit codes:
  0  Ran successfully (regardless of whether any resources/issues were found)
  2  Error (missing files, invalid JSON, gh CLI failure, etc.)
"""

import argparse
import json
import re
import subprocess
import sys


def load_generated_list(path):
    """Load the generated resource names. The file is a JSON array, as written
    by `make generated-lst`."""
    with open(path) as f:
        return set(json.load(f))


def load_schema_resource_names(path):
    """Load the resource type names from a Terraform provider schema JSON file."""
    with open(path) as f:
        data = json.load(f)
    provider_schemas = data.get("provider_schemas", {})
    for _, provider_data in provider_schemas.items():
        return set(provider_data.get("resource_schemas", {}).keys())
    return set()


def run_gh(args, check=True):
    """Run a `gh` CLI command and return its stdout."""
    result = subprocess.run(["gh", *args], capture_output=True, text=True, check=False)
    if check and result.returncode != 0:
        raise RuntimeError(
            f"gh {' '.join(args)} failed (exit {result.returncode}): {result.stderr.strip()}"
        )
    return result.stdout


def fetch_open_issues(repo):
    """Fetch all open issues (number, title, body, url) for the given repo."""
    out = run_gh(
        [
            "issue", "list",
            "--repo", repo,
            "--state", "open",
            "--limit", "500",
            "--json", "number,title,body,url",
        ]
    )
    return json.loads(out)


def resource_mentioned(resource, text):
    """True if resource appears as a whole token in text, so that
    coderd_agents_model does not match coderd_agents_default_model."""
    if not text:
        return False
    pattern = r"(?<![\w-])" + re.escape(resource) + r"(?![\w-])"
    return re.search(pattern, text) is not None


def find_existing_issue(resource, issues):
    """Return the first open issue that already mentions this resource, if any."""
    for issue in issues:
        if resource_mentioned(resource, issue.get("title", "")) or resource_mentioned(
            resource, issue.get("body", "")
        ):
            return issue
    return None


def build_issue_body(resource):
    return f"""The Terraform provider resource `{resource}` is present in
`config/schema.json` but missing from `config/generated.lst`, so it is not yet
exposed as a Crossplane managed resource.

## Steps to add it

1. Add an entry for `{resource}` to `config/external_name.go`, mapping its
   external name onto the Terraform ID.
2. Give it an API group, kind and any cross-resource references in
   `config/resources/resources.go`.
3. Run `make generate` to regenerate the API types, controllers, CRDs and
   examples. Both the cluster-scoped (`*.coderd.crossplane.io`) and the
   namespaced (`*.coderd.m.crossplane.io`) flavour are generated from the same
   configuration.
4. Add a hand-authored example manifest under `examples/cluster/<group>/` and
   `examples/namespaced/<group>/`.
5. Add the kind to the managed-resource table in `README.md` and to the docs
   site under `docs/content/docs/`.

_This issue was filed automatically by the `schema-diff-issues` workflow._
"""


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo", required=True, help="owner/repo to search/create issues in")
    parser.add_argument("--generated-lst", default="config/generated.lst")
    parser.add_argument("--schema", default="config/schema.json")
    parser.add_argument(
        "--dry-run", action="store_true", help="Only print what would happen; do not create issues"
    )
    parser.add_argument("--label", default="enhancement", help="Label to apply to created issues")
    args = parser.parse_args()

    try:
        generated = load_generated_list(args.generated_lst)
    except FileNotFoundError:
        print(f"Error: generated list file not found: {args.generated_lst}", file=sys.stderr)
        sys.exit(2)
    except json.JSONDecodeError as e:
        print(f"Error loading {args.generated_lst}: {e}", file=sys.stderr)
        sys.exit(2)

    try:
        schema_resources = load_schema_resource_names(args.schema)
    except (FileNotFoundError, json.JSONDecodeError) as e:
        print(f"Error loading schema: {e}", file=sys.stderr)
        sys.exit(2)

    missing = sorted(schema_resources - generated)
    if not missing:
        print(
            f"No missing resources: {args.generated_lst} covers all "
            f"{len(schema_resources)} resources in {args.schema}."
        )
        return

    print(f"Found {len(missing)} resource(s) missing from {args.generated_lst}:")
    for r in missing:
        print(f"  - {r}")
    print()

    try:
        issues = fetch_open_issues(args.repo)
    except RuntimeError as e:
        print(f"Error fetching open issues: {e}", file=sys.stderr)
        sys.exit(2)

    created = 0
    for resource in missing:
        existing = find_existing_issue(resource, issues)
        if existing:
            print(f"[skip] {resource}: already tracked by #{existing['number']} ({existing['url']})")
            continue

        title = f"feat: expose {resource} as a managed resource"
        body = build_issue_body(resource)

        if args.dry_run:
            print(f'[dry-run] would create issue: "{title}"')
            created += 1
            continue

        try:
            out = run_gh(
                [
                    "issue", "create",
                    "--repo", args.repo,
                    "--title", title,
                    "--body", body,
                    "--label", args.label,
                ]
            )
            print(f"[created] {resource}: {out.strip()}")
            created += 1
        except RuntimeError as e:
            print(f"[error] failed to create issue for {resource}: {e}", file=sys.stderr)

    print()
    print(
        f"Done. {created} issue(s) {'would be ' if args.dry_run else ''}created, "
        f"{len(missing) - created} resource(s) already tracked or skipped."
    )


if __name__ == "__main__":
    main()
