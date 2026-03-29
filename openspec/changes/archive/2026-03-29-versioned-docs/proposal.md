## Why

The docs site currently deploys a single version from `main`. Users on different release channels (stable vs dev) see the same docs, which may describe features that don't exist in their version. As the project evolves, there's no way to reference docs for a prior release. Adding versioned documentation lets each release carry its own snapshot and gives users a version selector to find docs matching their installed version.

## What Changes

- Add `mike` as the versioning tool for MkDocs Material, replacing direct `mkdocs gh-deploy`
- Modify `docs.yml` workflow to deploy `dev` version on push to main via `mike deploy dev`
- Add a docs deployment step to `release.yml` that deploys the tagged version with a `latest` alias via `mike deploy <version> latest --update-aliases`
- Add `extra.version.provider: mike` to `mkdocs.yml` to enable the built-in version selector dropdown
- Root URL (`/`) redirects to `/latest/`, which resolves to the most recent stable release
- No backfill of historical versions — the selector accumulates versions naturally as future releases are tagged

## Capabilities

### New Capabilities
- `versioned-docs`: Version selector and multi-version deployment for the docs site using mike. Covers the versioning scheme, URL structure, alias management, and workflow integration.

### Modified Capabilities
- `docs-site`: The GitHub Pages deployment requirement changes from single-version `mkdocs gh-deploy --force` to multi-version deployment via mike. Deployment now happens from two workflows (docs.yml for dev, release.yml for stable).

## Impact

- **Workflows**: `.github/workflows/docs.yml` changes from `mkdocs gh-deploy` to `mike deploy`. `.github/workflows/release.yml` gains a docs deployment step.
- **Config**: `mkdocs.yml` gains `extra.version.provider: mike`.
- **gh-pages branch**: Structure changes from flat site to versioned subdirectories. First mike deployment will restructure the branch (existing unversioned content is replaced).
- **Dependencies**: `mike` pip package added to CI workflows alongside `mkdocs-material`.
- **No Go code changes** — this is entirely CI/config.
