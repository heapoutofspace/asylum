## 1. MkDocs Configuration

- [x] 1.1 Add `extra.version.provider: mike` to `mkdocs.yml`

## 2. Docs Workflow (dev channel)

- [x] 2.1 Update `.github/workflows/docs.yml` to install `mike` alongside `mkdocs-material`
- [x] 2.2 Replace `mkdocs gh-deploy --force` with `mike deploy dev --push` in docs.yml
- [x] 2.3 Configure git user identity in docs.yml (required by mike for gh-pages commits)

## 3. Release Workflow (stable channel)

- [x] 3.1 Add a docs deployment job to `.github/workflows/release.yml` that runs `mike deploy <version> latest --update-aliases --push`
- [x] 3.2 Strip the `v` prefix from the tag name to derive the version identifier
- [x] 3.3 Set the default version to `latest` via `mike set-default latest --push`
- [x] 3.4 Configure git user identity and checkout of gh-pages branch access in release.yml

## 4. Verification

- [x] 4.1 Check if a CNAME file exists on gh-pages and ensure it's preserved across mike deployments
- [x] 4.2 Add CHANGELOG entry under Unreleased
