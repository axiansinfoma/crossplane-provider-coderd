# Documentation

The `provider-coderd` documentation site, built with [Hugo](https://gohugo.io/)
and the [Hextra](https://github.com/imfing/hextra) theme.

## Local development

```bash
cd docs
hugo server --buildDrafts
```

Serves the site at <http://localhost:1313/>. Hugo extended is required, and Go
must be on `PATH` — the theme is pulled in as a Hugo module.

## Build

```bash
cd docs
hugo --minify
```

## Content model

- Write guides by hand, for workflows people actually follow.
- Keep the resource pages curated. `../package/crds/` is the generated source of
  truth for the complete field schema of every kind; do not duplicate it here,
  because it goes stale the moment the upstream provider is bumped.
- A new managed resource needs a mention in `content/docs/resources.md`.

## Deployment

Pushes to `main` that touch `docs/**` are deployed to GitHub Pages by
`.github/workflows/deploy-docs.yml`. Pull requests that touch `docs/**` get a
preview under `/pr-preview/pr-<number>/` from
`.github/workflows/preview-docs.yml`.

Live site: <https://axiansinfoma.github.io/crossplane-provider-coderd/>
