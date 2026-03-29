## Context

The docs site uses MkDocs Material, deployed to GitHub Pages via `mkdocs gh-deploy --force` on pushes to `main`. This overwrites the entire `gh-pages` branch with a single unversioned build. The project has two release channels (stable tags `v*` and a rolling `dev` pre-release) but docs don't distinguish between them.

mike is the standard versioning tool for MkDocs Material. It deploys each version into its own subdirectory on the `gh-pages` branch and maintains a `versions.json` manifest that the Material theme's built-in version selector reads.

## Goals / Non-Goals

**Goals:**
- Version selector dropdown in the docs site header
- Dev channel docs deployed on every push to `main`
- Stable version docs deployed on every release tag
- Root URL redirects to the latest stable version
- Versions accumulate naturally — no manual cleanup needed

**Non-Goals:**
- Backfilling docs for historical releases (v0.1.0–v0.5.0)
- Pruning old versions automatically
- Separate docs content per version (all versions build from the same `docs/` directory at the point-in-time of that commit)

## Decisions

### Use mike for versioning

mike manages versioned deployments to gh-pages by writing each version into a subdirectory and maintaining a `versions.json` index. MkDocs Material has built-in support for it via `extra.version.provider: mike`.

**Alternatives considered:**
- **Manual subdirectory management**: Build into `site/<version>/` and push. Works but requires reimplementing alias management, redirects, and the version selector manifest. mike does all of this.
- **Docusaurus/other tool**: Would require migrating all existing docs. Massive overhead for one feature.

### Version naming scheme

- Stable releases: bare version number without `v` prefix (e.g., `0.6.0`, not `v0.6.0`). Matches how users think about versions and keeps URLs clean.
- Dev channel: literal string `dev`. Always overwritten on each push to main.
- The `latest` alias always points to the most recent stable release.

### Two-workflow deployment

- `docs.yml` handles `dev` — triggered on push to `main` (paths: `docs/**`, `mkdocs.yml`)
- `release.yml` handles stable — triggered on tag push, already runs the release build, gains a docs step

**Alternative considered:** Single dedicated docs workflow triggered by both events. Rejected because the release workflow already has the tag context and checkout. Adding a step there is simpler than coordinating between workflows.

### gh-pages branch migration

The first `mike deploy` will restructure the `gh-pages` branch. The existing flat content (from `mkdocs gh-deploy --force`) will be replaced by mike's versioned directory structure. This is a one-time, non-destructive transition — the old unversioned docs are simply superseded.

`mike set-default latest` ensures the root URL redirects to the latest stable version. This runs once during the first deployment and persists in the gh-pages branch.

## Risks / Trade-offs

**[First deploy creates a gap]** After the workflow changes merge but before the next release tag, only `dev` will exist in the version selector. The root URL redirect to `latest` won't resolve until a stable version is deployed.
**Mitigation:** The docs.yml workflow can run `mike set-default dev` initially. Once the first release deploys, it sets the default to `latest`. Alternatively, just accept the brief gap — the next release will fix it.

**[mike is an additional CI dependency]** Adds `mike` to pip install in both workflows.
**Mitigation:** mike is a lightweight, well-maintained package (only dependency is mkdocs itself). Pin or leave unpinned — either is fine given the simplicity.

**[gh-pages branch history rewrite]** mike's first deploy effectively replaces the branch content. Any custom files on gh-pages (CNAME, etc.) need to be preserved.
**Mitigation:** Check if a CNAME file exists on gh-pages. If so, ensure mike's deployment preserves it (mike does not delete files it didn't create, but `--force` on the first deploy might). Verify during implementation.
